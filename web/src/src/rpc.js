import axios from 'axios';

import * as util from './util';

const rpcAddr = '/rpc/admin';

const methods = {
  getAdminToken: { method: 'getAdminToken' },
  getAddrs: { method: 'getAddrs' },
  setAddrs: { method: 'setAddrs' },
  addAddrs: { method: 'addAddrs' },
  removeAddrs: { method: 'removeAddrs' },
  getLocalIP: { method: 'getLocalIP' },
  getInfo: { method: 'getInfo' },
}

var rpc = {};
for (let method in methods) {
  if (methods.hasOwnProperty(method)) {
    rpc[method] = (addr, params) => {
      params = util.assignDefined({}, methods[method].defaultParams, params)
      return rpcCall(addr, methods[method].method, params);
    }
  }
}

async function rpcCall(addr, method, params = {}) {
  let headers;
  try {
    headers = await window.rpcHeaders;
  } catch (e) {
    console.error('Await rpc headers error:', e);
  }

  let response = await axios({
    url: addr,
    method: 'POST',
    timeout: 10000,
    headers,
    data: {
      id: 'nConnect-web',
      jsonrpc: '2.0',
      method: method,
      params: params,
    },
  });

  let data = response.data;

  if (data.error) {
    throw data.error;
  }

  if (data.result !== undefined) {
    return data.result;
  }

  throw new Error('rpc response contains no result or error field');
}

export async function getAdminToken() {
  return rpc.getAdminToken(rpcAddr);
}

export async function getAddrs() {
  return rpc.getAddrs(rpcAddr);
}

export async function setAddrs(acceptAddrs, adminAddrs) {
  let params = {};
  if (acceptAddrs) {
    params.acceptAddrs = acceptAddrs;
  }
  if (adminAddrs) {
    params.adminAddrs = adminAddrs;
  }
  return rpc.setAddrs(rpcAddr, params);
}

export async function addAddrs(acceptAddrs, adminAddrs) {
  let params = {};
  if (acceptAddrs) {
    params.acceptAddrs = acceptAddrs;
  }
  if (adminAddrs) {
    params.adminAddrs = adminAddrs;
  }
  return rpc.addAddrs(rpcAddr, params);
}

export async function removeAddrs(acceptAddrs, adminAddrs) {
  let params = {};
  if (acceptAddrs) {
    params.acceptAddrs = acceptAddrs;
  }
  if (adminAddrs) {
    params.adminAddrs = adminAddrs;
  }
  return rpc.removeAddrs(rpcAddr, params);
}

export async function getLocalIP() {
  return rpc.getLocalIP(rpcAddr);
}

export async function getInfo() {
  return rpc.getInfo(rpcAddr);
}
