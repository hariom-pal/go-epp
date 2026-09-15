package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/epp"
)

type operationPlan struct {
	names     []string
	transform bool
	objects   []string
}

func planOperations(options cliOptions) operationPlan {
	var plan operationPlan
	add := func(name string, transform bool, objects ...string) {
		plan.names = append(plan.names, name)
		plan.transform = plan.transform || transform
		for _, object := range objects {
			if strings.TrimSpace(object) != "" {
				plan.objects = append(plan.objects, object)
			}
		}
	}

	if options.ConnectOnly {
		add("connect", false)
	}
	if options.LoginOnly {
		add("login", false)
	}
	if options.Hello {
		add("hello", false)
	}
	if strings.TrimSpace(options.CheckDomains) != "" {
		add("domain-check", false, splitDomains(options.CheckDomains)...)
	}
	if strings.TrimSpace(options.InfoDomain) != "" {
		add("domain-info", false, options.InfoDomain)
	}
	if strings.TrimSpace(options.CreateDomain) != "" {
		add("domain-create", true, options.CreateDomain)
	}
	if strings.TrimSpace(options.DomainUpdateName) != "" {
		add("domain-update", true, options.DomainUpdateName)
	}
	if strings.TrimSpace(options.DomainRenewName) != "" {
		add("domain-renew", true, options.DomainRenewName)
	}
	if strings.TrimSpace(options.DomainDeleteName) != "" {
		add("domain-delete", true, options.DomainDeleteName)
	}
	if strings.TrimSpace(options.DomainTransferName) != "" {
		add("domain-transfer-"+strings.TrimSpace(options.DomainTransferOperation), isTransformTransfer(options.DomainTransferOperation), options.DomainTransferName)
	}
	if len(options.ContactCheckIDs) > 0 {
		add("contact-check", false, options.ContactCheckIDs...)
	}
	if strings.TrimSpace(options.ContactInfoID) != "" {
		add("contact-info", false, options.ContactInfoID)
	}
	if strings.TrimSpace(options.ContactCreateID) != "" {
		add("contact-create", true, options.ContactCreateID)
	}
	if strings.TrimSpace(options.ContactUpdateID) != "" {
		add("contact-update", true, options.ContactUpdateID)
	}
	if strings.TrimSpace(options.ContactDeleteID) != "" {
		add("contact-delete", true, options.ContactDeleteID)
	}
	if strings.TrimSpace(options.ContactTransferID) != "" {
		add("contact-transfer-"+strings.TrimSpace(options.ContactTransferOperation), isTransformTransfer(options.ContactTransferOperation), options.ContactTransferID)
	}
	if len(options.HostCheckNames) > 0 {
		add("host-check", false, options.HostCheckNames...)
	}
	if strings.TrimSpace(options.HostInfoName) != "" {
		add("host-info", false, options.HostInfoName)
	}
	if strings.TrimSpace(options.HostCreateName) != "" {
		add("host-create", true, options.HostCreateName)
	}
	if strings.TrimSpace(options.HostUpdateName) != "" {
		add("host-update", true, options.HostUpdateName)
	}
	if strings.TrimSpace(options.HostDeleteName) != "" {
		add("host-delete", true, options.HostDeleteName)
	}
	if options.Poll {
		add("poll-request", false)
	}
	if strings.TrimSpace(options.PollAckID) != "" {
		add("poll-ack", true, "msgID="+strings.TrimSpace(options.PollAckID))
	}

	return plan
}

func isTransformTransfer(operation string) bool {
	switch strings.ToLower(strings.TrimSpace(operation)) {
	case "request", "approve", "reject", "cancel":
		return true
	default:
		return false
	}
}

func validateUATSafety(cfg *epp.Config, options cliOptions, plan operationPlan) (string, error) {
	environment := strings.TrimSpace(options.Environment)
	if environment == "" && cfg != nil {
		environment = strings.TrimSpace(cfg.Environment)
	}
	if environment == "" {
		environment = "UNKNOWN"
	}

	if options.UAT && !isSafeUATEnvironment(environment) {
		return environment, fmt.Errorf("UAT safety mode requires environment to be clearly UAT/OT&E/test, got: %s", environment)
	}
	if options.UAT && len(plan.names) > 1 {
		return environment, fmt.Errorf("UAT safety mode allows exactly one requested API operation per run, got: %s", strings.Join(plan.names, ", "))
	}
	if plan.transform && !isSafeUATEnvironment(environment) {
		return environment, fmt.Errorf("refusing transform operation because environment is not clearly UAT/OT&E/test: %s", environment)
	}
	if plan.transform && !options.ConfirmTransform {
		return environment, errors.New("refusing transform operation without --confirm-transform")
	}

	return environment, nil
}

