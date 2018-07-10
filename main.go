package main

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"gopkg.in/mup.v0/ldap"
	"strconv"
)

//func is_manager()
//func who_has_this_manager()
//func find_root()

/*func find_reports_for(client, uid) {
  manager_dn = "uid=" + uid + "ou=users,dc=puppetlabs,dc=com"
  client.UserFilter="(manager=%s)"
}
*/

func main() {

	viper.BindEnv("LDAP_USERNAME")
	viper.BindEnv("LDAP_PASSWORD")
	ldap_username := viper.Get("LDAP_USERNAME").(string)
	ldap_password := viper.Get("LDAP_PASSWORD").(string)
	bind_dn := "uid=" + ldap_username + ",ou=users,dc=puppetlabs,dc=com"
	base_dn := "dc=puppetlabs,dc=com"

	config := &ldap.Config{
		URL:      "ldaps://ldap.puppetlabs.com:636",
		BaseDN:   base_dn,
		BindDN:   bind_dn,
		BindPass: ldap_password,
	}

	search := &ldap.Search{
		Filter: "(manager=uid=stahnma,ou=users,dc=puppetlabs,dc=com)",
		Attrs:  []string{"sn", "mail", "uid", "manager"},
	}

	conn, err := ldap.Dial(config)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	results, err := conn.Search(search)
	//pretty.Println(results)
	if err != nil {
		panic(err)
	}
	fmt.Println("Found " + strconv.Itoa(len(results)) + " managers")
	for _, item := range results {
		pretty.Println(item.DN)
	}

}
