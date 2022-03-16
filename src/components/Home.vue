<template>
<div class="hometop">
  <div id="desktop">
    <img src="/assets/logo-y.png" class="logo">
    <p>Every Galaxiator has two names,<br />the one they're given and the one they earn.</p>
    <p>It's time for you to earn yours.</p>
    <p>
      <select v-model="race">
        <option value="" disabled selected hidden>Select your race...</option>
        <option value="c">Cylindurus</option>
        <option value="a" disabled>Abysstea (coming soon)</option>
        <option value="e" disabled>Egregore Z (coming soon)</option>
        <option value="l" disabled>Litagiar (coming soon)</option>
        <option value="s" disabled>Satyne (coming soon)</option>
        <option value="t" disabled>Thandrac (coming soon)</option>
      </select>
    </p>
    <p>
      <button @click="begin">Begin your story...</button>
    </p>
  </div>
  <!--div id="mobile">
    <img src="/assets/logo-y.png" style="width:280px;">
    <p><strong>***<br>MOBILE VIDSCREEN DEVICE DETECTED<br>***</strong></p>
    <p>Your Galaxiators deserve better.</p><p>Please access this service from a large-format vidscreen device.</p>
  </div-->
</div>
</template>

<script>
import {ethers} from "ethers";
import {glxNameStore} from "@/store"

export default {
  name: "Home",
  emits: ['flash', 'thinkon', 'thinkoff', 'clicked', 'logged'],
  data(){
    return{
      race: ""
    }
  },
  methods:{
    async begin(){
      this.$emit('clicked')
      if (!window.ethereum){
        this.$emit('flash', "You need an ethereum wallet to continue. Maybe MetaMask?")
      } else if(!this.race) {
        this.$emit('flash', 'Please select a race to continue.')
      } else if (glxNameStore.state.authToken){
        // Already logged in
        await this.$router.push(`/glx/${this.race}`)
      } else {
        this.$emit('thinkon', "Connect your MetaMask wallet...")
        // Connect to browser web3 provider - works with MetaMask who knows which other wallets
        const provider = new ethers.providers.Web3Provider(window.ethereum)
        try{
          await window.ethereum.request({ method: 'eth_requestAccounts' });
        }catch(e){
          this.$emit('flash', "Error connecting to MetaMask.")
          this.$emit('thinkoff')
          return
        }

        // Sign a message and send it to the backend to login
        const signer = provider.getSigner();
        let sig
        try{
          sig = await signer.signMessage("Log in to begin your Galaxiators story.")
        }catch(e){
          this.$emit('flash', "Error authorizing login. Please check your MetaMask.")
          this.$emit('thinkoff')
          return
        }

        this.$emit('thinkon', "Getting ludus roster...")
        let resp = await fetch(this.$backendhost + "/login", {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({sig})
        });

        if (resp.status !== 200) {
          this.$emit('flash', "Error logging in - no token returned.")
          this.$emit('thinkoff')
          return
        } else {
          // Save address and auth token
          // Also gets glx tokens for this user
          // And gets minigame data
          await glxNameStore.login(await resp.json())
          this.$emit('logged')
          this.$emit('thinkoff')
          await this.$router.push(`/glx/${this.race}`)
        }
      }
    }
  }
}
</script>

<style scoped>
.hometop{
  color: #fff;
  text-shadow: black 2px 2px 4px;
  text-align: center;
  padding: 18vh 18px 18px !important;
}

button{
  background: #e3d947;
  transition: background-color 0.15s ease-in-out;
}
button:hover{
  background: #ff8f00;
}

/*#mobile {
  display: none;
}

@media only screen and (max-width: 600px) {
  #mobile {
    display: block;
  }
  #desktop{
    visibility: hidden;
    height: 0;
  }
}*/
</style>