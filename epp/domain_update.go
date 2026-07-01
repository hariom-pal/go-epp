package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

func (c *Client) DomainUpdate(
	req types.DomainUpdateRequest,
) (*types.DomainUpdateResponse, error) {

	domain, requestXML, err := buildDomainUpdateRequestXML(
		req,
		c.nextTRID("UPDATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseDomainUpdateResponseXML(responseXML, domain)
}

func buildDomainUpdateRequestXML(
	req types.DomainUpdateRequest,
	clientTRID string,
) (string, []byte, error) {

	domain := strings.TrimSpace(req.DomainName)
	if domain == "" {
		domain = strings.TrimSpace(req.Domain)
	}
	domain = strings.TrimSuffix(domain, ".")

	if domain == "" {
		return "", nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "domain name is required",
		}
	}

	ascii, err := idn.ToASCII(domain)
	if err != nil {
		return "", nil, err
	}

	add, err := domainUpdateAddRemove(
		req.AddNameServers,
		req.AddNameServerInfo,
		req.AddContacts,
		req.AddAdminContacts,
		req.AddTechContacts,
		req.AddBillingContacts,
		req.AddStatuses,
		req.AddStatusInfo,
	)
	if err != nil {
		return "", nil, err
	}

	remove, err := domainUpdateAddRemove(
		req.RemoveNameServers,
		req.RemoveNameServerInfo,
		req.RemoveContacts,
		req.RemoveAdminContacts,
		req.RemoveTechContacts,
		req.RemoveBillingContacts,
		req.RemoveStatuses,
		req.RemoveStatusInfo,
	)
	if err != nil {
		return "", nil, err
	}

	change := domainUpdateChange(req)

	if add == nil &&
		remove == nil &&
		change == nil {

		return "", nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "at least one domain update operation is required",
		}
	}

	request := domainUpdateRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainUpdateCommandXML{
			ClientTRID: clientTRID,
			Update: domainUpdateXML{
				Domain: domainUpdateObjectXML{
					Name:   ascii,
					Add:    add,
					Remove: remove,
					Change: change,
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
		return "", nil, err
	}

	requestXML = append([]byte(xml.Header), requestXML...)

	return domain, requestXML, nil
}

func parseDomainUpdateResponseXML(
	responseXML []byte,
	domain string,
) (*types.DomainUpdateResponse, error) {

	var response domainUpdateResponseXML

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

	return &types.DomainUpdateResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Result: types.DomainUpdateResult{
			Domain: domain,
		},
	}, nil
}

func domainUpdateAddRemove(
	nameServers []string,
	nameServerInfo []types.DomainNameServer,
	contacts []types.DomainContact,
	adminContacts []string,
	techContacts []string,
	billingContacts []string,
	statuses []string,
	statusInfo []types.DomainStatus,
) (*domainUpdateAddRemoveXML, error) {

	nameServerXML, err := domainUpdateNameServers(nameServers, nameServerInfo)
	if err != nil {
		return nil, err
	}

	contactXML, err := domainUpdateContacts(
		contacts,
		adminContacts,
		techContacts,
		billingContacts,
	)
	if err != nil {
		return nil, err
	}

	statusXML, err := domainUpdateStatuses(statuses, statusInfo)
	if err != nil {
		return nil, err
	}

	if nameServerXML == nil &&
		len(contactXML) == 0 &&
		len(statusXML) == 0 {

		return nil, nil
	}

	return &domainUpdateAddRemoveXML{
		NameServers: nameServerXML,
		Contacts:    contactXML,
		Statuses:    statusXML,
	}, nil
}

func domainUpdateNameServers(
	nameServers []string,
	nameServerInfo []types.DomainNameServer,
) (*domainCreateNameServersXML, error) {

	return domainCreateNameServers(types.DomainCreateRequest{
		NameServers:    nameServers,
		NameServerInfo: nameServerInfo,
	})
}

func domainUpdateContacts(
	contacts []types.DomainContact,
	adminContacts []string,
	techContacts []string,
	billingContacts []string,
) ([]domainCreateContactXML, error) {

	result := make([]domainCreateContactXML, 0,
		len(contacts)+
			len(adminContacts)+
			len(techContacts)+
			len(billingContacts),
	)

	for _, contact := range contacts {
		next, err := domainCreateContact(contact.Type, contact.ID, true)
		if err != nil {
			return nil, err
		}
		result = append(result, next)
	}

	legacyContacts := []struct {
		Type string
		IDs  []string
	}{
		{Type: "admin", IDs: adminContacts},
		{Type: "tech", IDs: techContacts},
		{Type: "billing", IDs: billingContacts},
	}

	for _, group := range legacyContacts {
		for _, id := range group.IDs {
			next, err := domainCreateContact(group.Type, id, true)
			if err != nil {
				return nil, err
			}
			result = append(result, next)
		}
	}

	return result, nil
}

func domainUpdateStatuses(
	statuses []string,
	statusInfo []types.DomainStatus,
) ([]domainUpdateStatusXML, error) {

	result := make([]domainUpdateStatusXML, 0, len(statuses)+len(statusInfo))

	for _, status := range statuses {
		status = strings.TrimSpace(status)
		if status == "" {
			return nil, &Error{
				Code:    constants.ResultParameterError,
				Message: "domain status is required",
			}
		}

		result = append(result, domainUpdateStatusXML{
			Status: status,
		})
	}

	for _, status := range statusInfo {
		statusValue := strings.TrimSpace(status.Status)
		if statusValue == "" {
			return nil, &Error{
				Code:    constants.ResultParameterError,
				Message: "domain status is required",
			}
		}

		result = append(result, domainUpdateStatusXML{
			Status: statusValue,
			Lang:   strings.TrimSpace(status.Lang),
			Text:   strings.TrimSpace(status.Text),
		})
	}

	return result, nil
}

func domainUpdateChange(
	req types.DomainUpdateRequest,
) *domainUpdateChangeXML {

	change := &domainUpdateChangeXML{}

	registrant := strings.TrimSpace(req.Registrant)
	if registrant == "" {
		registrant = strings.TrimSpace(req.NewRegistrant)
	}
	if registrant != "" {
		change.Registrant = registrant
	}

	if authInfo := strings.TrimSpace(req.AuthInfo); authInfo != "" {
		change.AuthInfo = &domainCreateAuthInfoXML{
			Password: authInfo,
		}
	}

	if change.Registrant == "" &&
		change.AuthInfo == nil {

		return nil
	}

	return change
}
