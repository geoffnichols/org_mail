package main

import (
	"github.com/kr/pretty"
	"github.com/spf13/viper"
	"gopkg.in/mup.v0/ldap"
	//	"strconv"
	//	"reflect"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

//func find_root()

func is_manager(conn ldap.Conn, uid string) bool {
	manager_dn := "uid=" + uid + ",ou=users,dc=puppetlabs,dc=com"
	search := &ldap.Search{
		Filter: "(&(manager=" + manager_dn + ")(!(employeeType=Intern))(!(objectclass=exPuppetPerson)))",
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

func shrink_dn(uid string) (string, string) {
	dn := uid
	if strings.Contains(uid, "dc=") {
		idx := strings.Index(uid, "ou=")
		uid = uid[0:idx]
		uid = strings.Trim(uid, ",")
		uid = strings.Replace(uid, "uid=", "", 1)
	}
	if strings.Contains(uid, "dc=") == false {
		dn = "uid=" + uid + ",ou=users,dc=puppetlabs,dc=com"
	}

	return dn, uid
}

// need to rework to a struct vs string array
func build_tree(conn ldap.Conn, uid string) []string {
	log.Debug("In build tree, and uid string passed is " + uid)
	var dn string
	dn, uid = shrink_dn(uid)
	peers := []string{}
	if is_manager(conn, uid) {
		log.Debug("In build tree, and uid " + uid + " is a manager")
		for _, res := range who_has_this_manager(conn, uid) {
			// evalute if any of these people are managers
			//		pretty.Println(res)
			//fmt.Println(reflect.TypeOf(res))
			peers = append(peers, build_tree(conn, res)...)
			//pretty.Println(build_tree(conn, res))
			//pretty.Println(peers)
		}
		return peers

	} else {
		peers = append(peers, dn)
		log.Debug("In build tree, and uid " + uid + " is *not* a manager")
		return peers
	}
	return peers

}

func main() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	//TODO log level configurable via envvar
	log.SetLevel(log.InfoLevel)

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
	//	pretty.Println(who_has_this_manager(conn, "erict"))
	//	pretty.Println(is_manager(conn, "bradejr"))
	pretty.Println(build_tree(conn, "stahnma"))

}
