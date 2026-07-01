package epp

import (
	"encoding/xml"
	"strings"
	"time"

	"github.com/hariom-pal/go-epp/constants"
	feeext "github.com/hariom-pal/go-epp/extensions/fee"
	launchext "github.com/hariom-pal/go-epp/extensions/launch"
	rgpext "github.com/hariom-pal/go-epp/extensions/rgp"
	secdnsext "github.com/hariom-pal/go-epp/extensions/secdns"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

// DomainInfo retrieves RFC5731 information for a domain.
func (c *Client) DomainInfo(
	req types.DomainInfoRequest,
) (*types.DomainInfoResponse, error) {

	domain := strings.TrimSpace(req.Domain)
	domain = strings.TrimSuffix(domain, ".")

	if domain == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "domain is required",
		}
	}

	ascii, err := idn.ToASCII(domain)
	if err != nil {
		return nil, err
	}

	hosts := strings.TrimSpace(req.Hosts)
	if hosts == "" {
		hosts = constants.HostsDefault
	}

	if !constants.IsHostsValue(hosts) {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "invalid hosts value",
		}
	}

	if !launchext.ValidInfo(req.Launch) {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "invalid launch info extension",
		}
	}

	request := domainInfoRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainInfoCommandXML{
			ClientTRID: c.nextTRID("INFO"),
			Extension:  domainInfoExtension(req),
			Info: domainInfoXML{
				Domain: domainInfoObjectXML{
					Name: domainInfoNameXML{
						Hosts: hosts,
						Value: ascii,
					},
				},
			},
		},
	}

	requestXML, err := xml.MarshalIndent(
		request,
		"",
		"    ",
	)
	if err != nil {
		return nil, err
	}

	requestXML = append([]byte(xml.Header), requestXML...)

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	var response domainInfoResponseXML

	if err := xml.Unmarshal(responseXML, &response); err != nil {
		return nil, err
	}

	if response.Response.Result.Code != constants.ResultSuccess &&
		response.Response.Result.Code != constants.ResultSuccessPending {

		return nil, &Error{
			Code:       response.Response.Result.Code,
			Message:    response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		}
	}

	info := response.Response.ResData.InfoData
	rgpInfo := rgpext.InfoDataFromXML(response.Response.Extension.RGPInfoData)
	secDNSInfo := secdnsext.InfoDataFromXML(response.Response.Extension.SecDNSInfoData)
	feeInfo := feeext.InfoDataFromXML(response.Response.Extension.FeeInfoData)
	launchInfo := launchext.InfoDataFromXML(response.Response.Extension.LaunchInfoData)
	idnInfo := response.Response.Extension.IDNInfoData

	unicode, err := idn.ToUnicode(info.Name)
	if err != nil {
		unicode = info.Name
	}

	resp := &types.DomainInfoResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Result: types.DomainInfoResult{
			Domain:       unicode,
			ASCII:        info.Name,
			ROID:         info.ROID,
			Registrant:   info.Registrant,
			Registrar:    info.Registrar,
			CreatedBy:    info.CreatedBy,
			UpdatedBy:    info.UpdatedBy,
			AuthInfo:     strings.TrimSpace(info.AuthInfo.Password.Value),
			AuthInfoROID: info.AuthInfo.Password.ROID,
			RGP:          rgpInfo,
			DNSSEC: types.DomainDNSSECInfo{
				MaxSigLife: secDNSInfo.MaxSigLife,
				DSData:     make([]types.DomainDSData, 0, len(secDNSInfo.DSData)),
				KeyData:    make([]types.DomainKeyData, 0, len(secDNSInfo.KeyData)),
			},
			SecDNS: secDNSInfo,
			Fee: types.DomainFeeInfo{
				Currency: strings.TrimSpace(feeInfo.Currency),
				Commands: make([]types.DomainFeeCommand, 0, 1),
				Fees:     make([]types.DomainFeeAmount, 0, len(feeInfo.Fees)),
				Credits:  make([]types.DomainFeeAmount, 0, len(feeInfo.Credits)),
			},
			Launch: types.DomainLaunchInfo{
				Phase:         strings.TrimSpace(launchInfo.Phase.Value),
				ApplicationID: strings.TrimSpace(launchInfo.ApplicationID),
				Status:        launchInfo.Status.Status,
				StatusText:    strings.TrimSpace(launchInfo.Status.Text),
				StatusLang:    launchInfo.Status.Lang,
			},
			LaunchData: launchInfo,
			IDN: types.DomainIDNInfo{
				Table: strings.TrimSpace(idnInfo.Table),
			},

			Statuses:       make([]string, 0, len(info.Statuses)),
			StatusDetails:  make([]types.DomainStatus, 0, len(info.Statuses)),
			Contacts:       make([]types.DomainContact, 0, len(info.Contacts)),
			NameServers:    make([]string, 0, len(info.NameServers.HostObjects)+len(info.NameServers.HostAttrs)),
			NameServerInfo: make([]types.DomainNameServer, 0, len(info.NameServers.HostObjects)+len(info.NameServers.HostAttrs)),
			RGPStatuses:    make([]string, 0, len(rgpInfo.Statuses)),

			CreatedDate:  parseEPPDateTime(info.CreatedDate),
			UpdatedDate:  parseEPPDateTime(info.UpdatedDate),
			ExpiryDate:   parseEPPDateTime(info.ExpiryDate),
			TransferDate: parseEPPDateTime(info.TransferDate),
		},
	}

	for _, status := range info.Statuses {
		if status.Value == "" {
			continue
		}
		resp.Result.Statuses = append(resp.Result.Statuses, status.Value)
		resp.Result.StatusDetails = append(resp.Result.StatusDetails, types.DomainStatus{
			Status: status.Value,
			Lang:   status.Lang,
			Text:   strings.TrimSpace(status.Text),
		})
	}

	for _, status := range rgpInfo.Statuses {
		resp.Result.RGPStatuses = append(resp.Result.RGPStatuses, status.Status)
	}

	for _, ds := range secDNSInfo.DSData {
		dsData := types.DomainDSData{
			KeyTag:     ds.KeyTag,
			Algorithm:  ds.Algorithm,
			DigestType: ds.DigestType,
			Digest:     strings.TrimSpace(ds.Digest),
		}

		if ds.KeyData != nil &&
			(ds.KeyData.PublicKey != "" ||
				ds.KeyData.Flags != 0 ||
				ds.KeyData.Protocol != 0 ||
				ds.KeyData.Algorithm != 0) {

			dsData.KeyData = &types.DomainKeyData{
				Flags:     ds.KeyData.Flags,
				Protocol:  ds.KeyData.Protocol,
				Algorithm: ds.KeyData.Algorithm,
				PublicKey: strings.TrimSpace(ds.KeyData.PublicKey),
			}
		}

		resp.Result.DNSSEC.DSData = append(resp.Result.DNSSEC.DSData, dsData)
	}

	for _, key := range secDNSInfo.KeyData {
		resp.Result.DNSSEC.KeyData = append(resp.Result.DNSSEC.KeyData, types.DomainKeyData{
			Flags:     key.Flags,
			Protocol:  key.Protocol,
			Algorithm: key.Algorithm,
			PublicKey: strings.TrimSpace(key.PublicKey),
		})
	}

	if feeInfo.Command.Name != "" {
		resp.Result.Fee.Commands = append(resp.Result.Fee.Commands, types.DomainFeeCommand{
			Name:     feeInfo.Command.Name,
			Phase:    feeInfo.Command.Phase,
			Subphase: feeInfo.Command.Subphase,
		})
	}

	for _, fee := range feeInfo.Fees {
		resp.Result.Fee.Fees = append(resp.Result.Fee.Fees, types.DomainFeeAmount{
			Amount:      strings.TrimSpace(fee.Amount),
			Description: fee.Description,
			Refundable:  fee.Refundable,
			GracePeriod: fee.GracePeriod,
		})
	}

	for _, credit := range feeInfo.Credits {
		resp.Result.Fee.Credits = append(resp.Result.Fee.Credits, types.DomainFeeAmount{
			Amount:      strings.TrimSpace(credit.Amount),
			Description: credit.Description,
			Refundable:  credit.Refundable,
			GracePeriod: credit.GracePeriod,
		})
	}

	for _, contact := range info.Contacts {
		id := strings.TrimSpace(contact.ID)
		if id == "" {
			continue
		}
		resp.Result.Contacts = append(resp.Result.Contacts, types.DomainContact{
			Type: contact.Type,
			ID:   id,
		})
	}

	for _, host := range info.NameServers.HostObjects {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}
		resp.Result.NameServers = append(resp.Result.NameServers, host)
		resp.Result.NameServerInfo = append(resp.Result.NameServerInfo, types.DomainNameServer{
			HostName: host,
		})
	}

	for _, hostAttr := range info.NameServers.HostAttrs {
		host := strings.TrimSpace(hostAttr.HostName)
		if host == "" {
			continue
		}
		resp.Result.NameServers = append(resp.Result.NameServers, host)

		nameServer := types.DomainNameServer{
			HostName:  host,
			Addresses: make([]types.DomainHostAddress, 0, len(hostAttr.HostAddr)),
		}

		for _, addr := range hostAttr.HostAddr {
			ip := strings.TrimSpace(addr.IP)
			if ip == "" {
				continue
			}
			nameServer.Addresses = append(nameServer.Addresses, types.DomainHostAddress{
				IP:      ip,
				Version: addr.Version,
			})
		}

		resp.Result.NameServerInfo = append(resp.Result.NameServerInfo, nameServer)
	}

	return resp, nil
}

func domainInfoExtension(
	req types.DomainInfoRequest,
) *domainInfoExtensionXML {

	launchInfo := launchext.NewInfo(req.Launch)
	if launchInfo == nil {
		return nil
	}

	return &domainInfoExtensionXML{
		LaunchInfo: launchInfo,
	}
}

func parseEPPDateTime(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}

	return &parsed
}
