// Package bankaccount provides the /bank_accounts APIs
package bankaccount

import (
	"errors"

	stripe "github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/form"
)

// Client is used to invoke /bank_accounts APIs.
type Client struct {
	B   stripe.Backend
	Key string
}

// New POSTs a new bank account.
func New(params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().New(params)
}

func (c Client) New(params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	body := &form.Values{}

	// Note that we call this special append method instead of the standard one
	// from the form package. We should not use form's because doing so will
	// include some parameters that are undesirable here.
	params.AppendToAsSourceOrExternalAccount(body)

	ba := &stripe.BankAccount{}
	var err error
	if params.Customer != nil {
		err = c.B.Call("POST", stripe.FormatURLPath("/customers/%s/sources", stripe.StringValue(params.Customer)), c.Key, body, &params.Params, ba)
	} else {
		err = c.B.Call("POST", stripe.FormatURLPath("/accounts/%s/external_accounts", stripe.StringValue(params.Account)), c.Key, body, &params.Params, ba)
	}

	return ba, err
}

// Get returns the details of a bank account.
func Get(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().Get(id, params)
}

func (c Client) Get(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	var body *form.Values
	var commonParams *stripe.Params

	if params != nil {
		commonParams = &params.Params
		body = &form.Values{}
		form.AppendTo(body, params)
	}

	ba := &stripe.BankAccount{}
	var err error

	if params != nil && params.Customer != nil {
		err = c.B.Call("GET", stripe.FormatURLPath("/customers/%s/sources/%s", stripe.StringValue(params.Customer), id), c.Key, body, commonParams, ba)
	} else if params != nil && params.Account != nil {
		err = c.B.Call("GET", stripe.FormatURLPath("/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id), c.Key, body, commonParams, ba)
	} else {
		err = errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	return ba, err
}

// Update updates a bank account.
func Update(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().Update(id, params)
}

func (c Client) Update(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	var body *form.Values
	var commonParams *stripe.Params

	if params != nil {
		commonParams = &params.Params
		body = &form.Values{}
		form.AppendTo(body, params)
	}

	ba := &stripe.BankAccount{}
	var err error

	if params.Customer != nil {
		err = c.B.Call("POST", stripe.FormatURLPath("/customers/%s/sources/%s", stripe.StringValue(params.Customer), id), c.Key, body, commonParams, ba)
	} else if params.Account != nil {
		err = c.B.Call("POST", stripe.FormatURLPath("/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id), c.Key, body, commonParams, ba)
	} else {
		err = errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	return ba, err
}

// Del removes a bank account.
func Del(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	return getC().Del(id, params)
}

func (c Client) Del(id string, params *stripe.BankAccountParams) (*stripe.BankAccount, error) {
	var body *form.Values
	var commonParams *stripe.Params

	if params != nil {
		commonParams = &params.Params
		body = &form.Values{}
		form.AppendTo(body, params)
	}

	ba := &stripe.BankAccount{}
	var err error

	if params.Customer != nil {
		err = c.B.Call("DELETE", stripe.FormatURLPath("/customers/%s/sources/%s", stripe.StringValue(params.Customer), id), c.Key, body, commonParams, ba)
	} else if params.Account != nil {
		err = c.B.Call("DELETE", stripe.FormatURLPath("/accounts/%s/external_accounts/%s", stripe.StringValue(params.Account), id), c.Key, body, commonParams, ba)
	} else {
		err = errors.New("Invalid bank account params: either Customer or Account need to be set")
	}

	return ba, err
}

// List returns a list of bank accounts.
func List(params *stripe.BankAccountListParams) *Iter {
	return getC().List(params)
}

func (c Client) List(params *stripe.BankAccountListParams) *Iter {
	body := &form.Values{}
	var lp *stripe.ListParams
	var p *stripe.Params

	form.AppendTo(body, params)
	lp = &params.ListParams
	p = params.ToParams()

	return &Iter{stripe.GetIter(lp, body, func(b *form.Values) ([]interface{}, stripe.ListMeta, error) {
		list := &stripe.BankAccountList{}
		var err error

		if params.Customer != nil {
			err = c.B.Call("GET", stripe.FormatURLPath("/customers/%s/sources?object=bank_account", stripe.StringValue(params.Customer)), c.Key, b, p, list)
		} else if params.Account != nil {
			err = c.B.Call("GET", stripe.FormatURLPath("/accounts/%s/external_accounts?object=bank_account", stripe.StringValue(params.Account)), c.Key, b, p, list)
		} else {
			err = errors.New("Invalid bank account params: either Customer or Account need to be set")
		}

		ret := make([]interface{}, len(list.Data))
		for i, v := range list.Data {
			ret[i] = v
		}

		return ret, list.ListMeta, err
	})}
}

// Iter is an iterator for lists of BankAccount.
// The embedded Iter carries methods with it;
// see its documentation for details.
type Iter struct {
	*stripe.Iter
}

// BankAccount returns the most recent BankAccount
// visited by a call to Next.
func (i *Iter) BankAccount() *stripe.BankAccount {
	return i.Current().(*stripe.BankAccount)
}

func getC() Client {
	return Client{stripe.GetBackend(stripe.APIBackend), stripe.Key}
}
