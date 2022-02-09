<template>
  <v-app dark>
    <BackgroundLinear/>
    <v-main>
      <Nuxt/>
    </v-main>
  </v-app>
</template>

<script>
import BackgroundLinear from '../components/BackgroundLinear';
import '~/assets/global.scss'
export default {
  name: 'DefaultLayout',
  components: {BackgroundLinear},
  data() {
    return {
      clipped: false,
      drawer: false,
      fixed: false,
      items: [
        {
          icon: 'mdi-apps',
          title: 'Welcome',
          to: '/'
        },
        {
          icon: 'mdi-chart-bubble',
          title: 'Inspire',
          to: '/inspire'
        }
      ],
      miniVariant: false,
      right: true,
      rightDrawer: false,
      title: 'Vuetify.js'
    }
  }
}
</script>
