<template>
<div class="playtop">
  <div class="back">
    <router-link :to="glx ? `/glx/${glx.name.toLowerCase().charAt(0)}` : '/' ">&lt; back to ludus</router-link>
  </div>
    <div v-if="glx && glx.givenname" class="playtext">
        <h1>
          The story of {{ glx["Title"] ? `${glx["Title"]} ${glx.givenname}` : glx.givenname }}
        </h1>
      <img :src="`/assets/${glx.name.toLowerCase().charAt(0)}.png`" class="imgtopper">
        <div v-if="level && level['name']" :key="level['level_id']">
          <p v-html="level['outc'].replaceAll('\r\n', '<br>')" />
          <p>Whenever a new Galaxiator is 'recruited', it is the privilege of the Hunter to provide them with a gladiatorial name for the Galaxiators arena. From today onwards you will be known as</p>
          <h2>{{level['name']}}</h2>
          <p>Your initial Galaxiator recruitment stipend has been awarded at <strong>$GLXR {{level['$glxr'] ? level['$glxr'].toString() : "0"}}</strong> CyberSand</p>
          <!-- TODO button for download recruitment certificate? -->
          <ShareNetwork
              network="twitter"
              :url="`https://names.galaxiators.net/#/story/${glx.token_id}`"
              title="My Galaxiator just earned their arena nickname! Check out their story here."
              twitter-user="galaxiators"
              class="sharer twitter"
          >
            <img src="/assets/twitter.svg"/> Share on Twitter
          </ShareNetwork>
          <button @click="shareDiscord" class="sharer discord"><img src="/assets/discord.svg"/> Share on Discord</button>
          <button @click="$router.push('/glx/'+glx['name'].charAt(0).toLowerCase())">Play with another Galaxiator</button>
        </div>
        <div v-else-if="level && level['ans']" :key="level['level_id']">
          <p v-html="level['text'].replaceAll('\r\n', '<br>')" />
          <p v-if="level['level_id'].length === 1">
            BE WARNED. YOUR DECISIONS CANNOT BE UNDONE. THE CONSEQUENCES OF YOUR ACTIONS LIVE WITH YOUR GALAXIATOR FOR THE REST OF THEIR LIFE.
            <br><br>
            Five decisions stand between you and your new name. This is the point of no return.
          </p>
          <button v-for="(value, key) in level['ans']" :key="key" @click="setLevel(value['next'])">{{value["text"]}}</button>
        </div>
        <div v-else key="err">
          Error.
        </div>
    </div>
    <div v-else-if="glx">
      <h1>
        Welcome, {{ glx["Title"] ? `${glx["Race"]} ${glx["Title"]} #${glx.token_id}` : glx.name }}
      </h1>
      <p>The first of a Galaxiator's names is the one they're are given.</p>
      <p>So, what is your name?</p>
      <input type="text" v-model="nameInput" @keydown="keyDown">
      <p><small>(I hope your name is polite and appropriate, or Fleek will have to provide you with a new one of his choosing... Fluffy McRainbow-Gnome, for example)</small></p>
      <button @click="saveGiven">Save... Forever!</button>
    </div>
  <div v-else>
    <p>Login again to continue.</p>
    <button @click="$router.push('/')">Restart</button>
  </div>
</div>
</template>

<script>
import {glxNameStore} from "@/store"

