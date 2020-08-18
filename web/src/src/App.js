import React from 'react';
import QRCode from 'qrcode';
import Button from '@material-ui/core/Button';
import Container from '@material-ui/core/Container';
import TextField from '@material-ui/core/TextField';

import * as rpc from './rpc';

import './App.css';

class App extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      adminTokenQRCode: null,
      acceptAddrs: '',
      adminAddrs: '',
      localIP: [],
    };
    this.handleAcceptAddrsChange = this.handleAcceptAddrsChange.bind(this);
    this.handleAdminAddrsChange = this.handleAdminAddrsChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleAcceptAddrsChange(event) {
    this.setState({ acceptAddrs: event.target.value });
  }

  handleAdminAddrsChange(event) {
    this.setState({ adminAddrs: event.target.value });
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

  componentDidMount() {
    rpc.getAdminToken().then((adminToken) => {
      if (adminToken) {
        return QRCode.toDataURL(JSON.stringify(adminToken))
      }
    }).then((url) => {
      if (url) {
        this.setState({
          adminTokenQRCode: url,
        });
      }
    }).catch((e) => {
      console.error(e);
      alert(e);
    });

    rpc.getAddrs().then((addrs) => {
      this.setState({
        acceptAddrs: addrsToStr(addrs.acceptAddrs),
        adminAddrs: addrsToStr(addrs.adminAddrs),
      });
    }).catch((e) => {
      console.error(e);
      alert(e);
    });

    rpc.getLocalIP().then((localIP) => {
      this.setState({
        localIP: localIP.ipv4,
      });
    }).catch((e) => {
      console.error(e);
    });
  }

  render() {
    return (
      <div className="App">
        <Container>
          <div>
            <img
              src={this.state.adminTokenQRCode}
              />
          </div>
          <div>
            <TextField
              disabled
              multiline
              label="Local IP address"
              value={this.state.localIP.join('\n')}
              />
          </div>
          <div>
            <TextField
              multiline
              variant="filled"
              label="Accept addresses"
              value={this.state.acceptAddrs}
              onChange={this.handleAcceptAddrsChange}
              style={{width: '100%'}}
              />
          </div>
          <div>
            <TextField
              multiline
              variant="filled"
              label="Admin addresses"
              value={this.state.adminAddrs}
              onChange={this.handleAdminAddrsChange}
              style={{width: '100%'}}
              />
          </div>
          <div>
            <Button
              variant="contained"
              color="primary"
              onClick={this.handleSubmit}
              >
              Save
            </Button>
          </div>
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

export default App;
