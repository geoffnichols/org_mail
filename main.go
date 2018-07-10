package main

import (
	"github.com/jtblin/go-ldap-client"
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"log"
)

func main() {

	viper.BindEnv("LDAP_USERNAME")
	viper.BindEnv("LDAP_PASSWORD")
	ldap_username := viper.Get("LDAP_USERNAME").(string)
	ldap_password := viper.Get("LDAP_PASSWORD").(string)
	bind_dn := "uid=" + ldap_username + ",ou=users,dc=puppetlabs,dc=com"

	client := &ldap.LDAPClient{
		Base:         "dc=puppetlabs,dc=com",
		Host:         "ldap.puppetlabs.com",
		ServerName:   "ldap.puppetlabs.com",
		Port:         636,
		UseSSL:       true,
		BindDN:       bind_dn,
		BindPassword: ldap_password,
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"sn", "mail", "uid"},
	}
	// It is the responsibility of the caller to close the connection

	defer client.Close()

	client.ServerName = "ldap.puppetlabs.com"

	ok, user, err := client.Authenticate(ldap_username, ldap_password)
	pretty.Println(user)
	if err != nil {
		log.Fatalf("Error authenticating user %s: %+v", ldap_username, err)
	}
	if !ok {
		log.Fatalf("Authenticating failed for user %s", ldap_username)
	}
	log.Printf("User: %+v", user)

}
