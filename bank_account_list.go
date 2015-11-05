/*
!!! DO NOT MODIFY !!!

autogenerated
on: Thu Nov 05 02:49:57 +0000 2015
by: chakrit
*/

package omise

type BankAccountList struct {
	List
  Data []*BankAccount `json:"data"`
}

func (list *BankAccountList) Find(id string) *BankAccount {
	for _, item := range list.Data {
		if item.ID == id {
			return item
		}
	}

	return nil
}
