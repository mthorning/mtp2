import { Switch, Match } from "solid-js";
import type { Component } from "solid-js";
import { Gallery } from "./components";

const App: Component = () => {
  const { pathname } = window.location;

  return (
    <div class="p-16">
      <div class="container px-4 mx-auto">
        <Switch fallback={<h1>Not Found</h1>}>
          <Match when={pathname === "/"}>
            <Gallery selected={pathname.replace("/image/", "")} />
          </Match>
        </Switch>
      </div>
    </div>
  );
};

export default App;
