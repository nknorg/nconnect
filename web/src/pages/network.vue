<template>
    <v-container>
      <v-parallax :src="require('~/static/img/grid.png')" height="800"
                  style="position: absolute; top: 0; max-width: 1600px;min-width: 640px;width: 80%">
        <v-img :src="require('~/static/img/point.png')" width="40" height="40" contain
               style="position: absolute; right: 38%; bottom: 130px;"/>
        <v-img :src="require('~/static/img/point.png')" width="67" height="67" contain
               style="position: absolute; right: 20%; bottom: 80px;"/>
      </v-parallax>
      <v-container class="container">
        <header>
          <v-row class="" align="center">
            <v-col>
              <img src="~/static/img/nkn-logo.png" alt="">
            </v-col>
  
            <v-col class="ml-auto" cols="auto">
              <span>Language</span>
              <span style="max-width: 200px; display: inline-flex;">
                <v-select class="nav-item select-language" @change="onChangeSwitchLanguage(lang)"
                          v-model="lang" :items="availableLocales" item-value="code"
                          item-text="name" label="Language" solo dark style="width: 200px;">
                </v-select>
              </span>
            </v-col>
          </v-row>
        </header>

        <v-row>
          <v-col cols="20">

            <div class="bg-linear-1 pa-4 mb-4">
              <div class="w-100 d-flex justify-space-between mb-4">
                <h3>Network Manager</h3>
              </div>
              <div class="d-flex flex-wrap justify-space-between mt-5">
                <div>Address: {{ managerAddress }}</div>
                <div>NKN Balance: {{ Number(managerBalance).toFixed(2) }}</div>
              </div>
            </div>

            <div class="bg-linear-1 pa-4 mb-4">
              <div class="w-100 d-flex justify-space-between mb-4">
                <h3>Network Configuration</h3>
                <div class="expand" @click="showConfig=!showConfig">▼</div>
              </div>
              <div v-if="showConfig">
                <v-text-field v-model="networkData.networkInfo.domain" clearable hide-details="auto" label="Network domain" class="mx-6"></v-text-field>
                <div class="d-flex flex-wrap mt-5">
                  <v-text-field v-model="networkData.ipStart" clearable hide-details="auto" label="IP Start" class="mx-6"></v-text-field>
                  <v-text-field v-model="networkData.ipEnd" clearable hide-details="auto" label="IP End" class="mx-6"></v-text-field>
                  <v-text-field v-model="networkData.netmask" clearable hide-details="auto" label="Network mask" class="mx-6"></v-text-field>
                </div>
                <div class="d-flex flex-wrap mt-5">
                  <v-text-field v-model="networkData.networkInfo.gateway" clearable hide-details="auto" label="Gateway" class="mx-6"></v-text-field>
                  <v-text-field v-model="networkData.networkInfo.dns" clearable hide-details="auto" label="DNS" class="mx-6"></v-text-field>
                </div>
                <div class="d-flex justify-center mt-12">
                  <v-btn color="primary" style="width:40%;" @click="setNetworkConfig">Submit</v-btn>
                </div>
              </div>
            </div>
            
            <div class="bg-linear-1 pa-4 mb-4" style="overflow-x:scroll;">
              <div class="w-100 d-flex justify-space-between mb-4">
                <h3>Waiting for Authorization</h3>
                <v-btn color="primary" @click="getNetworkConfig">Refresh</v-btn>
              </div>
              <table>
                  <tr>
                    <th>Name</th>
                    <th>Address</th>
                    <th>Accept</th>
                    <th>Reject</th>
                  </tr>
                  <tr v-for="item in networkData.waiting" v-bind:key="item.address">
                    <td>{{ item.name }}</td>
                    <td >{{ item.address }}</td>
                    <td>
                      <v-btn class="ma-2" color="primary" @click="authorizeMember(item.address)"> Accept </v-btn> 
                    </td>
                    <td><v-btn color="secondary" @click="deleteWaiting(item.address)">Reject</v-btn></td>
                  </tr>
              </table>
            </div>
  
            <div class="bg-linear-1 pa-4 mb-4" style="overflow-x:scroll;">
              <div class="w-100 d-flex justify-space-between mb-4">
                <h3>Network Members</h3>
                <div class="expand" @click="showMembers=!showMembers">▼</div>
              </div>
              <table v-if="showMembers">
                  <tr>
                    <th>Name</th>
                    <th>IP</th>
                    <th>LastSeen</th>
                    <th>Server</th>
                    <th>Balance</th>
                    <th>Address</th>
                    <th>Accept</th>
                    <th>Send Token</th>
                    <th>Ping</th>
                    <th>Remove</th>
                  </tr>
                  <tr v-for="item in networkData.member" v-bind:key="item.address">
                    <td>{{ item.name }}</td>
                    <td>{{ item.ip }}</td>
                    <td>{{ item.lastSeen.substring(2,19).replace("T", " ") }}</td>
                    <td :style="item.server?'background:green':''">{{ item.server? 'Yes' : 'No' }}</td>
                    <td :style="item.server?(item.balance>0.1?'background:green':'background:orange'):''">
                      {{ Number(item.balance)>0 ? Number(item.balance).toFixed(2): item.balance }}</td>
                    <td style="width:260px; word-break: break-all;">{{ item.address }}</td>
                    <td><v-btn color="primary" @click="toSetAcceptAddress(item.address)">Set</v-btn></td>
                    <td><v-btn v-if="item.server" color="primary" @click="sendToken(item.address)">Send</v-btn></td>
                    <td><v-btn color="primary" @click="nknPing(item.address)">Ping</v-btn></td>
                    <td><v-btn @click="removeMember(item.address)" color="secondary">Remove</v-btn></td>
                  </tr>
              </table>
              
            </div>
  
          </v-col>
        </v-row>

        <v-dialog v-model="dialog" width="auto">
          <v-card>
            <v-card-text>
              <h3 style="margin-top:1rem;">Set Accept Address</h3>
              <h5 class="mt-5">Accept All</h5>
              <v-checkbox v-model="allAddress" label="All members" color="red" value="allMembers" @click="clickAllAddress"></v-checkbox>
              <v-divider></v-divider>
              <div>
                <table >
                  <tr>
                    <th>Name</th>
                    <th>IP</th> 
                    <th>Address</th> 
                  </tr>
                  <tr v-for="item in networkData.member" v-bind:key="item.address" v-if="item.address != setAddress">
                    <td><v-checkbox v-model="acceptAddress" :label="item.name" color="red" :value="item.address" @click="selectAddress(item.address)" ></v-checkbox></td>
                    <td>{{ item.ip }}</td>
                    <td>{{ item.address }}</td>

                  </tr>
                </table>
              </div>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="secondary" @click="dialog = false" > Cancel </v-btn>
              <v-btn color="primary" @click="setAcceptAddress" > Save </v-btn>
            </v-card-actions>

          </v-card>
        </v-dialog>

        <v-dialog v-model="showResponse" transition="dialog-top-transition" width="400" >
          <v-card class="mx-auto">
            <div class="blankline"></div>
            <v-card-text>
              <h3 class="mt-4 mb-4">Tips</h3>
              <h3>{{response}}</h3>
            </v-card-text>

            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="showResponse = false" > Close </v-btn>
            </v-card-actions>

          </v-card>
        </v-dialog>

        <v-dialog v-model="confirmDialog" width="300" >
          <v-card>
            <v-card-text> {{ confirmTitle }} </v-card-text>
            <v-card-actions>
              <v-spacer></v-spacer>
              <v-btn color="secondary" @click="(confirmDialog=false); (confirm = false)">Cancel</v-btn>
              <v-btn color="primary" @click="(confirmDialog=false); (confirm = true)">Confirm</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>

        <v-dialog v-model="sendDialog" width="300">
          <v-card>
            <div class="blankline"></div>
            <v-card-text> Send NKN token to member </v-card-text>
            <div class="blankline"></div>
            <v-text-field v-model="amount" clearable hide-details="auto" label="NKN Amount" class="mx-6"></v-text-field>
            <div class="blankline"></div>
            <v-card-actions>
              <v-btn color="secondary" @click="(sendDialog=false); (confirm = false)">Cancel</v-btn>
              <v-spacer></v-spacer>
              <v-btn color="primary" @click="(sendDialog=false); (confirm = true)">Confirm</v-btn>
            </v-card-actions>
          </v-card>
        </v-dialog>

      </v-container>
    </v-container>
  </template>
  
  <script>
  import * as rpc from '../assets/network_rpc'
  import Cookies from "js-cookie"
  
  export default {
    name: 'network',
    computed: {
      availableLocales() {
        return this.$i18n.locales
      }
    },
    data() {
      return {
        lang: this.$i18n.locale,
        managerAddress: '',
        managerBalance: '',
  
        showConfig: false,
        showMembers: true,

        networkData: {networkInfo: {}, member:[{name: "bill", ip: "10.0.86.3", address: "aaaabbbbccccdddd", lastSeen:"2023-09-11 13:00:00"}]},
        waitingCheck: [],
        memberCheck: [],

        setAddress: '',
        acceptAddress: [],
        allAddress: [],
        dialog: false,

        showResponse: false,
        response: '',
        
        confirmDialog: false,
        confirmTitle: '',
        confirm: false,

        sendDialog: false,
        amount: 1,
      }
    },
    async mounted() {
      this.getNetworkConfig()
    },
  
    async created() {
      
    },
  
    methods: {
      onChangeSwitchLanguage(event) {
        this.$i18n.locale = event
        Cookies.set('language', event)
      },
  
      async getNetworkConfig(){
        try {
          let network= await rpc.getNetworkConfig()
          this.networkData = network.networkData
          this.managerAddress = network.managerAddress
          this.managerBalance = network.managerBalance
        } catch (e) {
          console.error(e)
          window.alert(e)
        }
      },

      async setNetworkConfig(){
        try {
          let resp = await rpc.setNetworkConfig(this.networkData)
          this.response = resp
          this.showResponse = true
        } catch (e) {
          console.error(e)
          window.alert(e)
        }
      },

      async authorizeMember(address){
        let resp = await rpc.authorizeMember(address)
        if (resp == 'success') {
          this.getNetworkConfig() // to update ip address assigned to the authorized node.
          this.toSetAcceptAddress(address)
        } 
      },

      async removeMember(address){
        let that = this
        that.confirm = false
        that.confirmTitle = 'Are you sure to remove this member?'
        that.confirmDialog = true
        await new Promise((resolve, reject) => {
          let interval = setInterval(() => {
            if (that.confirmDialog == false) {
              clearInterval(interval)
              if (that.confirm) {
                resolve()
              } else {
                reject()
              }
            }
          }, 100)
        }).then(async () => {
          let resp = await rpc.removeMember(address)
          if (resp == 'success') {
            let item = that.networkData.member[address]
            delete that.networkData.member[address]
            that.networkData.waiting[address] = item
            that.response = resp
            that.showResponse = true
          }
        }).catch(() => {
          console.log('cancel')
        })
        
      },

      async sendToken(address){
        let that = this
        that.confirm = false
        that.sendDialog = true
        await new Promise((resolve, reject) => {
          let interval = setInterval(() => {
            if (that.sendDialog == false) {
              clearInterval(interval)
              console.log("confirm", that.confirm)
              if (that.confirm) {
                resolve()
              } else {
                reject()
              }
            }
          }, 100)
        }).then(async () => {
          if (this.amount>0) {
            console.log("send", this.amount, " token to ", address)
            let resp = await rpc.sendToken(address, this.amount)
            that.response = resp
            that.showResponse = true
          }
        }).catch(() => {
          console.log('cancel')
        })
      },

      async deleteWaiting(address){
        let resp = await rpc.deleteWaiting(address)
        if (resp == 'success') {
          delete this.networkData.waiting[address]
          this.response = resp
          this.showResponse = true
        }
      },

      async nknPing(address){
        let resp = await rpc.nknPing(address)
        if (resp.includes('success')) {
          this.response = resp
          this.showResponse = true
        }
      },

      clickAllAddress(){
        if (this.allAddress.length>0) {
          this.acceptAddress = []
          for (const address in this.networkData.member) {
            this.acceptAddress.push(address)
          }
        } else {
          this.acceptAddress = []
        }
      },

      selectAddress(address){
        if (!this.acceptAddress.includes(address)) {
          this.allAddress = []
        } 
      },

      toSetAcceptAddress(address){
        this.setAddress = address
        this.allAddress = []
        this.acceptAddress = this.networkData.acceptAddress[address]
        if (!this.acceptAddress) {
          this.acceptAddress = []
        }
        if (this.acceptAddress && this.acceptAddress.length>0 ){
          if ( this.acceptAddress[0] == 'allMembers'){
            this.allAddress = ['allMembers']
            this.acceptAddress = []
            for (let address in this.networkData.member) {
              this.acceptAddress.push(address)
            }
          }
        }
        
        this.dialog = true
      },

      async setAcceptAddress(){
        this.dialog = false
        let acceptAddress = this.acceptAddress
        if (this.allAddress.length>0 ) {
          acceptAddress = ['allMembers']
        }
        let resp = await rpc.setAcceptAddress(this.setAddress, acceptAddress)
        
        if (resp == 'success') {
          this.networkData.acceptAddress[this.setAddress] = acceptAddress
        }
        this.response = resp
        this.showResponse = true
      }

    }
  }
</script>

<style lang="scss">
  header {
    height: 140px;
    display: flex;
    align-items: center;
  }
  
  .container {
    * {
      z-index: 1 !important;
    }
  }
  
  .nav {
    .nav-item {
      margin-bottom: 10px !important;
      margin-left: 10px !important;
    }
  }

  table {
    width: 100%;
    margin-top: 1rem;
    border: 1px solid #888;
    border-collapse: collapse;
    text-align: center;
    padding-left: 10px;
    padding-right: 10px;
    font-size: small;
  }
  th, td {
    width: 30px;
    margin-top: 1rem;
    border: 1px solid #888;
    border-collapse: collapse;
    text-align: center;
    padding-left: 10px;
    padding-right: 10px;
    font-size: small;
  }

  tr {
    height: 3rem;
  }
  tr:hover{
    background-color: #555;
  }

  .expand{
    width: 4rem; 
    text-align: center;
  }
  .expand:hover{
    background-color: #555;
    cursor: pointer;
  }

  .blankline{
    display: block;
    height: 1rem;
  }

</style>
  