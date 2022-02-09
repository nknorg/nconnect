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
              <v-select class="nav-item select-language"
                        @change="onChangeSwitchLanguage(lang)"
                        v-model="lang"
                        :items="availableLocales"
                        item-value="code"
                        item-text="name"
                        label="Language"
                        solo dark
                        style="width: 200px;">
              </v-select>
            </span>
          </v-col>
        </v-row>
      </header>
      <v-row class="nav pr-2" align="center">
        <v-btn class="bg-linear-1 nav-item" :class="activeTab === 0 ? 'active' : ''" text @click="activeTab = 0">{{
            $t('mobile tab')
          }}
        </v-btn>
        <v-btn class="bg-linear-1 nav-item" :class="activeTab === 1 ? 'active' : ''" text @click="activeTab = 1">
          {{ $t('desktop tab') }}
        </v-btn>
        <div class="d-flex flex-column justify-center align-center" v-if="tuna && currentServerRegionName">
          <div class="mb-1">{{ $t('serverRegion') }}</div>
          <div style="max-width: 180px;">
            <v-select class="nav-item select-language"
                      @change="handleTunaConfigChoiceChange(tunaConfigSelectedValue)"
                      v-model="tunaConfigSelectedValue"
                      :items="tunaConfigChoices"
                      item-value="textId"
                      :label="$t('customized')"
                      solo dark
                      style="width: 200px;">
              <template v-slot:selection="{ item, index }">
                {{ $t(item.textId) || item.textId }}
              </template>
              <template v-slot:item="{ item, index }">
                {{ $t(item.textId) || item.textId }}
              </template>
            </v-select>
          </div>
        </div>
        <v-spacer/>
        <v-btn class="bg-linear-1 nav-item" :class="activeTab === 2 ? 'active' : ''" text @click="activeTab = 2">{{
            $t('help')
          }}
        </v-btn>
        <v-btn class="bg-linear-1 nav-item" :class="activeTab === 3 ? 'active' : ''" text @click="activeTab = 3">
          {{ $t('advance tab') }}
        </v-btn>
      </v-row>
      <v-row>
        <v-col cols="7" v-show="activeTab === 0">
          <v-row class="download-container">
            <v-col class="d-flex align-center">
              <strong>{{ $t('download nConnect part1') }} </strong>
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <img v-bind="attrs" v-on="on" src="~/static/img/qr_code.png" alt="QR Code" class="mx-1">
                </template>
                <img :src="downloadQrcode" alt="Download QR Code">
              </v-tooltip>
              <strong>{{ $t('download nConnect part2') }}</strong>
            </v-col>
          </v-row>
          <v-row class="mb-4">
            <v-col cols="auto">
              <img class="br-10" :src="adminTokenQRCode" alt="Admin Token QR Code">
            </v-col>
            <v-col>
              {{ $t('export tip') }}
            </v-col>
          </v-row>
          <div class="bg-linear-1 pa-4">
            <ul>
              <li class="mb-4">{{ $t('add device from mobile') }}</li>
              <li>{{ $t('connect from mobile') }}</li>
            </ul>
          </div>
          <v-btn class="bg-linear-2 pa-8 mt-4"
                 href="https://forum.nkn.org/t/nconnect-user-manual-video-nconnect/2457"
                 target="_blank"
                 text color="black" width="100%">
            {{ $t('desktop guide') }}
          </v-btn>
        </v-col>
        <v-col cols="7" v-show="activeTab === 1">
          <div class="bg-linear-1 pa-4 mb-4">
            <div v-html="$t('add device in mobile first')"></div>
          </div>
          <div class="bg-linear-1 pa-4 mb-4">
            <div v-html="$t('add server from desktop')"></div>
          </div>
          <div class="bg-linear-1 pa-4 mb-4">
            <div v-html="$t('scan QR code to add server to desktop')"></div>
          </div>
          <div class="bg-linear-1 pa-4 mb-8">
            <div v-html="$t('connect from desktop')"></div>
          </div>
          <v-btn class="bg-linear-2 pa-8 mt-4"
                 href="https://forum.nkn.org/t/nconnect-user-manual-video-nconnect/2457"
                 target="_blank"
                 text color="black" width="100%">
            {{ $t('desktop guide') }}
          </v-btn>
        </v-col>
        <v-col cols="7" v-show="activeTab === 2">
          <div class="bg-linear-1 pa-4 mb-8">
            <p class="mb-12" v-html="$t('need help method')"></p>
            <p v-html="$t('create forum post')"></p>
            <p v-html="$t('Q&A')"></p>
            <p v-html="$t('send email', {email: 'nconnect@nkn.org'})"></p>
            <p v-html="$t('mobile customer service')"></p>
          </div>
        </v-col>
        <v-col cols="7" v-show="activeTab === 3">
          <h3>{{ $t('local IP address') }}</h3>
          <v-textarea solo rows="3" disabled style="width: 244px" :value="localIP.join('\n')"></v-textarea>
          <h3>{{ $t('access key') }}</h3>
          <v-textarea solo rows="4" disabled :value="adminTokenStr"></v-textarea>
          <h3>{{ $t('accept addresses') }}</h3>
          <v-textarea solo v-model="acceptAddrs"></v-textarea>
          <h3>{{ $t('admins') }}</h3>
          <v-textarea solo v-model="adminAddrs"></v-textarea>
          <v-row>
            <v-col class="d-flex">
              <v-btn class="ml-auto btn-1" color="#00A3FF" @click="handleSubmit">
                {{ $t('save') }}
              </v-btn>
            </v-col>
          </v-row>
        </v-col>
        <v-col>
          <div class="bg-linear-1 pa-4 remaining-data">
            <strong>{{ $t('estimatedRemainingData') }}</strong>
            <div class="text-h3" v-if="tuna"><strong>{{ remainingData }}</strong></div>
            <div class="text-h3" v-else><strong>{{ $t('unlimited') }}</strong></div>
          </div>
          <v-flex class="mt-4 d-flex justify-end">
            <v-btn class="bg-linear-2 pa-8" text color="black" width="300" target="_blank"
                   :href="$t('paymentLink',{addr: addr, lng: lang, additionalParams: paymentAdditionalParams})">
              <strong>{{ $t('data plan tab') }}</strong>
            </v-btn>
          </v-flex>

          <v-row class="mt-4 justify-end mt-16" v-if="activeTab === 3">
            <v-col class="text-right">
              <v-btn class="bg-linear-1 mb-2" width="300" text @click="handleExportAccount">
                {{ $t('export account') }}
              </v-btn>
              <br>
              <v-btn class="bg-linear-1 mb-2" width="300" text @click="handleImportAccount">
                {{ $t('import account') }}
              </v-btn>
              <br>
              <v-btn class="bg-linear-1 mb-2" width="300" text @click="downloadLog">
                {{ $t('download log') }}
              </v-btn>
            </v-col>
          </v-row>

        </v-col>
      </v-row>
    </v-container>
  </v-container>
