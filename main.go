package main

import (
	//	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gopkg.in/mup.v0/ldap"
	"os"
	//	"strconv"
	//	"reflect"
	"fmt"
	"sort"
	"strings"
)

//func find_root()

func is_manager(conn ldap.Conn, uid string) bool {
	manager_dn := "uid=" + uid + ",ou=users,dc=puppetlabs,dc=com"
	search := &ldap.Search{
		Filter: "(&(manager=" + manager_dn + ")(!(employeeType=Intern))(!(objectclass=exPuppetPerson)))",
		Attrs:  []string{"cn", "mail", "uid", "manager", "title"},
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
		Attrs:  []string{"cn", "mail", "uid", "manager", "title"},
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

type Ldapentry struct {
	dn      string
	uid     string
	mail    string
	manager string
	title   string
	cn      string
}

func build_entry(conn ldap.Conn, uid string) Ldapentry {
	log.Debug("In build entry for " + uid)
	var entry Ldapentry
	search := &ldap.Search{
		Filter: "(uid=" + uid + ")",
		Attrs:  []string{"uid", "mail", "cn", "title", "manager"},
	}
	results, err := conn.Search(search)
	if err != nil {
		panic(err)
	}
	for _, item := range results {
		entry.dn = item.DN
		entry.manager = item.Value("manager")
		entry.mail = item.Value("mail")
		entry.uid = item.Value("uid")
		entry.title = item.Value("title")
		entry.cn = item.Value("cn")
	}
	return entry
}

func peer_sort(peers []Ldapentry) []Ldapentry {
	sort.Slice(peers[:], func(i, j int) bool {
		return peers[i].mail < peers[j].mail
	})
	return peers
}

// need to rework to a struct vs string array
func build_tree(conn ldap.Conn, uid string) []Ldapentry {
	log.Debug("In build tree, and uid string passed is " + uid)
	//	var dn string
	var entry Ldapentry
	_, uid = shrink_dn(uid)
	peers := []Ldapentry{}
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
		return peer_sort(peers)

	} else {
		entry = build_entry(conn, uid)
		peers = append(peers, entry)
		log.Debug("In build tree, and uid " + uid + " is *not* a manager")
		return peers
	}

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
	reports := build_tree(conn, "stahnma")
	//build_tree(conn, "bradejr")

	//rsize := strconv.Itoa(len(reports))
	//	fmt.Println("I found " + rsize + " reports")
	// sort by mail?
	for _, v := range reports {
		fmt.Println(v.mail)
	}

}
