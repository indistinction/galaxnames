<template>
<div class="choosetop">
  <div class="back">
    <router-link to="/">&lt; back to race selection</router-link>
  </div>
  <h1>Who are you?</h1>
  <div id="glxcontainer" v-if="chooseGlx.length > 0">
    <div v-for="glx in chooseGlx.slice(0, show)" :key="glx.token_id">
      <img :src="glx.image" class="glximg">
      <h2>
        <span v-if="glx.givenname">
          {{ glx["Title"] ? `${glx["Title"]} ${glx.givenname}` : glx.givenname }}
        </span>
        <span v-else>
          {{ glx["Title"] ? `${glx["Race"]} ${glx["Title"]} #${glx.token_id}` : glx.name }}
        </span>
      </h2>
      <p v-if="glx.nickname">"{{ glx.nickname }}"</p>
      <p v-else>No name earned yet</p>
      <button v-if="glx.nickname" @click="$emit('clicked');$router.push(`/story/${glx.token_id}`)">Re-read your story</button>
      <button v-else @click="$emit('clicked');$router.push(`/play/${glx.token_id}`)">{{glx.level ? "Continue!" : "Begin!"}}</button>
    </div>
    <button @click="$emit('clicked');show+=step" v-if="show < chooseGlx.length">Show more...</button>
  </div>
  <div v-else>
    <p>ERROR: NO GALAXIATORS OF THAT RACE FOUND!</p>
    <button @click="$router.push('/')">Return to Race Selection</button>
  </div>
</div>
</template>

<script>
import {glxNameStore} from "@/store"

export default {
  name: "Choose",
  emits: ['clicked', 'flash', 'play'],
  data(){
    return{
      chooseGlx:[],
      show: 5,
      step: 5,
      needUpdate: true
    }
  },
  methods:{
    update(){
      if (!glxNameStore.state.authToken){
        this.$emit('flash', "Please log in.")
        this.$router.push("/")
      } else if(this.$activeRaces.indexOf(this.$route.params.race)===-1){
        // Haven't activated the game for this race yet
        this.$router.push("/")
      } else {
        this.$emit('play', this.$route.params.race)
          this.chooseGlx = []
          for (const glx in glxNameStore.state.glxList) {
            if(glxNameStore.state.glxList[glx].name.toLowerCase().startsWith(this.$route.params.race)){
              this.chooseGlx.push(glxNameStore.state.glxList[glx])
            }
          }
      }
    }
  },
  mounted() {
    this.update()
  }
}
</script>

<style scoped>

.choosetop{
  color: #fff;
  text-shadow: black 2px 2px 4px;
  text-align: center;
  padding: 18px;
  max-width: 800px;
  box-sizing: border-box;
  width: 100%;
}

#glxcontainer{
  display: flex;
  flex-wrap: wrap;
  justify-content: space-between;
  align-items: flex-start;
}

#glxcontainer div{
  width: 325px;
  padding: 16px;
  box-sizing: border-box;
}

.glximg{
  width: 100%;
  display: block;
  margin: 0 auto;
  border-radius: 3px;
}

button{
  background: #e3d947;
  transition: background-color 0.15s ease-in-out;
}
button:hover{
  background: #ff8f00;
}

p{
  margin-top: 0;
  box-sizing: border-box;
}

h2{
  margin-bottom: 0;
}
a{
  color: #e3d947;
  text-decoration: none;
}
</style>