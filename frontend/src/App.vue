<template>
  <div id="app" unselectable="on" onselectstart="return false;">
    <img src="" id="bg">
    <TabBar
      v-on:send="switchToSend"
      v-on:wallet="switchToWallet"
      v-on:settings="switchToSettings"
      v-if="((screen === 'welcome' && manualStop) || screen !== 'welcome') && (screen !== 'checks' || tabBarVisible)"
    />
    <Welcome  v-if="screen === 'welcome'" v-on:start-mining="switchToChecks" />
    <Checks v-if="screen === 'checks'" v-on:mining="switchToMining" v-on:checksFailed="showTabBar" />
    <Send v-if="screen === 'send'" v-on:back="switchToMining" v-on:cancel="switchToMining" />
    <Mining v-show="screen === 'mining'" v-on:stop-mining="stopMining" />
    <Settings v-if="screen === 'settings'" v-on:committed="restartMining" />
    <Update v-if="screen === 'update'" v-on:back="restartMiningIfNotStopped" />
    <Tracking v-on:update="switchToUpdate" />
    <div class='buttonTheme' >
      <img v-if="theme === 'theme-light'" src="../public/images/theme_dark_icon.png" @click="changeTheme">
      <img v-if="theme === 'theme-dark'" src="../public/images/theme_light_icon.png" @click="changeTheme">
    </div>

  </div>
  
</template>

<script>
import Welcome from "./components/Welcome.vue";
import Mining from "./components/Mining.vue";
import Checks from "./components/Checks.vue";
import Send from "./components/Send.vue";
import Settings from "./components/Settings.vue";
import TabBar from "./components/TabBar.vue";
import Tracking from "./components/Tracking.vue";
import Update from "./components/Update.vue";

import "./assets/css/main.css";

export default {
 
  data() {
    return {
      theme: '',
      screen: "welcome",
      manualStop: false,
      tabBarVisible: false
    };
  },
  
  mounted() {
    var self = this;
    window.wails.Events.On("minerRapidFail", () => {
      window.backend.Backend.StopMining().then(() => {
        self.switchToChecks();
      });
    });
    window.backend.Backend.ReadBg().then((value) => {
      document.getElementById('bg').src = value;
    })
    window.backend.Backend.ReadTheme().then((value) => {
      this.theme = (value === '' || value == null) ? 'theme-dark' : value.trim();
      let body = document.getElementsByTagName('body')[0];
      body.classList.remove("theme-light", "theme-dark");
      body.classList.add(this.theme);
    })
  },
  methods: {
    changeTheme(){
      this.theme = (this.theme === 'theme-light') ? 'theme-dark' : 'theme-light';
      let body = document.getElementsByTagName('body')[0];
      body.classList.remove("theme-light", "theme-dark");
      body.classList.add(this.theme);

      window.backend.Backend.SaveTheme(this.theme);
    },
    

    stopMining: function() {
      this.manualStop = true;
      this.switchToWelcome();
    },
    // Target for the wallet tab (meta between welcome (if stopped) and mining (if mining))
    switchToWallet: function() {
      var self = this;
      if (this.tabBarVisible === true) {
        this.tabBarVisible = false;
        window.backend.Backend.StopMining().then(() => {
          self.switchToChecks();
        });
      } else {
        if (this.manualStop) {
          this.switchToWelcome();
        } else {
          this.switchToMining();
        }
      }
    },
    showTabBar: function() {
      this.tabBarVisible = true;
    },
    switchToChecks: function() {
      this.screen = "checks";
    },
    switchToSettings: function() {
      this.screen = "settings";
    },
    switchToSend: function() {
      this.screen = "send";
    },
    switchToMining: function() {
      this.manualStop = false;
      this.screen = "mining";
    },
    switchToUpdate: function() {
      this.screen = "update";
    },
    switchToWelcome: function() {
      this.screen = "welcome";
    },
    restartMiningIfNotStopped: function() {
      if(this.manualStop) {
        this.switchToWelcome();
      } else {
        this.restartMining();
      }
    },
    restartMining: function() {
      var self = this;
      window.backend.Backend.StopMining().then(() => {
        self.switchToChecks();
      });
    }
  },
  name: "app",
  components: {
    Welcome,
    Mining,
    Checks,
    Send,
    TabBar,
    Tracking,
    Settings,
    Update
  }
};
</script>
<style>
#bg {
  min-height: 100%;
  min-width: 512px;

  width: 100%;
  height: auto;

  position: fixed;
  top: 0;
  left: 0;
  pointer-events: none;
  z-index: -1;
}
</style>
