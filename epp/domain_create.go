package epp

import (
	"encoding/xml"
	"strings"

	"github.com/hariom-pal/go-epp/constants"
	"github.com/hariom-pal/go-epp/pkg/idn"
	"github.com/hariom-pal/go-epp/types"
)

const (
	domainCreatePeriodUnitYears  = "y"
	domainCreatePeriodUnitMonths = "m"

	domainCreateMaxPeriod = 99
)

func (c *Client) DomainCreate(
	req types.DomainCreateRequest,
) (*types.DomainCreateResponse, error) {

	requestXML, err := buildDomainCreateRequestXML(
		req,
		c.nextTRID("CREATE"),
	)
	if err != nil {
		return nil, err
	}

	responseXML, err := c.Execute(requestXML)
	if err != nil {
		return nil, err
	}

	return parseDomainCreateResponseXML(responseXML)
}

func buildDomainCreateRequestXML(
	req types.DomainCreateRequest,
	clientTRID string,
) ([]byte, error) {

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

	if req.Period < 1 || req.Period > domainCreateMaxPeriod {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "period must be between 1 and 99",
		}
	}

	unit := strings.ToLower(strings.TrimSpace(req.Unit))
	if unit == "" {
		unit = domainCreatePeriodUnitYears
	}

	if unit != domainCreatePeriodUnitYears &&
		unit != domainCreatePeriodUnitMonths {

		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "period unit must be y or m",
		}
	}

	registrant := strings.TrimSpace(req.Registrant)
	if registrant == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "registrant contact is required",
		}
	}

	authInfo := strings.TrimSpace(req.AuthInfo)
	if authInfo == "" {
		return nil, &Error{
			Code:    constants.ResultParameterError,
			Message: "authInfo is required",
		}
	}

	contacts, err := domainCreateContacts(req)
	if err != nil {
		return nil, err
	}

	nameServers, err := domainCreateNameServers(req)
	if err != nil {
		return nil, err
	}

	request := domainCreateRequestXML{
		XMLNS:       constants.EPPNamespace,
		DomainXMLNS: constants.DomainNamespace,
		Command: domainCreateCommandXML{
			ClientTRID: clientTRID,
			Create: domainCreateXML{
				Domain: domainCreateObjectXML{
					Name: ascii,
					Period: &domainCreatePeriodXML{
						Unit:  unit,
						Value: req.Period,
					},
					NameServers: nameServers,
					Registrant:  registrant,
					Contacts:    contacts,
					AuthInfo: domainCreateAuthInfoXML{
						Password: authInfo,
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

	return requestXML, nil
}

func parseDomainCreateResponseXML(
	responseXML []byte,
) (*types.DomainCreateResponse, error) {

	var response domainCreateResponseXML

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

	createData := response.Response.ResData.CreateData

	unicode, err := idn.ToUnicode(createData.Name)
	if err != nil {
		unicode = createData.Name
	}

	resp := &types.DomainCreateResponse{
		Response: types.Response{
			ResultCode: response.Response.Result.Code,
			ResultMsg:  response.Response.Result.Msg,
			ClientTRID: response.Response.TRID.ClientTRID,
			ServerTRID: response.Response.TRID.ServerTRID,
		},
		Result: types.DomainCreateResult{
			Domain: unicode,
		},
	}

	if createdDate := parseEPPDateTime(createData.CreatedDate); createdDate != nil {
		resp.Result.CreatedDate = *createdDate
	}

	if expiryDate := parseEPPDateTime(createData.ExpiryDate); expiryDate != nil {
		resp.Result.ExpiryDate = *expiryDate
	}

	return resp, nil
}

func domainCreateContacts(
	req types.DomainCreateRequest,
) ([]domainCreateContactXML, error) {

	contacts := make([]domainCreateContactXML, 0,
		len(req.Contacts)+
			len(req.AdminContacts)+
			len(req.TechContacts)+
			len(req.BillingContacts)+3,
	)

	for _, contact := range req.Contacts {
		next, err := domainCreateContact(contact.Type, contact.ID, true)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, next)
	}

	legacyContacts := []types.DomainContact{
		{Type: "admin", ID: req.AdminContact},
		{Type: "tech", ID: req.TechContact},
		{Type: "billing", ID: req.BillingContact},
	}

	for _, contact := range legacyContacts {
		next, err := domainCreateContact(contact.Type, contact.ID, false)
		if err != nil {
			return nil, err
		}
		if next.ID == "" {
			continue
		}
		contacts = append(contacts, next)
	}

	for _, id := range req.AdminContacts {
		next, err := domainCreateContact("admin", id, true)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, next)
	}

	for _, id := range req.TechContacts {
		next, err := domainCreateContact("tech", id, true)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, next)
	}

	for _, id := range req.BillingContacts {
		next, err := domainCreateContact("billing", id, true)
		if err != nil {
			return nil, err
		}
		contacts = append(contacts, next)
	}

	return contacts, nil
}

func domainCreateContact(
	contactType string,
	id string,
	required bool,
) (domainCreateContactXML, error) {

	contactType = strings.TrimSpace(contactType)
	id = strings.TrimSpace(id)

	if !required && id == "" {
		return domainCreateContactXML{}, nil
	}

	if contactType == "" || id == "" {
		return domainCreateContactXML{}, &Error{
			Code:    constants.ResultParameterError,
			Message: "contact type and id are required",
		}
	}

	return domainCreateContactXML{
		Type: contactType,
		ID:   id,
	}, nil
}

func domainCreateNameServers(
	req types.DomainCreateRequest,
) (*domainCreateNameServersXML, error) {

	if len(req.NameServers) == 0 &&
		len(req.NameServerInfo) == 0 {

		return nil, nil
	}

	nameServers := &domainCreateNameServersXML{
		HostObjects: make([]string, 0, len(req.NameServers)),
		HostAttrs:   make([]domainCreateHostAttrXML, 0, len(req.NameServerInfo)),
	}

	for _, host := range req.NameServers {
		ascii, err := domainCreateHostName(host)
		if err != nil {
			return nil, err
		}
		nameServers.HostObjects = append(nameServers.HostObjects, ascii)
	}

	for _, host := range req.NameServerInfo {
		hostName, err := domainCreateHostName(host.HostName)
		if err != nil {
			return nil, err
		}

		hostAttr := domainCreateHostAttrXML{
			HostName:  hostName,
			HostAddrs: make([]domainCreateHostAddrXML, 0, len(host.Addresses)),
		}

		for _, addr := range host.Addresses {
			ip := strings.TrimSpace(addr.IP)
			if ip == "" {
				return nil, &Error{
					Code:    constants.ResultParameterError,
					Message: "host address is required",
				}
			}

			version := strings.TrimSpace(addr.Version)
			if version != "" &&
				version != "v4" &&
				version != "v6" {

				return nil, &Error{
					Code:    constants.ResultParameterError,
					Message: "host address version must be v4 or v6",
				}
			}

			hostAttr.HostAddrs = append(hostAttr.HostAddrs, domainCreateHostAddrXML{
				Version: version,
				Value:   ip,
			})
		}

		nameServers.HostAttrs = append(nameServers.HostAttrs, hostAttr)
	}

	return nameServers, nil
}

func domainCreateHostName(
	host string,
) (string, error) {

	host = strings.TrimSpace(host)
	host = strings.TrimSuffix(host, ".")

	if host == "" {
		return "", &Error{
			Code:    constants.ResultParameterError,
			Message: "name server host name is required",
		}
	}

	return idn.ToASCII(host)
}
