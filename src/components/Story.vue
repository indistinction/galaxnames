<template>
  <div>
    <div class="back">
      <a href="javascript:history.back()">&lt; back</a>
    </div>
    <div id="storycont">
      <div v-if="error">
        <h1>NO RECORD FOUND</h1>
        <p>Begin your journey <router-link to="/" style="color: #a700ff;">here</router-link>.</p>
      </div>
      <div v-else>
        <h1>Recruitment Record:<br>Galaxiator #{{$route.params.glxid}}</h1>
        <p style="text-align: center;font-size: 24pt;"><strong>Name: {{name}}</strong></p>
        <img :src="`https://storage.googleapis.com/glx-images/${$route.params.glxid}.png`" class="reportimg" />
        <ShareNetwork
            network="twitter"
            :url="`https://names.galaxiators.net/#/story/${$route.params.glxid}`"
            title="A Galaxiator just earned their arena nickname! Check out their story here."
            twitter-user="galaxiators"
            class="sharer twitter"
        ><img src="/assets/twitter.svg"/> Share on Twitter</ShareNetwork>
        <p class="spacer">--- begin recruitment report --</p>
        <div v-html="story" />
        <p class="spacer">--- end recruitment report --</p>
        <p>Arena nickname assigned by recruiting Hunter...</p>
        <h1>{{nick}}</h1>
        <br><br><ShareNetwork
          network="twitter"
          :url="`https://names.galaxiators.net/#/story/${$route.params.glxid}`"
          title="A Galaxiator just earned their arena nickname! Check out their story here."
          twitter-user="galaxiators"
          class="sharer twitter"
      ><img src="/assets/twitter.svg"/> Share on Twitter</ShareNetwork>
      </div>
    </div>
  </div>

</template>

<script>
export default {
  name: "Story",
  emits:['thinkon', 'thinkoff'],
  data(){
    return{
      story:"",
      name:"",
      nick: "",
      error: false
    }
  },
  async mounted(){
    this.$emit('thinkon', "Rewinding history...")
    let resp2 = await fetch(this.$backendhost + "/story/" + this.$route.params.glxid, {
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      }
    })
    if (resp2.status !== 200) {
      // TODO this didn't seem to work with a 404 returned? Do we need a try/catch?
      const errJson = await resp2.json()
      this.$emit('flash', errJson.message)
      this.error = true
      this.$emit('thinkoff')
    } else {
      let respJson = await resp2.json()
      if (!respJson.nickname){
        this.story = "(recruitment in process - awaiting outcome)"
        this.nick = "(pending)"
      } else {
        this.story = respJson.story.replaceAll("\r\n", "<br>")
        this.nick = respJson.nickname
      }
      this.name = respJson.givenname
      this.$emit('thinkoff')
    }
  }
}
</script>

<style scoped>
#storycont{
  background: #fff;
  color: #000;
  box-shadow: #000 2px 2px 5px;
  padding: 48px 18px;
  margin: 0 0 96px;
  border: 1px solid black;
}

#storycont h1{
  text-shadow: none;
  color: #000;
}

.spacer{
  width:100%;
  text-align: center;
  font-weight: bold;
}
a{
  color: #e3d947;
  text-decoration: none;
}
.reportimg{
  width: 300px;
  border-radius: 3px;
  margin: 12px auto;
  display: block;
}
</style>