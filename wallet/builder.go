package wallet

import (
	"errors"

	"github.com/dabankio/wallet-core/bip44"
)

type WalletOptions struct {
	options []WalletOption
}

func (opts *WalletOptions) Add(opt WalletOption) {
	opts.options = append(opts.options, opt)
}

func (opts *WalletOptions) getOptions() []WalletOption {
	return opts.options
}

type WalletOption interface {
	Visit(*Wallet) error
}

// Clone makes a copy of existing Wallet instance, with original attributes override by the given options
func (c Wallet) Clone(options *WalletOptions) (wallet *Wallet, err error) {
	cloned := c
	for _, opt := range options.getOptions() {
		err = opt.Visit(&cloned)
		if err != nil {
			return nil, err
		}
	}
	//TODO verify wallet
	return &cloned, nil
}

type shareAccountWithParentChainOpt struct {
	shareAccountWithParentChain bool
}

func (opt *shareAccountWithParentChainOpt) Visit(wallet *Wallet) error {
	wallet.ShareAccountWithParentChain = opt.shareAccountWithParentChain
	return nil
}

func WithShareAccountWithParentChain(shareAccountWithParentChain bool) WalletOption {
	return &shareAccountWithParentChainOpt{
		shareAccountWithParentChain: shareAccountWithParentChain,
	}
}

type pathFormatOpt struct {
	pathFormat string
}

func (f *pathFormatOpt) Visit(wallet *Wallet) error {
	wallet.path = f.pathFormat
	return nil
}

func WithPathFormat(pathFormat string) WalletOption {
	return &pathFormatOpt{
		pathFormat: pathFormat,
	}
}

type flagOpt struct {
	flag string
}

func (f flagOpt) Visit(wallet *Wallet) error {
	wallet.flags[f.flag] = struct{}{}
	return nil
}

// WithFlag 该选项添加特殊配置flag, flag参考 FlagXXX 常量
func WithFlag(flag string) WalletOption {
	return flagOpt{flag}
}

type passwordOpt struct {
	password string
}

func (f *passwordOpt) Visit(wallet *Wallet) error {
	wallet.password = f.password
	return nil
}

func WithPassword(password string) WalletOption {
	return &passwordOpt{password: password}
}

// NewWalletBuilder normal builder pattern, not so good in golang
func NewWalletBuilder() *WalletBuilder {
	return &WalletBuilder{}
}

type WalletBuilder struct {
	mnemonic                    string
	shareAccountWithParentChain bool
	seed                        []byte
	testNet                     bool
	password                    string
	pathFormat                  string
}

func (wb *WalletBuilder) SetMnemonic(mnemonic string) *WalletBuilder {
	wb.mnemonic = mnemonic
	return wb
}

func (wb *WalletBuilder) SetTestNet(testNet bool) *WalletBuilder {
	wb.testNet = testNet
	return wb
}

func (wb *WalletBuilder) SetPassword(password string) *WalletBuilder {
	wb.password = password
	return wb
}

func (wb *WalletBuilder) SetShareAccountWithParentChain(shareAccountWithParentChain bool) *WalletBuilder {
	wb.shareAccountWithParentChain = shareAccountWithParentChain
	return wb
}

func (wb *WalletBuilder) SetUseShortestPath(useShortestPath bool) *WalletBuilder {
	var pathFormat string
	if useShortestPath {
		pathFormat = bip44.PathFormat
	} else {
		pathFormat = bip44.FullPathFormat
	}
	wb.pathFormat = pathFormat
	return wb
}

func (wb *WalletBuilder) Wallet() (wallet *Wallet, err error) {
	if wb.mnemonic == "" {
		return nil, errors.New("mnemonic should not be empty")
	}
	wallet, err = NewHDWalletFromMnemonic(wb.mnemonic, wb.testNet)
	if err != nil {
		return nil, err
	}
	wallet.path = wb.pathFormat
	wallet.ShareAccountWithParentChain = wb.shareAccountWithParentChain
	wallet.password = wb.password
	//TODO verify wallet
	return
}

// BuildWallet create a Wallet instance with fixed args (here is mnemonic and testNet) and other options
func BuildWalletFromMnemonic(mnemonic string, testNet bool, options *WalletOptions) (wallet *Wallet, err error) {
	wallet, err = NewHDWalletFromMnemonic(mnemonic, testNet)
	if err != nil {
		return
	}
	for _, opt := range options.getOptions() {
		err = opt.Visit(wallet)
		if err != nil {
			return
		}
	}
	//TODO verify wallet
	return
}

// TODO not implemented
// BuildWallet create a Wallet instance with fixed args (here is privateKey and testNet) and other options
func BuildWalletFromPrivateKey(privateKey string, testNet bool, options WalletOptions) (wallet *Wallet, err error) {
	panic("implement me")
}
