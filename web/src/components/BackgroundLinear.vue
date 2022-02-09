<template>
  <v-container fluid class="absolute">
    <div class="background-1"></div>
    <div class="background-2"></div>
    <!--    <div class="background-3"><img src="~/static/img/grid.png" alt=""></div>-->
  </v-container>
</template>

<style lang="scss">
.absolute {
  position: absolute;
  height: 100%;
}

.background-1 {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  background: linear-gradient(350.57deg, #4F86DD -7.44%, #4D00A7 56.56%, #DE00FF 127.39%);
  transform: rotate(-180deg);
}

.background-2 {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  background: linear-gradient(0.41deg, #000000 -6.11%, rgba(18, 48, 105, 0.703257) 62.1%, rgba(79, 139, 255, 0.62) 92.53%);
}

.background-3 {
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;

  img {
    position: absolute;
    right: 0;
    bottom: 0;
  }
}
</style>
