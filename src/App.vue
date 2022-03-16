<template>
  <!-- TODO Preload some mp3/images -->
  <div id="thinking" v-if="thinking">
    <div>
      <img src="/assets/g-y.png" />
      <p class="blink">{{thinkmsg}}</p>
    </div>
  </div>
  <div id="container">
    <div>
      <div id="topbar" :class="{'topbg': showsound || stateGlxCount > 0 }">
        <div v-if="showsound && playing">
          <img src="/assets/vol-y.png" class="soundicon" @click="soundOff"/>
        </div>
        <div v-else-if="showsound && !playing">
          <img src="/assets/mute-y.png" class="soundicon" @click="soundOn"/>
        </div>
        <div v-else>
          &nbsp;
        </div>
        <div v-if="stateGlxCount > 0">
          <button id="walletbutton" @click="signOut">Sign out</button>
        </div>
        <div v-else>
          &nbsp;
        </div>
      </div>
      <router-view
          id="content"
          @flash="flashMsg"
          @play="playMusic"
          @thinkon="thinkon"
          @clicked="click"
          @logged="loggedin=true"
          @thinkoff="thinking=false;thinkmsg=''"
      />
    </div>
    <div id="flash-container">
      <div class="flash-box" v-for="(msg, id) in flash" :key="id">{{msg}}</div>
    </div>
    <!-- music by https://www.fesliyanstudios.com/royalty-free-music/downloads-c/rock-music/16 -->
    <audio loop ref="c-bg">
      <source src="/assets/c-bg.mp3">
    </audio>
    <audio ref="c-click">
      <source src="/assets/c-click.mp3">
    </audio>
    <audio ref="-click">
      <source src="/assets/c-click.mp3">
    </audio>
  </div>
</template>
<script>
import {glxNameStore} from "@/store"

export default {
  name: 'App',
  computed: {
    stateGlxCount() {
      return Object.keys(glxNameStore.state.glxList).length
    }
  },
  data(){
    return{
      loggedin: false,
      thinking: false,
      thinkmsg: "",
      showsound: false,
      playing: false,
      flash: {}
    }
  },
  methods: {
    thinkon(msg){
      this.thinking=true
      this.thinkmsg=msg
    },
    flashMsg(msg){
      const id = Date.now()
      this.flash[id] = msg
      setTimeout(() => this.deleteMsg(id), 3000)
    },
    deleteMsg(id){
      delete this.flash[id]
    },
    soundOn(){
      this.playing = true
      for(let i = 0;i < this.$activeRaces.length;i++){
        this.$refs[`${this.$activeRaces[i]}-bg`].volume = 0.5
        this.$refs[`${this.$activeRaces[i]}-click`].volume = 1.0
      }
      this.$refs['-click'].volume = 1.0
    },
    soundOff(){
      this.playing = false
      for(let i = 0;i < this.$activeRaces.length;i++){
        this.$refs[`${this.$activeRaces[i]}-bg`].volume = 0.0
        this.$refs[`${this.$activeRaces[i]}-click`].volume = 0.0
      }
      this.$refs['-click'].volume = 0.0
    },
    playMusic(race){
      console.log(`${race}-bg`)
      if(!this.playing){
        this.playing = true
        this.showsound = true
        this.$refs[`${race}-bg`].play()
      }
    },
    click(race){
      if(!race){
        race = ""
      }
      this.$refs[`${race}-click`].play()
    },
    signOut(){
      glxNameStore.logout()
      this.$router.push("/")
    }
  }
}
</script>

<style>
#content{
  max-width: 800px;
  margin: 0 auto;
  box-sizing: border-box;
}
#topbar{
  height: 64px;
  position: fixed;
  box-sizing: border-box;
  padding: 18px;
  top:0;
  display: flex;
  justify-content: space-between;
  width: 100%;
}
#topbar div{
  height: 100%;
  display: block;
}

.topbg{
  background: rgba(0, 0, 0, 0.35);
}

.soundicon{
  height: 100%;
}

#walletbutton{
  width: 100px;
  height: 100%;
  font-size: 12pt;
  margin: 0 18pt 0 0;
  box-sizing: border-box;
  padding: 0;
  background: #e3d947;
  transition: background-color 0.15s ease-in-out;
}

#thinking{
  position: fixed;
  top:0;
  left:0;
  height: 100vh;
  width: 100vw;
  backdrop-filter: blur(3px);
  z-index: 10001;
  background-color: rgba(0, 0, 0, 0.85);
}

#thinking div{
  width:100%;
  height: 100%;
  padding: 40vh 0 0;
}

#thinking div img {
  padding: 0;
  display: block;
  height: 30px;
  margin: 0 auto;
  animation: rotateAnimation 1.5s linear infinite;
}

#thinking div p {
  width: 100%;
  color: #fff;
  text-align: center;
  padding: 8px;
  box-sizing: border-box;
}

@keyframes rotateAnimation {
  from { -webkit-transform: rotateY(0deg);    }
  to   { -webkit-transform: rotateY(-360deg); }
}

</style>
