import axios from 'axios';
import React from 'react';
import QRCode from 'qrcode';
import { withTranslation, Trans } from 'react-i18next';
import {
  Button, Container, MenuItem, List, ListItem, ListItemText, Tab, TextField,
  Tooltip, Select, Grid, Dialog, DialogTitle, DialogContent, RadioGroup,
  FormControlLabel, Radio, DialogActions,
} from '@material-ui/core';
import { TabContext, TabList, TabPanel } from '@material-ui/lab';
import { ArrowDropDown } from '@material-ui/icons';

import i18n, { resources as languages } from './i18n';
import * as rpc from './rpc';

import './App.css';

const tunaConfigChoicesAddr = '/static/tuna-config-choices.json';

class HoverQRCode extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      qrCode: '',
    };
  }

  componentDidMount() {
    QRCode.toDataURL(this.props.rawData).then(qrCode => {
      this.setState({qrCode});
    }).catch(console.error);
  }

  render() {
    return (
      <Tooltip title={<img src={this.state.qrCode} alt="QR Code"/>}>
        <img src="/static/media/qr_code.png" alt="QR Code" style={{height: '24px', verticalAlign: 'middle'}}/>
      </Tooltip>
    );
  }
}

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      initialized: false,
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
      tags: [],
      balance: '',
      tunaServiceName: '',
      tunaCountry: [],
      tunaConfigChoices: [],
      tunaConfigSelected: -1,
      isTunaConfigChoiceOpen: false,
      currentTunaConfig: -1,
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
    this.handleExportAccount = this.handleExportAccount.bind(this);
    this.handleImportAccount = this.handleImportAccount.bind(this);
    this.openTunaConfigChoice = this.openTunaConfigChoice.bind(this);
    this.handleTunaConfigChoiceCancel = this.handleTunaConfigChoiceCancel.bind(this);
    this.handleTunaConfigChoiceOK = this.handleTunaConfigChoiceOK.bind(this);
    this.handleTunaConfigChoiceChange = this.handleTunaConfigChoiceChange.bind(this);
    this.downloadLog = this.downloadLog.bind(this);
  }

  handleTabChange(event, value) {
    this.setState({activeTab: value});
    if (value === '4') {
      this.updateInfo();
    }
  }

  handleAcceptAddrsChange(event) {
    this.setState({acceptAddrs: event.target.value});
  }

  handleAdminAddrsChange(event) {
    this.setState({adminAddrs: event.target.value});
  }

  handleLanguageChange(event) {
    this.setState({language: event.target.value});
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
      window.alert(this.props.t('save success'));
    } catch (e) {
      console.error(e);
      window.alert(e);
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

  updateInfo() {
    this.updateAdminToken();

    rpc.getAddrs().then((addrs) => {
      this.setState({
        acceptAddrs: addrsToStr(addrs.acceptAddrs),
        adminAddrs: addrsToStr(addrs.adminAddrs),
      });
    }).catch((e) => {
      console.error(e);
      window.alert(e);
    });

    rpc.getInfo().then((info) => {
      let initialized = this.state.initialized;
      this.setState({
        initialized: true,
        addr: info.addr,
        localIP: info.localIP.ipv4,
        inPrice: info.inPrice,
        outPrice: info.outPrice,
        tags: info.tags,
        tunaServiceName: info.tunaServiceName || '',
        tunaCountry: info.tunaCountry || [],
      });
      if (!initialized) {
        axios.get(tunaConfigChoicesAddr + '?t=' + Date.now()).then((response) => {
          if (response.data && response.data.length) {
            this.setState({
              tunaConfigChoices: response.data,
            });
            for (let i = 0; i < this.state.tunaConfigChoices.length; i++) {
              if (this.state.tunaConfigChoices[i].config.serviceName === this.state.tunaServiceName) {
                if (JSON.stringify(this.state.tunaConfigChoices[i].config.country) === JSON.stringify(this.state.tunaCountry)) {
                  this.setState({
                    tunaConfigSelected: i,
                    currentTunaConfig: i,
                  });
                  break;
                }
              }
            }
            if (this.state.tunaConfigSelected < 0 && !this.state.tunaServiceName && !(this.state.tunaCountry && this.state.tunaCountry.length)) {
              this.openTunaConfigChoice();
            }
          }
        }).catch((e) => {
          if (e.response.status !== 404) {
            console.error(e);
          }
        });
      }
    }).catch((e) => {
      console.error(e);
    });

    if (!this.state.balance) {
      rpc.getBalance().then((balance) => {
        this.setState({balance});
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

    let maxPrice = 0;
    for (let i = 0; i < this.state.inPrice.length; i++) {
      let price = parseFloat(this.state.inPrice[i]);
      if (!isNaN(price) && price > maxPrice) {
        maxPrice = price;
      }
    }
    for (let i = 0; i < this.state.outPrice.length; i++) {
      let price = parseFloat(this.state.outPrice[i]);
      if (!isNaN(price) && price > maxPrice) {
        maxPrice = price;
      }
    }

    if (isNaN(maxPrice)) {
      return null;
    }

    if (maxPrice === 0) {
      return null;
    }

    let mb = balance / maxPrice;
    let gb = mb / 1024;
    if (gb > 1) {
      return gb.toFixed(1) + ' GB';
    }
    return mb.toFixed(0) + ' MB';
  }

  async handleExportAccount(event) {
    event.preventDefault();

    if (!window.confirm(this.props.t('exportConfirm'))) {
      return;
    }

    try {
      let seed = await rpc.getSeed();
      window.alert(this.props.t('exportSuccess', {seed}));
    } catch (e) {
      console.error(e);
      window.alert(e);
    }
  }

  async handleImportAccount(event) {
    event.preventDefault();

    if (!window.confirm(this.props.t('importConfirm'))) {
      return;
    }

    let currentSeedInput = window.prompt(this.props.t('importPromptCurrent'));
    if (!currentSeedInput) {
      return;
    }

    try {
      let currentSeed = await rpc.getSeed();
      if (currentSeed !== currentSeedInput.trim()) {
        window.alert(this.props.t('importWrongCurrent'));
        return;
      }
    } catch (e) {
      console.error(e);
      window.alert(e);
      return;
    }

    let newSeed = window.prompt(this.props.t('importPromptNew'));
    if (!newSeed) {
      return;
    }

    try {
      await rpc.setSeed(newSeed.trim());
      window.alert(this.props.t('importSuccess'));
    } catch (e) {
      console.error(e);
      window.alert(e);
    }
  }

  openTunaConfigChoice() {
    this.setState({
      isTunaConfigChoiceOpen: true,
    });
  }

  handleTunaConfigChoiceCancel() {
    this.setState({
      isTunaConfigChoiceOpen: false,
    });
  }

  async handleTunaConfigChoiceOK() {
    this.setState({
      isTunaConfigChoiceOpen: false,
      currentTunaConfig: this.state.tunaConfigSelected,
    });
    if (this.state.tunaConfigSelected >= 0) {
      try {
        await rpc.setTunaConfig(this.state.tunaConfigChoices[this.state.tunaConfigSelected].config);
        window.alert(this.props.t('setTunaConfigSuccess'));
      } catch (e) {
        console.error(e);
        window.alert(e);
      }
    }
  }

  handleTunaConfigChoiceChange(event) {
    this.setState({
      tunaConfigSelected: event.target.value,
    });
  }

  async downloadLog(event) {
    event.preventDefault();

    try {
      let log = await rpc.getLog();
      if (log && log.length) {
        downloadFile('nConnect.log', log);
      } else {
        window.alert(this.props.t('no log available'))
      }
    } catch (e) {
      console.error(e);
      window.alert(e);
    }
  }

  componentDidMount() {
    this.updateInfo();
    setInterval(this.updateAdminToken, 5 * 60 * 1000);
  }

  render() {
    let remainingData = this.estimatedRemainingData();
    let paymentAdditionalParams = '';
    if (this.state.tags && this.state.tags.length) {
      for (let i = 0; i < this.state.tags.length; i++) {
        paymentAdditionalParams += '&tag=' + this.state.tags[i];
      }
    }
    let currentServerRegionName = '';
    if (this.state.tunaConfigChoices && this.state.tunaConfigChoices.length) {
      if (this.state.currentTunaConfig >= 0) {
        let item = this.state.tunaConfigChoices[this.state.currentTunaConfig];
        currentServerRegionName = this.props.t(item.textId) || item.textId;
      } else {
        currentServerRegionName = this.props.t('customized');
      }
    }
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
              <Grid item xs={12} sm={6} className="text-left">
                <img src="/static/media/nkn_logo.png" alt="NKN logo"/>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Grid container
                      direction="row"
                      justify="flex-end"
                      alignItems="center">
                  <Grid item xs>
                    <div className="row text-right">
                      {
                        currentServerRegionName && (
                          <span>
                        {this.props.t('currentServerRegion') + ': '}
                            <span onClick={this.openTunaConfigChoice} style={{cursor: 'pointer'}}>
                          {currentServerRegionName}
                              <ArrowDropDown style={{verticalAlign: 'middle'}}/>
                        </span>
                      </span>
                        )
                      }
                    </div>
                    <div className="row text-right">
                      {remainingData && (this.props.t('estimatedRemainingData') + ': ' + remainingData)}
                    </div>
                  </Grid>
                  <Grid item xs={4} className="text-right">
                    <Button variant="contained" color="primary" target="_blank" href={this.props.t('paymentLink', {
                      addr: addrToPubKey(this.state.addr),
                      lng: this.state.language,
                      additionalParams: paymentAdditionalParams
                    })}>
                      {this.props.t('data plan tab')}
                    </Button>
                  </Grid>
                </Grid>
              </Grid>
            </Grid>
          </div>

          <TabContext value={this.state.activeTab}>
            <TabList onChange={this.handleTabChange} className="bottom-line">
              <Tab label={this.props.t('mobile tab')} value="0"/>
              <Tab label={this.props.t('desktop tab')} value="1"/>
              <Tab label={this.props.t('advanced tab')} value="2"/>
              <Tab className="margin-left-auto" label={this.props.t('need help tab')} value="3"/>
            </TabList>
            <TabPanel value="0">
              <Grid container spacing={10} justify="center" alignItems="center">
                <Grid item xs={12} sm={6} className="text-right">
                  <img src={this.state.adminTokenQRCode} alt="QR Code"/>
                </Grid>
                <Grid item xs={12} sm={6} className="text-left">
                  <div>{this.props.t('export tip')}</div>
                </Grid>
              </Grid>
              <List>
                <ListItem>
                  <ListItemText>
                    <Trans
                      i18nKey="download nConnect"
                      components={{
                        nConnectLink: <a target="_blank" rel="noopener noreferrer"
                                         href={this.props.t('nConnectLink')}/>,
                        QRCode: <HoverQRCode rawData={this.props.t('nConnectLink')}/>,
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
                        guideLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('getStartedLink')}/>,
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
                      nConnectLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nConnectLink')}/>,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="add server from desktop"
                    components={{
                      nConnectClientDesktopLink: <a target="_blank" rel="noopener noreferrer"
                                                    href={this.props.t('nConnectClientDesktopLink')}/>,
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
                      guideLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('getStartedLink')}/>,
                    }}
                  />
                </ListItemText>
              </ListItem>
            </TabPanel>
            <TabPanel value="2">
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
              <div className="advanced-row">
                <div style={{paddingBottom: '5px'}}>
                  <Button
                    variant="contained"
                    onClick={this.handleExportAccount}
                    style={{width: '100%'}}
                  >
                    {this.props.t('export account')}
                  </Button>
                </div>
                <div style={{paddingTop: '5px'}}>
                  <Button
                    variant="contained"
                    onClick={this.handleImportAccount}
                    style={{width: '100%'}}
                  >
                    {this.props.t('import account')}
                  </Button>
                </div>
                <div className="advanced-row">
                  <Button
                    variant="contained"
                    onClick={this.downloadLog}
                    style={{width: '100%'}}
                  >
                    {this.props.t('download log')}
                  </Button>
                </div>
              </div>
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
                      forumLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('forumLink')}/>,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="Q&A"
                    components={{
                      'QALink': <a target="_blank" rel="noopener noreferrer" href={this.props.t('QALink')}/>,
                    }}
                  />
                </ListItemText>
              </ListItem>
              <ListItem>
                <ListItemText>
                  <Trans
                    i18nKey="send email"
                    components={{
                      emailLink: <a href={'mailto:' + this.props.t('emailAddress')}/>,
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
                      nConnectLink: <a target="_blank" rel="noopener noreferrer" href={this.props.t('nConnectLink')}/>,
                    }}
                  />
                </ListItemText>
              </ListItem>
            </TabPanel>
          </TabContext>
          <Dialog
            disableBackdropClick
            disableEscapeKeyDown
            keepMounted
            maxWidth="xs"
            open={this.state.isTunaConfigChoiceOpen}
          >
            <DialogTitle>{this.props.t('tunaConfigChoiceTitle')}</DialogTitle>
            <DialogContent dividers>
              <RadioGroup
                value={`${this.state.tunaConfigSelected}`}
                onChange={this.handleTunaConfigChoiceChange}
              >
                {
                  this.state.tunaConfigChoices.map((item, index) => (
                    <FormControlLabel
                      value={`${index}`}
                      key={`${index}`}
                      control={<Radio/>}
                      label={this.props.t(item.textId) || item.textId}
                    />
                  ))
                }
              </RadioGroup>
            </DialogContent>
            <DialogActions>
              <Button autoFocus onClick={this.handleTunaConfigChoiceCancel} color="primary">
                {this.props.t('cancel')}
              </Button>
              <Button onClick={this.handleTunaConfigChoiceOK} color="primary">
                {this.props.t('save')}
              </Button>
            </DialogActions>
          </Dialog>
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
  return s[s.length - 1];
}

function downloadFile(filename, content) {
  var element = document.createElement('a');
  element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(content));
  element.setAttribute('download', filename);
  element.style.display = 'none';
  document.body.appendChild(element);
  element.click();
  document.body.removeChild(element);
}

export default withTranslation()(App);
