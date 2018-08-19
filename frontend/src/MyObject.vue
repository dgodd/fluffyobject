<template>
  <li style="display: block; width: 200px; height: 200px; position: relative;">
    <drop @drop="handleDrop" style="width:100%;height:100%">
      <div class="object" v-bind:style="{ backgroundImage: 'url('+obj.Image+')' }"></div>
      <b>{{obj.Name}}</b>
    </drop>
  </li>
</template>

<script>
import { Drop } from 'vue-drag-drop';

export default {
  components: { Drop },
  props: { obj: { type: Object, required: true } },
  methods: {
    handleDrop (data) {
      alert(`You dropped with data: ${data.userID} => ${this.obj.ID}`);
    },
  },
}
</script>

<style lang="css" scoped>
  div.object {
    width: 100%;
    height: 100%;
    background-size: cover;
    background-position: center;
    border-radius: 100%;
    box-shadow: rgba(0,0,0,0.8) 0 0 20px;
    border-collapse: collapse;
  }
  b {
    display: block;
    text-align: center;
    line-height: 200px;
    position: absolute;
    top: 0; left: 0; right: 0;
    font-size: 200%;
    color: red;
    text-shadow: 2px 2px rgba(0,0,0,0.8);
  }
</style>

