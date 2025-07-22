module analyzer

go 1.22.4

require (
	golang.org/x/tools v0.12.0
	golang.org/x/tools/go/pointer v0.1.0-deprecated
)

require (
	golang.org/x/mod v0.12.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
)

// prevent blueprint in go.work to force higher version
replace golang.org/x/tools => golang.org/x/tools v0.12.0

//replace analyzer => .

//replace github.com/blueprint-uservices/blueprint/plugins => /home/vagrant/blueprint/plugins

//replace github.com/blueprint-uservices/blueprint/blueprint => /home/vagrant/blueprint/blueprint

//replace github.com/blueprint-uservices/blueprint/runtime => /home/vagrant/blueprint/runtime

//replace github.com/blueprint-uservices/blueprint/examples/postnotification_simple => /home/vagrant/blueprint/examples/postnotification_simple

//replace github.com/blueprint-uservices/blueprint/examples/postnotification_simple/workflow => /home/vagrant/blueprint/examples/postnotification_simple/workflow

//replace github.com/blueprint-uservices/blueprint/examples/postnotification_simple/wiring => /home/vagrant/blueprint/examples/postnotification_simple/wiring
