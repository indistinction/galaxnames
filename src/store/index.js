import {reactive} from 'vue'

export const glxNameStore = {

  state: reactive({
    user: null,
    authToken: null,
    glxList: {}
  }),

  async login(loginJsonResponse) {
    this.state.user = loginJsonResponse.user
    this.state.authToken = loginJsonResponse.token

    let glxArray = []
    let assetRequest
    let assetJson
    let assetCursor

    // Get this users Glx token from IMX
    do {
      assetRequest = await fetch(`https://api.x.immutable.com/v1/assets?collection=0x6c82e53cbbd8a6afaf9663d58547cfc1a43be7aa&user=${loginJsonResponse.user}`)
      assetJson = await assetRequest.json()
      glxArray = glxArray.concat(assetJson.result)
      if(assetJson.remaining > 0){
        assetCursor = assetRequest.cursor
      } else {
        assetCursor = null
      }
    } while (assetCursor)

    // Need temp obj as non-zero number of keys triggers stuff elsewhere
    let tempGlxObj = {}
    let tempTokenList = []
    for (let i=0; i < glxArray.length; i++) {
      tempTokenList.push(glxArray[i]["token_id"])
      tempGlxObj[glxArray[i]["token_id"]] = glxArray[i]["metadata"]
      tempGlxObj[glxArray[i]["token_id"]]["token_id"] = glxArray[i]["token_id"]
      // Duplicate token_id as key and in data (makes life easier when appending to arrays etc)
    }

    // Get minigame data from DB
    let resp = await fetch("https://mini1.galaxiators.net/glxinfo"//"http://127.0.0.1:8001/glxinfo"
      , {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(tempTokenList)
    });

    if (resp.status !== 200){
      // TODO handle error
    } else {
      let dbGlxInfo = await resp.json()
      for (let glx in dbGlxInfo) {
        tempGlxObj[glx] = {...tempGlxObj[glx], ...dbGlxInfo[glx]}
      }
    }

    this.state.glxList = tempGlxObj
  },

  updateGlx(id, data){
    this.state.glxList[id] = {...this.state.glxList[id], ...data}
  },

  logout(){
    this.state.glxList = {}
    this.state.user = null
    this.state.authToken = null
  },

}