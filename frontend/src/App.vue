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
      objects: {},
      users: {},
      object_users: {},
    }
  },
  mounted() {
    const es = new EventSource("/api/events?stream=messages");
    es.addEventListener("users", x => this.users = JSON.parse(x.data));
    es.addEventListener("objects", x => this.objects = JSON.parse(x.data));
    // es.addEventListener("open", function() {
    //   console.log("ES OPEN 2");
    //   var httpRequest = new XMLHttpRequest();
    //   // if (!httpRequest) { alert('Giving up :( Cannot create an XMLHTTP instance'); }
    //   httpRequest.open('POST', '/api/senddata');
    //   httpRequest.send();
    // })

    var httpRequest = new XMLHttpRequest();
    httpRequest.open('POST', '/api/senddata');
    httpRequest.send();
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