export default {
  name: "Play",
  emits: ['flash', 'thinkon', 'thinkoff', 'clicked'],
  computed:{
    glx(){
      return glxNameStore.state.glxList[this.$route.params["glxid"]]
    }
  },
  data() {
    return {
      nameInput: "",
      level: null,
      allowedChars: "ABCDEFGHIJKLMNOPQRSTUVWXYZ abcdefghijklmnopqrstuvwxyz.-'"
    }
  },
  watch:{
    nameInput(val){
      let newVal = ""
      let char
      for(let i = 0;i < val.length; i++){
        char = val.charAt(i)
        if (this.testChar(char)){
          newVal += char
        } else {
          this.$emit('flash', `Cannot use the character ${char} in Galaxiator names. Sorry!`)
        }
      }
      this.nameInput = newVal
    }
  },
  methods: {
    // TODO In all of below need to deal with invalid/expired tokens
    // 403's return to main menu with a flash
    // Could DRY all requests (like secureFetch(data) returns success data and handles all fails)
    testChar(char){
      return !(this.allowedChars.indexOf(char) === -1)
    },
    async saveGiven() {
      this.$emit('clicked')
      this.$emit('thinkon', "Allocating given name...")
      let resp = await fetch(this.$backendhost + "/savegiven", {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'x-glx-token': glxNameStore.state.authToken,
          'x-glx-id': this.$route.params["glxid"]
        },
        body: JSON.stringify({
          "givenname": this.nameInput
        })
      });

      if (resp.status !== 204) {
        const errJson = await resp.json()
        this.$emit('thinkoff')
        this.$emit('flash', errJson.message)
      } else {
        glxNameStore.updateGlx(this.$route.params["glxid"], {
          "givenname": this.nameInput
        })
        this.$emit('thinkon', "Loading story...")
        let resp2 = await fetch(this.$backendhost + "/setl", {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'x-glx-token': glxNameStore.state.authToken,
            'x-glx-id': this.$route.params["glxid"]
          },
          body: JSON.stringify({
            "level": ""
          })
        })
        if (resp2.status !== 200) {
          const errJson = await resp2.json()
          this.$emit('thinkoff')
          this.$emit('flash', errJson.message)
        } else {
          this.level = await resp2.json()
          glxNameStore.updateGlx(this.$route.params["glxid"], {
            "level": this.level["level_id"]
          })
          this.$emit('thinkoff')
        }
      }
    },
    async updateLevel(){
      if (this.glx.hasOwnProperty('level')){
        this.$emit('thinkon', "Loading Galaxiator story...")
        let resp2 = await fetch(this.$backendhost + "/getl", {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'x-glx-token': glxNameStore.state.authToken,
            'x-glx-id': this.$route.params["glxid"]
          }
        })
        if (resp2.status !== 200) {
          const errJson = await resp2.json()
          this.$emit('thinkoff')
          this.$emit('flash', errJson.message)
        } else {
          let level = await resp2.json()
          if(level["name"]){
            glxNameStore.updateGlx(this.$route.params["glxid"], {
              "nickname": level["name"]
            })
          } else {
            let tempAns = {}
            let keys = Object.keys(level["ans"])
            keys = keys.sort(() => Math.random() - 0.5)
            for (let i = 0; i < keys.length; i++) {
              tempAns[keys[i]] = level["ans"][keys[i]]
            }
            level["ans"] = tempAns
          }
          this.level = level
          this.$emit('thinkoff')
        }
      }
    },
    async setLevel(levelId){
      this.$emit('clicked')
      this.$emit('thinkon', "Adventuring...")
      let resp2 = await fetch(this.$backendhost + "/setl", {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'x-glx-token': glxNameStore.state.authToken,
          'x-glx-id': this.$route.params["glxid"]
        },
        body: JSON.stringify({
          "level": levelId
        })
      })
      if (resp2.status !== 200) {
        const errJson = await resp2.json()
        this.$emit('flash', errJson.message)
        this.$emit('thinkoff')
      } else {
        let level = await resp2.json()
        let update = {"level": levelId}
        if(level["name"]){
          update["nickname"] = level["name"]
        } else {
          let tempAns = {}
          let keys = Object.keys(level["ans"])
          keys = keys.sort(() => Math.random() - 0.5)
          for (let i = 0; i < keys.length; i++) {
            tempAns[keys[i]] = level["ans"][keys[i]]
          }
          level["ans"] = tempAns
        }
        this.level = level
        glxNameStore.updateGlx(this.$route.params["glxid"], update)
        this.$emit('thinkoff')
      }
    },
    async shareDiscord(){
      this.$emit('thinkon', "Sharing to Discord...")
      let resp2 = await fetch(this.$backendhost + "/discoshare", {
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          'x-glx-token': glxNameStore.state.authToken,
          'x-glx-id': this.$route.params["glxid"]
        }
      })
      if (resp2.status !== 204) {
        const errJson = await resp2.json()
        this.$emit('flash', errJson.message)
        this.$emit('thinkoff')
      } else {
        this.$emit('flash', "Shared in the Galaxiators server!")
        this.$emit('thinkoff')
      }
    },
    keyDown(e){
      /*const specialKeys = [65]
      let char
      for(let i = 0;i < this.nameInput.length; i++){
        char = this.nameInput.charAt(i)
        if (!this.testChar(char)){
          val.replaceAll(char, "")
        }
      }


      if(this.allowedChars.indexOf(e.key) === -1){
        e.preventDefault()
        this.$emit('flash', "Invalid character for name. Sorry!")
      }*/
    }
  },
  mounted() {
    if(!glxNameStore.state.authToken){
      this.$emit('flash',"You are not logged in!")
      this.$router.push("/")
    } else {
      if (!this.glx) {
        this.$emit('flash', "No data found for that Galaxiator! Sign out, back in, and try again?")
      } else if (this.$activeRaces.indexOf(this.glx.name.toLowerCase().charAt(0))===-1) {
        this.$emit('flash', "That race's story hasn't begun yet.")
        this.$router.push("/")
      } else {
        this.updateLevel()
      }
    }
  }
}
</script>

<style scoped>
.playtop{
  color: #fff;
  text-shadow: black 2px 2px 4px;
  text-align: center;
  padding: 18px;
  margin: 0 auto !important;
  max-width: 800px;
  box-sizing: border-box;
  width: 100%;
  z-index: 2;
}
button{
  background: #e3d947;
  transition: background-color 0.15s ease-in-out;
}
button:hover{
  background: #ff8f00;
}

.back a{
  color: #e3d947;
  text-decoration: none;
}

.playtext{
  background: rgba(0, 0, 0, 0.75);
  color: #fff;
  text-align: justify;
  font-size: 14pt;
  padding: 36px 18px;
}

h2{
  width: 100%;
  display: block;
  text-align: center;
}

.discord{
  color: #fff;
  background: #7289da;
}

.imgtopper{
  width:100%;
  display: block;
  margin: 24px 0;
}

</style>