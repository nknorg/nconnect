import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector';

import translationEN from './translations/en.json';
import translationZHCN from './translations/zh-CN.json';
import translationZHTW from './translations/zh-TW.json';

export const resources = {
  en: {
    translation: translationEN
  },
  "zh-CN": {
    translation: translationZHCN
  },
  "zh-TW": {
    translation: translationZHTW
  },
};

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'en',
    interpolation: {
      escapeValue: false,
    }
  });

export default i18n;
