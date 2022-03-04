import Vue from 'vue'
import VueI18n from 'vue-i18n'
import en from '~/locales/en.json'
import zh from '~/locales/zh-CN.json'
import zhTW from '~/locales/zh-TW.json'

Vue.use(VueI18n)

export default ({ app, store }) => {
  let messages = {}
  messages = {...messages, en, zh, zhTW}
  let lang = 'en'
  if (typeof navigator !== 'undefined') {
    let navLang = navigator.language || navigator.userLanguage
    if (!!navLang) lang = navLang.substr(0, 2)
  }
  app.i18n = new VueI18n({
    locale: lang,
    fallbackLocale: 'en',
    messages,
  });

  let locales = []
  for (let code of app.i18n.availableLocales) {
    let name = messages[code].language
    locales.push({code, name})
  }
  app.i18n.locales = locales
  app.i18n.path = (link) => {
    if (app.i18n.locale === app.i18n.fallbackLocale) {
      return `/${link}`;
    }
    return `/${app.i18n.locale}/${link}`;
  }
}
