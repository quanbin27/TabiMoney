import { createVuetify } from 'vuetify'
import { aliases, mdi } from 'vuetify/iconsets/mdi'
import * as components from 'vuetify/components'
import * as directives from 'vuetify/directives'

import 'vuetify/styles'
import '@mdi/font/css/materialdesignicons.css'

import { lightTheme, darkTheme } from './theme'

export default createVuetify({
  components: { ...components },
  directives,
  icons: {
    defaultSet: 'mdi',
    aliases,
    sets: { mdi },
  },
  theme: {
    defaultTheme: 'light',
    themes: {
      light: lightTheme,
      dark: darkTheme,
    },
  },
  defaults: {
    VBtn: { style: 'text-transform: none;' },
    VCard: { elevation: 2 },
    VTextField: { variant: 'outlined' },
    VSelect: { variant: 'outlined' },
    VTextarea: { variant: 'outlined' },
  },
})


