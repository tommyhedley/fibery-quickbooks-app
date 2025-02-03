package quickbooks

type InvoiceBody struct {
	Id string
}

type InvoiceType struct {
	Id   string
	Body InvoiceBody
}

func (i *InvoiceType) QueryAll(c *Client) (InvoiceBody, error) {
	return InvoiceBody{}, nil
}
