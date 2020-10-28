import React from 'react';
import QRCode from 'qrcode';
import { withTranslation, Trans } from 'react-i18next';
import { Button, Collapse, Container, TextField, MenuItem, Select } from '@material-ui/core';
import { ExpandLess, ExpandMore } from '@material-ui/icons';

import i18n, { resources as languages } from './i18n';
import * as rpc from './rpc';

import './App.css';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      adminTokenStr: '',
      adminTokenQRCode: '',
      acceptAddrs: '',
      adminAddrs: '',
      addr: '',
      localIP: [],
      language: '',
      showAdvanced: false,
    };
    for (let i = 0; i < i18n.languages.length; i++) {
      if (languages[i18n.languages[i]]) {
        this.state.language = i18n.languages[i];
        break;
      }
    }
    this.handleAcceptAddrsChange = this.handleAcceptAddrsChange.bind(this);
    this.handleAdminAddrsChange = this.handleAdminAddrsChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.updateAdminToken = this.updateAdminToken.bind(this);
    this.handleAdvancedChange = this.handleAdvancedChange.bind(this);
    this.handleLanguageChange = this.handleLanguageChange.bind(this);
  }

  handleAcceptAddrsChange(event) {
    this.setState({ acceptAddrs: event.target.value });
  }

  handleAdminAddrsChange(event) {
    this.setState({ adminAddrs: event.target.value });
  }

  handleAdvancedChange(event) {
    this.setState({ showAdvanced: !this.state.showAdvanced });
  }

  handleLanguageChange(event) {
    this.setState({ language: event.target.value });
    i18n.changeLanguage(event.target.value);
  }

  async handleSubmit(event) {
    event.preventDefault();
    try {
      let addrs = await rpc.setAddrs(strToAddrs(this.state.acceptAddrs), strToAddrs(this.state.adminAddrs));
      this.setState({
        acceptAddrs: addrsToStr(addrs.acceptAddrs),
        adminAddrs: addrsToStr(addrs.adminAddrs),
      });
      alert('Save success!');
    } catch (e) {
      console.error(e);
      alert(e);
    }
  }

  async updateAdminToken() {
    try {
      let adminToken = await rpc.getAdminToken();
      if (adminToken) {
        let adminTokenStr = JSON.stringify(adminToken);
        let adminTokenQRCode = await QRCode.toDataURL(adminTokenStr);
        this.setState({
          adminTokenStr,
          adminTokenQRCode,
        });
      }
    } catch (e) {
      console.error(e);
    }
  }

  componentDidMount() {
    this.updateAdminToken();
    setInterval(this.updateAdminToken, 5 * 60 * 1000);

    rpc.getAddrs().then((addrs) => {
      this.setState({
        acceptAddrs: addrsToStr(addrs.acceptAddrs),
        adminAddrs: addrsToStr(addrs.adminAddrs),
      });
    }).catch((e) => {
      console.error(e);
      alert(e);
    });

    rpc.getInfo().then((info) => {
      this.setState({
        addr: info.addr,
        localIP: info.localIP.ipv4,
      });
    }).catch((e) => {
      console.error(e);
    });
  }

  render() {
    return (
      <div className="App">
        <Container>
          <div className="language-switcher">
            <Select
              value={this.state.language}
              onChange={this.handleLanguageChange}
            >
              {
                Object.keys(languages).map((lang) => {
                  return (
                    <MenuItem key={lang} value={lang}>{i18n.getFixedT(lang)('language')}</MenuItem>
                  )
                })
              }
            </Select>
          </div>
          <div className="row">
            <img src="/static/media/nkn_logo.png" alt="NKN logo" />
          </div>
          <div className="row">
            <img src={this.state.adminTokenQRCode} alt="QR Code" />
          </div>
          <div className="row">
            <Trans
              i18nKey="get started"
              components={{
                getStartedLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('getStartedLink')} />,
              }}
            />
          </div>
          <div className="row">
            <Trans
              i18nKey="QR code description"
              components={{
                nMobileProLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nMobileProLink')} />,
              }}
            />
          </div>
          <div className="row">
            <Trans
              i18nKey="desktop client description"
              components={{
                nConnectClientDesktopLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nConnectClientDesktopLink')} />,
              }}
            />
          </div>
          <div className="row">
            <Trans
              i18nKey="purchase description"
              components={{
                nMobileProLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nMobileProLink')} />,
                paymentLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('paymentLink', {addr: addrToPubKey(this.state.addr), lng: this.state.language})} />,
              }}
            />
          </div>
          <div className="row">
            <Trans
              i18nKey="custom service description"
              components={{
                customServiceLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('customServiceLink')} />,
              }}
            />
          </div>
          <div className="row">
            <Button
              variant="outlined"
              color="primary"
              onClick={this.handleAdvancedChange}
              style={{width: '100%'}}
              >
              {this.state.showAdvanced ? <ExpandLess /> : <ExpandMore /> }
              {this.state.showAdvanced ? this.props.t('hide advanced') : this.props.t('show advanced')}
            </Button>
          </div>
          <Collapse in={this.state.showAdvanced}>
            <div className="advanced-row">
              <TextField
                disabled
                multiline
                label={this.props.t('local IP address')}
                value={this.state.localIP.join('\n')}
                style={{width: '100%'}}
                />
            </div>
            <div className="advanced-row">
              <TextField
                disabled
                multiline
                label={this.props.t('access key')}
                value={this.state.adminTokenStr}
                style={{width: '100%'}}
                />
            </div>
            <div className="advanced-row">
              <TextField
                multiline
                variant="filled"
                label={this.props.t('accept addresses')}
                value={this.state.acceptAddrs}
                onChange={this.handleAcceptAddrsChange}
                style={{width: '100%'}}
                />
              <TextField
                multiline
                variant="filled"
                label={this.props.t('admins')}
                value={this.state.adminAddrs}
                onChange={this.handleAdminAddrsChange}
                style={{width: '100%'}}
                />
            </div>
            <div className="advanced-row">
              <Button
                variant="contained"
                color="primary"
                onClick={this.handleSubmit}
                style={{width: '100%'}}
                >
                {this.props.t('save')}
              </Button>
            </div>
          </Collapse>
        </Container>
      </div>
    );
  }
}

function addrsToStr(addrs) {
  if (!addrs) {
    return '';
  }
  return addrs.join('\n');
}

function strToAddrs(str) {
  if (!str) {
    return [];
  }
  return str.split('\n').filter(s => s.length > 0);
}

function addrToPubKey(addr) {
  let s = addr.split('.');
  return s[s.length-1];
}

export default withTranslation()(App);
