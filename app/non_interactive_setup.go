package main

import (
	"errors"
	"github.com/dedis/kyber/util/encoding"
	"github.com/dedis/kyber/util/key"
	"github.com/dedis/onet/app"
	"github.com/dedis/onet/log"
	"github.com/dedis/onet/network"
	"github.com/lca1/unlynx/lib"
	"gopkg.in/urfave/cli.v1"
)

// NonInteractiveSetup is used to setup the cothority node for unlynx in a non-interactive way (and without error checks)
func NonInteractiveSetup(c *cli.Context) error {

	// cli arguments
	serverBindingStr := c.String("serverBinding")
	description := c.String("description")
	privateTomlPath := c.String("privateTomlPath")
	publicTomlPath := c.String("publicTomlPath")

	if serverBindingStr == "" || description == "" || privateTomlPath == "" || publicTomlPath == "" {
		err := errors.New("arguments not OK")
		log.Error(err)
		return cli.NewExitError(err, 3)
	}

	kp := key.NewKeyPair(libunlynx.SuiTe)

	privStr, _ := encoding.ScalarToStringHex(libunlynx.SuiTe, kp.Private)
	pubStr, _ := encoding.PointToStringHex(libunlynx.SuiTe, kp.Public)
	public, _ := encoding.StringHexToPoint(libunlynx.SuiTe, pubStr)

	serverBinding := network.NewTLSAddress(serverBindingStr)
	conf := &app.CothorityConfig{
		Public:      pubStr,
		Private:     privStr,
		Address:     serverBinding,
		Description: description,
	}

	server := app.NewServerToml(libunlynx.SuiTe, public, serverBinding, conf.Description)
	group := app.NewGroupToml(server)

	conf.Save(privateTomlPath)
	group.Save(publicTomlPath)

	return nil
}
