package main

import (
	//	"fmt"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"gopkg.in/mup.v0/ldap"
	//	"strconv"
	//  "reflect"
)

//func find_root()

func is_manager(conn ldap.Conn, uid string) bool {
	manager_dn := "uid=" + uid + ",ou=users,dc=puppetlabs,dc=com"
	search := &ldap.Search{
		Filter: "(manager=" + manager_dn + ")",
		Attrs:  []string{"sn", "mail", "uid", "manager"},
	}
	results, err := conn.Search(search)
	if err != nil {
		panic(err)
	}
	if len(results) > 0 {
		return true
	}
	return false
}

func who_has_this_manager(conn ldap.Conn, uid string) []string {
	manager_dn := "uid=" + uid + ",ou=users,dc=puppetlabs,dc=com"
	search := &ldap.Search{
		Filter: "(manager=" + manager_dn + ")",
		Attrs:  []string{"sn", "mail", "uid", "manager"},
	}
	results, err := conn.Search(search)
	if err != nil {
		panic(err)
	}
	var dns []string
	for _, item := range results {
		dns = append(dns, item.DN)
	}
	return dns
}

func build_tree(conn *ldap.Conn, uid string) {
}

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

	conn, err := ldap.Dial(config)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	//  fmt.Println(reflect.TypeOf(conn))
	pretty.Println(who_has_this_manager(conn, "erict"))
	pretty.Println(is_manager(conn, "bradejr"))
	build_tree(&conn, "stahnma")

}