</template>

<script>
import Qrcode from 'qrcode'
import * as rpc from '../assets/rpc'

const tunaConfigChoicesAddr = '/static/tuna-config-choices.json';

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

export default {
  name: 'IndexPage',
  computed: {
    availableLocales() {
      return this.$i18n.locales
    }
  },
  data() {
    return {
      lang: this.$i18n.locale,
      tuna: false,
      activeTab: 0,
      initialized: false,
      downloadQrcode: '',
      adminTokenStr: '',
      adminTokenQRCode: '',
      acceptAddrs: '',
      adminAddrs: '',
      addr: '',
      localIP: [],
      inPrice: [],
      outPrice: [],
      tags: [],
      remainingData: 0,
      balance: '',
      tunaServiceName: '',
      tunaCountry: [],
      tunaConfigChoices: [],
      tunaConfigSelected: -1,
      tunaConfigSelectedValue: null,
      currentTunaConfig: -1,
      paymentAdditionalParams: '',
      currentServerRegionName: '',
    }
  },
  async mounted() {
    if (this.tags && this.tags.length) {
      for (let i = 0; i < this.tags.length; i++) {
        this.paymentAdditionalParams += '&tag=' + this.tags[i];
      }
    }

    await this.updateInfo();

    this.remainingData = this.estimatedRemainingData()

    if (this.tunaConfigChoices && this.tunaConfigChoices.length) {
      if (this.currentTunaConfig >= 0) {
        let item = this.tunaConfigChoices[this.currentTunaConfig];
        this.currentServerRegionName = this.$t(item.textId) || item.textId;
      } else {
        this.currentServerRegionName = this.$t('customized');
      }
    }

    setInterval(this.updateAdminToken, 5 * 60 * 1000);
  },
  async created() {
    this.downloadQrcode = await Qrcode.toDataURL('https://nconnect.org')
  },
  methods: {
    onChangeSwitchLanguage(event) {
      this.$router.replace(this.switchLocalePath(event));
    },
    async updateAdminToken() {
      try {
        let adminToken = await rpc.getAdminToken();
        if (adminToken) {
          this.adminTokenStr = JSON.stringify(adminToken);
          this.adminTokenQRCode = await Qrcode.toDataURL(this.adminTokenStr);
        }
      } catch (e) {
        console.error(e);
      }
    },
    async updateInfo() {
      this.updateAdminToken();

      let promise1 = rpc.getAddrs().then((addrs) => {
        this.acceptAddrs = addrsToStr(addrs.acceptAddrs)
        this.adminAddrs = addrsToStr(addrs.adminAddrs)
      }).catch((e) => {
        console.error(e);
        window.alert(e);
      })

      let promise2 = rpc.getInfo().then((info) => {
        console.log(info)
        let initialized = this.initialized

        this.initialized = true
        this.addr = info.addr
        this.localIP = info.localIP.ipv4
        this.inPrice = info.inPrice
        this.outPrice = info.outPrice
        this.tags = info.tags
        this.tunaServiceName = info.tunaServiceName || ''
        this.tunaCountry = info.tunaCountry || []
        this.tuna = info.tuna

        if (!initialized) {
          this.$axios.get(tunaConfigChoicesAddr + '?t=' + Date.now()).then((response) => {
            if (response.data && response.data.length) {
              this.tunaConfigChoices = response.data
              for (let i = 0; i < this.tunaConfigChoices.length; i++) {
                if (this.tunaConfigChoices[i].config.serviceName === this.tunaServiceName) {
                  if (JSON.stringify(this.tunaConfigChoices[i].config.country) === JSON.stringify(this.tunaCountry)) {
                    this.tunaConfigSelected = i
                    this.tunaConfigSelectedValue = this.tunaConfigChoices[i]
                    this.currentTunaConfig = i
                    break;
                  }
                }
              }
              // if (this.tunaConfigSelected < 0 && !this.tunaServiceName && !(this.tunaCountry && this.tunaCountry.length)) {
              //   this.openTunaConfigChoice();
              // }
            }
          }).catch((e) => {
            if (e.response.status !== 404) {
              console.error(e);
            }
          });
        }
      }).catch((e) => {
        console.error(e);
      })

      let promise3
      if (!this.balance) {
        promise3 = rpc.getBalance().then((balance) => {
          this.balance = balance
        }).catch((e) => {
          console.error(e);
        });
      } else {
        promise3 = Promise.resolve()
      }

      await Promise.all([promise1, promise2, promise3])
    },
    estimatedRemainingData() {
      if (!this.balance) {
        return null;
      }

      if (!(this.inPrice && this.inPrice.length) && !(this.outPrice && this.outPrice.length)) {
        return null;
      }

      let balance = parseFloat(this.balance);
      if (isNaN(balance)) {
        return null;
      }

      let maxPrice = 0;
      for (let i = 0; i < this.inPrice.length; i++) {
        let price = parseFloat(this.inPrice[i]);
        if (!isNaN(price) && price > maxPrice) {
          maxPrice = price;
        }
      }
      for (let i = 0; i < this.outPrice.length; i++) {
        let price = parseFloat(this.outPrice[i]);
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
    },
    handleTunaConfigChoiceChange(value) {
      this.tunaConfigSelected = this.tunaConfigChoices.findIndex((item) => item.textId === value)
      this.handleTunaConfigChoiceOK()
    },
    async handleTunaConfigChoiceOK() {
      this.currentTunaConfig = this.tunaConfigSelected
      if (this.tunaConfigSelected >= 0) {
        try {
          await rpc.setTunaConfig(this.tunaConfigChoices[this.tunaConfigSelected].config);
          window.alert(this.$t('setTunaConfigSuccess'));
        } catch (e) {
          console.error(e);
          window.alert(e);
        }
      }
    },
    async handleExportAccount() {
      if (!window.confirm(this.$t('exportConfirm'))) {
        return;
      }

      try {
        let seed = await rpc.getSeed();
        window.alert(this.$t('exportSuccess', {seed}));
      } catch (e) {
        console.error(e);
        window.alert(e);
      }
    },
    async handleImportAccount() {
      if (!window.confirm(this.$t('importConfirm'))) {
        return;
      }

      let currentSeedInput = window.prompt(this.$t('importPromptCurrent'));
      if (!currentSeedInput) {
        return;
      }

      try {
        let currentSeed = await rpc.getSeed();
        if (currentSeed !== currentSeedInput.trim()) {
          window.alert(this.$t('importWrongCurrent'));
          return;
        }
      } catch (e) {
        console.error(e);
        window.alert(e);
        return;
      }

      let newSeed = window.prompt(this.$t('importPromptNew'));
      if (!newSeed) {
        return;
      }

      try {
        await rpc.setSeed(newSeed.trim());
        window.alert(this.$t('importSuccess'));
      } catch (e) {
        console.error(e);
        window.alert(e);
      }
    },
    async downloadLog() {
      try {
        let log = await rpc.getLog();
        if (log && log.length) {
          downloadFile('nConnect.log', log);
        } else {
          window.alert(this.$t('no log available'))
        }
      } catch (e) {
        console.error(e);
        window.alert(e);
      }
    },
    async handleSubmit() {
      try {
        let addrs = await rpc.setAddrs(strToAddrs(this.acceptAddrs), strToAddrs(this.adminAddrs));

        this.acceptAddrs = addrsToStr(addrs.acceptAddrs)
        this.adminAddrs = addrsToStr(addrs.adminAddrs)

        window.alert(this.$t('save success'));
      } catch (e) {
        console.error(e);
        window.alert(e);
      }
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

  .select-language {

    .v-input__slot {

    }
  }
}

.v-parallax__content {
  justify-content: start;
}

.remaining-data {
  margin-left: auto;
  width: 300px;
  height: 125px;
}
</style>
