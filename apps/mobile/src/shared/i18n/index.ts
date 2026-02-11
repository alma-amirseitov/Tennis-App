import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import * as Localization from 'expo-localization';

import ru from './locales/ru.json';
import kk from './locales/kk.json';
import en from './locales/en.json';

const deviceLang = Localization.getLocales()[0]?.languageCode ?? 'ru';

i18n.use(initReactI18next).init({
  resources: { ru: { translation: ru }, kk: { translation: kk }, en: { translation: en } },
  lng: deviceLang === 'kk' ? 'kk' : deviceLang === 'en' ? 'en' : 'ru',
  fallbackLng: 'ru',
  interpolation: { escapeValue: false },
});

export default i18n;
