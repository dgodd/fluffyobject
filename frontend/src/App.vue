<template>
  <div id="app">
    <h1>Users</h1>
    <ul>
      <li v-for="user in users">
        <drag :transfer-data="{ userID: user.ID }">
          {{user.Email}}
        </drag>
      </li>
    </ul>

    <h1>Objects</h1>
    <ul class="objects">
      <MyObject v-for="object in objects" :key="object.ID" :obj="object"></MyObject>
    </ul>
  </div>
</template>

<script>
import { Drag } from 'vue-drag-drop';
import MyObject from './MyObject.vue';

export default {
  components: { Drag, MyObject },
  data () {
    return {
      users: [
        { ID: 1, Email: 'dave@goddard.id.au' },
        { ID: 2, Email: 'ecasey@pivotal.io' },
      ],
      objects: [
        { ID: 1, Name: 'Bastille', Image: 'https://c1.staticflickr.com/5/4309/35795982981_7d32a407a5_k.jpg' },
      ],
    }
  },
}
</script>

<style lang="css">
  body {
    margin:0; padding: 0;
    @import url('https://fonts.googleapis.com/css?family=Montserrat');
    font-family: 'Montserrat', sans-serif;
  }
  #app {
  }
</style>