func isSafeUATEnvironment(environment string) bool {
	value := strings.ToLower(strings.TrimSpace(environment))
	return strings.Contains(value, "uat") ||
		strings.Contains(value, "ote") ||
		strings.Contains(value, "ot&e") ||
		strings.Contains(value, "test") ||
		strings.Contains(value, "sandbox")
}

func printTargetSummary(cfg *epp.Config, environment string, plan operationPlan) {
	fmt.Println("========== UAT TARGET SUMMARY ==========")
	if cfg != nil {
		fmt.Printf("Target Host : %s\n", cfg.Server.Host)
		fmt.Printf("Target Port : %d\n", cfg.Server.Port)
		if cfg.TLS.ServerName != "" {
			fmt.Printf("TLS SNI     : %s\n", cfg.TLS.ServerName)
		}
	}
	fmt.Printf("Environment : %s\n", environment)
	if len(plan.names) > 0 {
		fmt.Printf("Operation   : %s\n", strings.Join(plan.names, ", "))
	}
	if len(plan.objects) > 0 {
		fmt.Printf("Test Objects: %s\n", strings.Join(plan.objects, ", "))
	}
	if plan.transform {
		fmt.Println("Transform   : yes")
	} else {
		fmt.Println("Transform   : no")
	}
	fmt.Println("Secrets     : password/authInfo/private keys are never printed")
	fmt.Println("========================================")
}

type captureLogger struct {
	file *os.File
}

type captureEvent struct {
	Timestamp  string `json:"timestamp"`
	Operation  string `json:"operation"`
	DurationMS int64  `json:"duration_ms"`
	ResultCode int    `json:"result_code,omitempty"`
	ClientTRID string `json:"clTRID,omitempty"`
	ServerTRID string `json:"svTRID,omitempty"`
	Success    bool   `json:"success"`
	ErrorClass string `json:"error_class,omitempty"`
	Error      string `json:"error,omitempty"`
}

func newCaptureLogger(dir string) (*captureLogger, error) {
	dir = strings.TrimSpace(dir)
	if dir == "" {
		return nil, nil
	}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	name := "uat-" + time.Now().UTC().Format("20060102T150405Z") + ".jsonl"
	file, err := os.OpenFile(filepath.Join(dir, name), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, err
	}
	return &captureLogger{file: file}, nil
}

func (l *captureLogger) Close() error {
	if l == nil || l.file == nil {
		return nil
	}
	return l.file.Close()
}

func (l *captureLogger) EPPEvent(event epp.Event) {
	if l == nil || l.file == nil {
		return
	}
	entry := captureEvent{
		Timestamp:  time.Now().UTC().Format(time.RFC3339Nano),
		Operation:  string(event.Type),
		DurationMS: event.Duration.Milliseconds(),
		ResultCode: event.ResultCode,
		ClientTRID: event.ClientTRID,
		ServerTRID: event.ServerTRID,
		Success:    event.Err == nil,
	}
	if event.Command != "" {
		entry.Operation = event.Command
	}
	if event.Err != nil {
		entry.ErrorClass = classifyError(event.Err)
		entry.Error = sanitizeText(event.Err.Error())
	}
	_ = json.NewEncoder(l.file).Encode(entry)
}

func classifyError(err error) string {
	var ambiguous *epp.AmbiguousTransformError
	if errors.As(err, &ambiguous) {
		return string(epp.ErrorKindAmbiguousTransform)
	}
	var sdkErr *epp.SDKError
	if errors.As(err, &sdkErr) {
		return string(sdkErr.Kind)
	}
	var resultErr *epp.Error
	if errors.As(err, &resultErr) {
		return string(resultErr.Kind)
	}
	return "unknown"
}

func sanitizeText(value string) string {
	replacements := []string{"password", "authInfo", "private key", "client key", "secret", "token"}
	sanitized := value
	for _, item := range replacements {
		sanitized = strings.ReplaceAll(sanitized, item, "[redacted]")
		sanitized = strings.ReplaceAll(sanitized, strings.ToUpper(item), "[redacted]")
	}
	return sanitized
}
