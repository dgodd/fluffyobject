<template>
  <div id="app">
    <h1>Users</h1>
    <ul>
      <li v-for="user in users">
        {{user.Email}}
      </li>
    </ul>

    <h1>Objects</h1>
    <ul class="objects">
      <li v-for="object in objects">
        <div class="object" v-bind:style="{ backgroundImage: 'url('+object.Image+')' }"></div>
        <b>{{object.Name}}</b>
      </li>
    </ul>
  </div>
</template>

<script>
export default {
  name: 'app',
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
  }
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
  ul.objects li {
    display: block;
    width: 200px; height: 200px;
    position: relative;
  }
  ul.objects li div.object {
    width: 100%;
    height: 100%;
    background-size: cover;
    background-position: center;
    border-radius: 100%;
    box-shadow: rgba(0,0,0,0.8) 0 0 20px;
    border-collapse: collapse;
  }
  ul.objects li b {
    display: block;
    text-align: center;
    line-height: 200px;
    position: absolute;
    top: 0; left: 0; right: 0;
    color: red;
    font-size: 200%;
  }
</style>
