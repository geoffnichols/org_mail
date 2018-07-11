# org_mail

A relatively simple program that by outputs mail addresses of everybody who
works under `MANAGER`.

`MANAGER` defaults to stahnma, because I am that vain. The value of `MANAGER`
is an LDAP uid (which at puppet is also the name before the `@` sign in a mail
address.)

# Installation

`go get github.com/puppetlabs/org_mail`

You may want to add `$GOPATH:/bin` to your `$PATH`.

# Configuration

You should export `LDAP_USERNAME` and `LDAP_PASSWORD` with appropriate values
for Puppet's LDAP servers.


# Usage

`./org_mail`


# Another Manager

    > MANAGER=geoff.nichols ./org_mail
    branan@puppet.com
    casey.williams@puppet.com
    enis.inan@puppet.com
    erick@puppet.com
    ethan@puppet.com
    glenn.sarti@puppet.com
    james.pogran@puppet.com
    john.oconnor@puppet.com
    jonathan.morris@puppet.com
    michael.lombardi@puppet.com
    scott.garman@puppet.com
    scott.mcclellan@puppet.com
    sean.mcdonald@puppet.com
    william.hurt@puppet.com

# License

Apache 2.0
