<script setup lang="ts"></script>

<template>
  <!-- TODO: use the README.md to generate parts of this page -->
  <div class="main">
    <p>
      Robot Game - Legendary Ultimate Arena is a competitive game featuring turn based battles where
      two teams of robots fight in an arena. Players write the lua script controlling robots of
      their team, to make them move, attack, defend or even explode. The goal is to have more robots
      than the opposing team at the end of the battle. Robot Game - Legendary Ultimate Arena is a
      competitive game featuring turn based battles where two teams of robots fight in an arena.
      Players write the lua script controlling robots of their team, to make them move, attack,
      defend or even explode. The goal is to have more robots than the opposing team at the end of
      the battle.
    </p>
    <p>Robots have 10ms to act, they have one action per turn:</p>
    <p>Robots have 10ms to act, they have one action per turn:</p>
    <ul>
      <li><code>MOVE</code> : moves the robot to an adjacent tile.</li>
      <li><code>ATTACK</code> : attack an adjacent tile.</li>
      <li>
        <code>GUARD</code> : stay in place, protect against collision damage and take half damage
        from attacks and suicides.
      </li>
      <li><code>SUICIDE</code> : dies but deals damage to adjacent tiles.</li>
    </ul>
    <p>
      There are no friendly damage (collision, attack and suicide only deal damage to enemy
      robots).<br />
      They are several waves in which robots spawn randomly, The team that have more robots at the
      end of the game wins. There are no friendly damage (collision, attack and suicide only deal
      damage to enemy robots).<br />
      They are several waves in which robots spawn randomly, The team that have more robots at the
      end of the game wins.
    </p>
    <p>
      Get more info in the documentation in the
      <a href="https://github.com/MoutonGrognon/robot-game-lua#documentation">github repository</a>
    </p>
  </div>
</template>

<style lang="css" scoped>
.main {
  padding: 24px;
}
.main {
  padding: 24px;
}

code {
  background-color: black;
  color: #00ff41;
  padding-left: 2px;
  padding-right: 2px;
  border-radius: 2px;
}
code {
  background-color: black;
  color: #00ff41;
  padding-left: 2px;
  padding-right: 2px;
  border-radius: 2px;
}
</style>
