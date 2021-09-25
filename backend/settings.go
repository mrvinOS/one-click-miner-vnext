package backend

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/tidwall/buntdb"
	"github.com/vertcoin-project/one-click-miner-vnext/logging"
	"github.com/vertcoin-project/one-click-miner-vnext/networks"
	"github.com/vertcoin-project/one-click-miner-vnext/pools"
	"github.com/vertcoin-project/one-click-miner-vnext/tracking"
	"github.com/vertcoin-project/one-click-miner-vnext/util"
)

func (m *Backend) getSetting(name string) bool {
	setting := "0"
	m.settings.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(name)
		setting = v
		return err
	})
	return setting == "1"
}

func (m *Backend) setSetting(name string, value bool) {
	setting := "0"
	if value {
		setting = "1"
	}
	m.settings.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(name, setting, nil)
		return err
	})
}

func (m *Backend) setIntSetting(name string, value int) {
	setting := fmt.Sprintf("%d", value)
	m.settings.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(name, setting, nil)
		return err
	})
}

func (m *Backend) getIntSetting(name string) int {
	setting := "0"
	m.settings.View(func(tx *buntdb.Tx) error {
		v, err := tx.Get(name)
		setting = v
		return err
	})
	i, _ := strconv.Atoi(setting)
	return i
}

func (m *Backend) GetPool() int {
	pool := m.getIntSetting("pool")
	if pool == 0 {
		if m.GetTestnet() {
			return 2 // Default P2Pool on testnet
		}
		// Default to a random pool
		rand.Seed(time.Now().UnixNano())
		pools := pools.GetPools(m.Address(), m.GetTestnet())
		pool := pools[rand.Intn(len(pools))].GetID()
		// Save this setting immediately so that we don't get
		// a different random pool in future calls to GetPool().
		m.setIntSetting("pool", pool)
		return pool
	}
	return pool
}

func (m *Backend) SetPool(pool int) {
	if m.GetPool() != pool {
		m.setIntSetting("pool", pool)
		m.ResetPool()
		logging.Infof("Calling WalletInitialized\n")
		m.WalletInitialized()
		logging.Infof("Done!")
	}
}

func (m *Backend) SetEnableIntegrated(enabled bool) {
	m.setSetting("enableIntegrated", enabled)
}

func (m *Backend) GetEnableIntegrated() bool {
	return m.getSetting("enableIntegrated")
}

type PoolChoice struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (m *Backend) GetPools() []PoolChoice {
	pc := make([]PoolChoice, 0)
	for _, p := range pools.GetPools(m.Address(), m.GetTestnet()) {
		pc = append(pc, PoolChoice{
			ID:   p.GetID(),
			Name: fmt.Sprintf("%s (%0.2f%% fee)", p.GetName(), p.GetFee()),
		})
	}
	return pc
}

func (m *Backend) GetTestnet() bool {
	return m.getSetting("testnet")
}

func (m *Backend) SetTestnet(newTestnet bool) {
	if m.GetTestnet() != newTestnet {
		logging.Infof("Setting testnet to [%b]\n", newTestnet)
		m.setSetting("testnet", newTestnet)

		logging.Infof("Setting network to testnet=%b\n", newTestnet)
		networks.SetNetwork(newTestnet)

		logging.Infof("Calling WalletInitialized\n")
		m.WalletInitialized()
		logging.Infof("Done!")
	}
}

func (m *Backend) GetSkipVerthashExtendedVerify() bool {
	return false // Verification is default - return false
	//return m.getSetting("skipverthashverify")
}

func (m *Backend) SetSkipVerthashExtendedVerify(newVerthashVerify bool) {
	logging.Infof("Setting skip verthash verify to [%b]\n", newVerthashVerify)
	m.setSetting("skipverthashverify", newVerthashVerify)
}

func (m *Backend) GetClosedSource() bool {
	return false // No closed source Verthash miners - return false
	//return m.getSetting("closedsource")
}

func (m *Backend) SetClosedSource(newClosedSource bool) {
	logging.Infof("Setting closed source to [%b]\n", newClosedSource)
	m.setSetting("closedsource", newClosedSource)
}

func (m *Backend) GetDebugging() bool {
	return m.getSetting("debugging")
}

func (m *Backend) SetDebugging(newDebugging bool) {
	logging.Infof("Setting debugging to [%b]\n", newDebugging)
	m.setSetting("debugging", newDebugging)
}

func (m *Backend) GetAutoStart() bool {
	return util.GetAutoStart()
}

func (m *Backend) SetAutoStart(newAutoStart bool) {
	util.SetAutoStart(newAutoStart)
}

func (m *Backend) GetVersion() string {
	return tracking.GetVersion()
}

func (m *Backend) PrerequisiteProxyLoop() {
	for pi := range m.prerequisiteInstall {
		send := "0"
		if pi {
			send = "1"
		}
		m.runtime.Events.Emit("prerequisiteInstall", send)
	}
}

func (m *Backend) SaveBg(bg string) {

	bgPath := filepath.Join(util.DataDirectory(), "background")
	f, err := os.OpenFile(bgPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil { logging.Errorf("Open error: %s\n", err) }
	f.WriteString(bg)
	f.WriteString("\n")
	err = f.Close();
	if err != nil { logging.Errorf("Close error: %s\n", err) }
}

func (m *Backend) ReadBg() string {

	bgPath := filepath.Join(util.DataDirectory(), "background")
	if _, err := os.Stat(bgPath); err == nil {
		f, err := os.Open(bgPath)
		if err != nil { return "" }

		reader := bufio.NewReader(f)
		var line string
		line, err = reader.ReadString('\n')
		err = f.Close();
		if err != nil { logging.Errorf("Close error: %s", err) }

		//logging.Infof(line)
		return line
	}

	return ""
}

func (m *Backend) SaveTheme(theme string) {

	themePath := filepath.Join(util.DataDirectory(), "theme")
	f, err := os.OpenFile(themePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil { logging.Errorf("Open error: %s\n", err) }
	f.WriteString(theme)
	f.WriteString("\n")
	err = f.Close();
	if err != nil { logging.Errorf("Close error: %s\n", err) }
}

func (m *Backend) ReadTheme() string {

	themePath := filepath.Join(util.DataDirectory(), "theme")
	if _, err := os.Stat(themePath); err == nil {
		f, err := os.Open(themePath)
		if err != nil { return "" }

		reader := bufio.NewReader(f)
		var line string
		line, err = reader.ReadString('\n')
		err = f.Close();
		if err != nil { logging.Errorf("Close error: %s", err) }

		//logging.Infof(line)
		return line
	}

	return ""
}
