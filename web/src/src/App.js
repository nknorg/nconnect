import React from 'react';
import QRCode from 'qrcode';
import { withTranslation, Trans } from 'react-i18next';
import { Button, Container, MenuItem, List, ListItem, ListItemText, Tab, TextField, Tooltip, Select, Grid } from '@material-ui/core';
import { TabContext, TabList, TabPanel } from '@material-ui/lab';

import i18n, { resources as languages } from './i18n';
import * as rpc from './rpc';

import './App.css';

class HoverQRCode extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      qrCode: '',
    };
  }

  componentDidMount() {
    QRCode.toDataURL(this.props.rawData).then(qrCode => {
      this.setState({ qrCode });
    }).catch(console.error);
  }

  render() {
    return (
      <Tooltip title={<img src={this.state.qrCode} alt="QR Code" />} >
        <img src="/static/media/qr_code.png" alt="QR Code" style={{height: '24px', verticalAlign: 'middle'}} />
      </Tooltip>
    );
  }
}

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      activeTab: '0',
      adminTokenStr: '',
      adminTokenQRCode: '',
      acceptAddrs: '',
      adminAddrs: '',
      addr: '',
      localIP: [],
      language: '',
      inPrice: [],
      outPrice: [],
      balance: '',
    };
    for (let i = 0; i < i18n.languages.length; i++) {
      if (languages[i18n.languages[i]]) {
        this.state.language = i18n.languages[i];
        break;
      }
    }
    this.handleTabChange = this.handleTabChange.bind(this);
    this.handleAcceptAddrsChange = this.handleAcceptAddrsChange.bind(this);
    this.handleAdminAddrsChange = this.handleAdminAddrsChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.updateAdminToken = this.updateAdminToken.bind(this);
    this.handleLanguageChange = this.handleLanguageChange.bind(this);
  }

  handleTabChange(event, value) {
    this.setState({ activeTab: value });
    if (value === '4') {
      this.updateAdvancedInfo();
    }
  }

  handleAcceptAddrsChange(event) {
    this.setState({ acceptAddrs: event.target.value });
  }

  handleAdminAddrsChange(event) {
    this.setState({ adminAddrs: event.target.value });
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
      alert(this.props.t('save success'));
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

  updateAdvancedInfo() {
    this.updateAdminToken();

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
        inPrice: info.inPrice,
        outPrice: info.outPrice,
      });
    }).catch((e) => {
      console.error(e);
    });

    if (!this.state.balance) {
      rpc.getBalance().then((balance) => {
        this.setState({ balance });
      }).catch((e) => {
        console.error(e);
      });
    }
  }

  estimatedRemainingData() {
    if (!this.state.balance) {
      return null;
    }

    if (!(this.state.inPrice && this.state.inPrice.length) && !(this.state.outPrice && this.state.outPrice.length)) {
      return null;
    }

    let balance = parseFloat(this.state.balance);
    if (isNaN(balance)) {
      return null;
    }

    let averagePrice = 0;
    for (let i = 0; i < this.state.inPrice.length; i++) {
      averagePrice += parseFloat(this.state.inPrice[i]);
    }
    for (let i = 0; i < this.state.outPrice.length; i++) {
      averagePrice += parseFloat(this.state.outPrice[i]);
    }
    averagePrice /= (this.state.inPrice.length + this.state.outPrice.length);
    if (isNaN(averagePrice)) {
      return null;
    }

    if (averagePrice == 0) {
      return null;
    }

    let mb = balance / averagePrice;
    let gb = mb / 1024;
    if (gb > 1) {
      return gb.toFixed(1) + ' GB';
    }
    return mb.toFixed(0) + ' MB';
  }

  componentDidMount() {
    this.updateAdvancedInfo();
    setInterval(this.updateAdminToken, 5 * 60 * 1000);
  }

  render() {
    let remainingData = this.estimatedRemainingData();
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
            <Grid container justify="center" alignItems="center">
              <Grid item xs={12} sm={6}>
                <img src="/static/media/nkn_logo.png" alt="NKN logo" />
              </Grid>
              <Grid item xs={12} sm={6}>
                <div className="row">
                  { remainingData && ('Estimated Remaining Data: ' + remainingData) }
                </div>
              </Grid>
            </Grid>
          </div>

          <TabContext value={this.state.activeTab}>
            <TabList centered onChange={this.handleTabChange}>
              <Tab label={this.props.t('mobile tab')} value="0" />
              <Tab label={this.props.t('desktop tab')} value="1" />
              <Tab label={this.props.t('data plan tab')} value="2" />
              <Tab label={this.props.t('need help tab')} value="3" />
              <Tab label={this.props.t('advanced tab')} value="4" />
            </TabList>
            <TabPanel value="0">
              <div className="row">
                <img src={this.state.adminTokenQRCode} alt="QR Code" />
              </div>
              <List>
                <ListItem>
                  <ListItemText>
                    <Trans
                      i18nKey="download nMobile pro"
                      components={{
                        nMobileProLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nMobileProLink')} />,
                        QRCode: <HoverQRCode rawData={this.props.t('nMobileProLink')} />,
                      }}
                    />
                  </ListItemText>
                </ListItem>
                <ListItem>
                  <ListItemText>
                    {this.props.t('add device from mobile')}
                  </ListItemText>
                </ListItem>
                <ListItem>
                  <ListItemText>
                    {this.props.t('connect from mobile')}
                  </ListItemText>
                </ListItem>
                <ListItem>
                  <ListItemText>
                    <Trans
                      i18nKey="mobile guide"
                      components={{
                        guideLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('getStartedLink')} />,
                      }}
                    />
                  </ListItemText>
                </ListItem>
              </List>
            </TabPanel>
            <TabPanel value="1">
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="add device in mobile first"
                    components={{
                      nMobileProLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nMobileProLink')} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="add server from desktop"
                    components={{
                      nConnectClientDesktopLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nConnectClientDesktopLink')} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  {this.props.t('scan QR code to add server to desktop')}
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  {this.props.t('connect from desktop')}
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="desktop guide"
                    components={{
                      guideLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('getStartedLink')} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
            </TabPanel>
            <TabPanel value="2">
              <ListItem>
                <ListItemText>
                  {this.props.t('purchase method')}
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="purchase from mobile"
                    components={{
                      nMobileProLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nMobileProLink')} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="purchase from web"
                    components={{
                      paymentLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('paymentLink', {addr: addrToPubKey(this.state.addr), lng: this.state.language})} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
            </TabPanel>
            <TabPanel value="3">
              <ListItem>
                <ListItemText>
                  {this.props.t('need help method')}
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="create forum post"
                    components={{
                      forumLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('forumLink')} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="send email"
                    components={{
                      emailLink: <a href={'mailto:'+this.props.t('emailAddress')} />,
                      emailAddress: this.props.t('emailAddress'),
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="mobile customer service"
                    components={{
                      nMobileProLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nMobileProLink')} />,
                    }}
                  />
                </ListItemText>
              </ListItem>
            </TabPanel>
            <TabPanel value="4">
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
            </TabPanel>
          </TabContext>
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
