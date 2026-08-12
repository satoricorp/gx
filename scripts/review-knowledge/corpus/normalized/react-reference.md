---
title: <Activity>
---

<Intro>

`<Activity>` lets you hide and restore the UI and internal state of its children.

```js
<Activity mode={visibility}>
 <Sidebar />
</Activity>
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `<Activity>` {/*activity*/}

You can use Activity to hide part of your application:

```js [[1, 1, "\\"hidden\\""], [2, 2, "<Sidebar />"], [3, 1, "\\"visible\\""]]
<Activity mode={isShowingSidebar ? "visible" : "hidden"}>
 <Sidebar />
</Activity>
```

When an Activity boundary is <CodeStep step={1}>hidden</CodeStep>, React will visually hide <CodeStep step={2}>its children</CodeStep> using the `display: "none"` CSS property. It will also destroy their Effects, cleaning up any active subscriptions.

While hidden, children still re-render in response to new props, albeit at a lower priority than the rest of the content.

When the boundary becomes <CodeStep step={3}>visible</CodeStep> again, React will reveal the children with their previous state restored, and re-create their Effects.

In this way, Activity can be thought of as a mechanism for rendering "background activity". Rather than completely discarding content that's likely to become visible again, you can use Activity to maintain and restore that content's UI and internal state, while ensuring that your hidden content has no unwanted side effects.

[See more examples below.](#usage)

#### Props {/*props*/}

* `children`: The UI you intend to show and hide.
* `mode`: A string value of either `'visible'` or `'hidden'`. If omitted, defaults to `'visible'`.

#### Caveats {/*caveats*/}

- If an Activity is rendered inside of a [ViewTransition](/reference/react/ViewTransition), and it becomes visible as a result of an update caused by [startTransition](/reference/react/startTransition), it will activate the ViewTransition's `enter` animation. If it becomes hidden, it will activate its `exit` animation.
- A *hidden* Activity that just renders text will not render anything rather than rendering hidden text, because there’s no corresponding DOM element to apply visibility changes to. For example, `<Activity mode="hidden"><ComponentThatJustReturnsText /></Activity>` will not produce any output in the DOM for `const ComponentThatJustReturnsText = () => "Hello, World!"`. `<Activity mode="visible"><ComponentThatJustReturnsText /></Activity>` will render visible text.

---

## Usage {/*usage*/}

### Restoring the state of hidden components {/*restoring-the-state-of-hidden-components*/}

In React, when you want to conditionally show or hide a component, you typically mount or unmount it based on that condition:

```jsx
{isShowingSidebar && (
 <Sidebar />
)}
```

But unmounting a component destroys its internal state, which is not always what you want.

When you hide a component using an Activity boundary instead, React will "save" its state for later:

```jsx
<Activity mode={isShowingSidebar ? "visible" : "hidden"}>
 <Sidebar />
</Activity>
```

This makes it possible to hide and then later restore components in the state they were previously in.

The following example has a sidebar with an expandable section. You can press "Overview" to reveal the three subitems below it. The main app area also has a button that hides and shows the sidebar.

Try expanding the Overview section, and then toggling the sidebar closed then open:

<Sandpack>

```js src/App.js active
import { useState } from 'react';
import Sidebar from './Sidebar.js';

export default function App() {
 const [isShowingSidebar, setIsShowingSidebar] = useState(true);

 return (
 <>
 {isShowingSidebar && (
 <Sidebar />
 )}

 <main>
 <button onClick={() => setIsShowingSidebar(!isShowingSidebar)}>
 Toggle sidebar
 </button>
 <h1>Main content</h1>
 </main>
 </>
 );
}
```

```js src/Sidebar.js
import { useState } from 'react';

export default function Sidebar() {
 const [isExpanded, setIsExpanded] = useState(false)

 return (
 <nav>
 <button onClick={() => setIsExpanded(!isExpanded)}>
 Overview
 <span className={`indicator ${isExpanded ? 'down' : 'right'}`}>
 &#9650;
 </span>
 </button>

 {isExpanded && (
 <ul>
 <li>Section 1</li>
 <li>Section 2</li>
 <li>Section 3</li>
 </ul>
 )}
 </nav>
 );
}
```

```css
body { height: 275px; margin: 0; }
#root {
 display: flex;
 gap: 10px;
 height: 100%;
}
nav {
 padding: 10px;
 background: #eee;
 font-size: 14px;
 height: 100%;
}
main {
 padding: 10px;
}
p {
 margin: 0;
}
h1 {
 margin-top: 10px;
}
.indicator {
 margin-left: 4px;
 display: inline-block;
 rotate: 90deg;
}
.indicator.down {
 rotate: 180deg;
}
```

</Sandpack>

The Overview section always starts out collapsed. Because we unmount the sidebar when `isShowingSidebar` flips to `false`, all its internal state is lost.

This is a perfect use case for Activity. We can preserve the internal state of our sidebar, even when visually hiding it.

Let's replace the conditional rendering of our sidebar with an Activity boundary:

```jsx {7,9}
// Before
{isShowingSidebar && (
 <Sidebar />
)}

// After
<Activity mode={isShowingSidebar ? 'visible' : 'hidden'}>
 <Sidebar />
</Activity>
```

and check out the new behavior:

<Sandpack>

```js src/App.js active
import { Activity, useState } from 'react';

import Sidebar from './Sidebar.js';

export default function App() {
 const [isShowingSidebar, setIsShowingSidebar] = useState(true);

 return (
 <>
 <Activity mode={isShowingSidebar ? 'visible' : 'hidden'}>
 <Sidebar />
 </Activity>

 <main>
 <button onClick={() => setIsShowingSidebar(!isShowingSidebar)}>
 Toggle sidebar
 </button>
 <h1>Main content</h1>
 </main>
 </>
 );
}
```

```js src/Sidebar.js
import { useState } from 'react';

export default function Sidebar() {
 const [isExpanded, setIsExpanded] = useState(false)

 return (
 <nav>
 <button onClick={() => setIsExpanded(!isExpanded)}>
 Overview
 <span className={`indicator ${isExpanded ? 'down' : 'right'}`}>
 &#9650;
 </span>
 </button>

 {isExpanded && (
 <ul>
 <li>Section 1</li>
 <li>Section 2</li>
 <li>Section 3</li>
 </ul>
 )}
 </nav>
 );
}
```

```css
body { height: 275px; margin: 0; }
#root {
 display: flex;
 gap: 10px;
 height: 100%;
}
nav {
 padding: 10px;
 background: #eee;
 font-size: 14px;
 height: 100%;
}
main {
 padding: 10px;
}
p {
 margin: 0;
}
h1 {
 margin-top: 10px;
}
.indicator {
 margin-left: 4px;
 display: inline-block;
 rotate: 90deg;
}
.indicator.down {
 rotate: 180deg;
}
```

</Sandpack>

Our sidebar's internal state is now restored, without any changes to its implementation.

---

### Restoring the DOM of hidden components {/*restoring-the-dom-of-hidden-components*/}

Since Activity boundaries hide their children using `display: none`, their children's DOM is also preserved when hidden. This makes them great for maintaining ephemeral state in parts of the UI that the user is likely to interact with again.

In this example, the Contact tab has a `<textarea>` where the user can enter a message. If you enter some text, change to the Home tab, then change back to the Contact tab, the draft message is lost:

<Sandpack>

```js src/App.js
import { useState } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Contact from './Contact.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('contact');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'contact'}
 onClick={() => setActiveTab('contact')}
 >
 Contact
 </TabButton>

 <hr />

 {activeTab === 'home' && <Home />}
 {activeTab === 'contact' && <Contact />}
 </>
 );
}
```

```js src/TabButton.js
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Contact.js active
export default function Contact() {
 return (
 <div>
 <p>Send me a message!</p>

 <textarea />

 <p>You can find me online here:</p>
 <ul>
 <li>admin@mysite.com</li>
 <li>+123456789</li>
 </ul>
 </div>
 );
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
```

</Sandpack>

This is because we're fully unmounting `Contact` in `App`. When the Contact tab unmounts, the `<textarea>` element's internal DOM state is lost.

If we switch to using an Activity boundary to show and hide the active tab, we can preserve the state of each tab's DOM. Try entering text and switching tabs again, and you'll see the draft message is no longer reset:

<Sandpack>

```js src/App.js active
import { Activity, useState } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Contact from './Contact.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('contact');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'contact'}
 onClick={() => setActiveTab('contact')}
 >
 Contact
 </TabButton>

 <hr />

 <Activity mode={activeTab === 'home' ? 'visible' : 'hidden'}>
 <Home />
 </Activity>
 <Activity mode={activeTab === 'contact' ? 'visible' : 'hidden'}>
 <Contact />
 </Activity>
 </>
 );
}
```

```js src/TabButton.js
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Contact.js
export default function Contact() {
 return (
 <div>
 <p>Send me a message!</p>

 <textarea />

 <p>You can find me online here:</p>
 <ul>
 <li>admin@mysite.com</li>
 <li>+123456789</li>
 </ul>
 </div>
 );
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
```

</Sandpack>

Again, the Activity boundary let us preserve the Contact tab's internal state without changing its implementation.

---

### Pre-rendering content that's likely to become visible {/*pre-rendering-content-thats-likely-to-become-visible*/}

So far, we've seen how Activity can hide some content that the user has interacted with, without discarding that content's ephemeral state.

But Activity boundaries can also be used to _prepare_ content that the user has yet to see for the first time:

```jsx [[1, 1, "\\"hidden\\""]]
<Activity mode="hidden">
 <SlowComponent />
</Activity>
```

When an Activity boundary is <CodeStep step={1}>hidden</CodeStep> during its initial render, its children won't be visible on the page — but they will _still be rendered_, albeit at a lower priority than the visible content, and without mounting their Effects.

This _pre-rendering_ allows the children to load any code or data they need ahead of time, so that later, when the Activity boundary becomes visible, the children can appear faster with reduced loading times.

Let's look at an example.

In this demo, the Posts tab loads some data. If you press it, you'll see a Suspense fallback displayed while the data is being fetched:

<Sandpack>

```js src/App.js
import { useState, Suspense } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Posts from './Posts.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('home');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'posts'}
 onClick={() => setActiveTab('posts')}
 >
 Posts
 </TabButton>

 <hr />

 <Suspense fallback={<h1>🌀 Loading...</h1>}>
 {activeTab === 'home' && <Home />}
 {activeTab === 'posts' && <Posts />}
 </Suspense>
 </>
 );
}
```

```js src/TabButton.js hidden
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Posts.js
import { use } from 'react';
import { fetchData } from './data.js';

export default function Posts() {
 const posts = use(fetchData('/posts'));

 return (
 <ul className="items">
 {posts.map(post =>
 <li className="item" key={post.id}>
 {post.title}
 </li>
 )}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/posts')) {
 return await getPosts();
 } else {
 throw Error('Not implemented');
 }
}

async function getPosts() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });
 let posts = [];
 for (let i = 0; i < 10; i++) {
 posts.push({
 id: i,
 title: 'Post #' + (i + 1)
 });
 }
 return posts;
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
video { width: 300px; margin-top: 10px; aspect-ratio: 16/9; }
```

</Sandpack>

This is because `App` doesn't mount `Posts` until its tab is active.

If we update `App` to use an Activity boundary to show and hide the active tab, `Posts` will be pre-rendered when the app first loads, allowing it to fetch its data before it becomes visible.

Try clicking the Posts tab now:

<Sandpack>

```js src/App.js
import { Activity, useState, Suspense } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Posts from './Posts.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('home');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'posts'}
 onClick={() => setActiveTab('posts')}
 >
 Posts
 </TabButton>

 <hr />

 <Suspense fallback={<h1>🌀 Loading...</h1>}>
 <Activity mode={activeTab === 'home' ? 'visible' : 'hidden'}>
 <Home />
 </Activity>
 <Activity mode={activeTab === 'posts' ? 'visible' : 'hidden'}>
 <Posts />
 </Activity>
 </Suspense>
 </>
 );
}
```

```js src/TabButton.js hidden
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Posts.js
import { use } from 'react';
import { fetchData } from './data.js';

export default function Posts() {
 const posts = use(fetchData('/posts'));

 return (
 <ul className="items">
 {posts.map(post =>
 <li className="item" key={post.id}>
 {post.title}
 </li>
 )}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/posts')) {
 return await getPosts();
 } else {
 throw Error('Not implemented');
 }
}

async function getPosts() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });
 let posts = [];
 for (let i = 0; i < 10; i++) {
 posts.push({
 id: i,
 title: 'Post #' + (i + 1)
 });
 }
 return posts;
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
video { width: 300px; margin-top: 10px; aspect-ratio: 16/9; }
```

</Sandpack>

`Posts` was able to prepare itself for a faster render, thanks to the hidden Activity boundary.

---

Pre-rendering components with hidden Activity boundaries is a powerful way to reduce loading times for parts of the UI that the user is likely to interact with next.

<Note>

Only data read from a source that [activates a Suspense boundary](/reference/react/Suspense#what-activates-a-suspense-boundary), such as a Promise read with [`use`](/reference/react/use), is fetched during pre-rendering. Activity does not detect data fetched inside an Effect.

</Note>

---

### Speeding up interactions during page load {/*speeding-up-interactions-during-page-load*/}

React includes an under-the-hood performance optimization called Selective Hydration. It works by hydrating your app's initial HTML _in chunks_, enabling some components to become interactive even if other components on the page haven't loaded their code or data yet.

Suspense boundaries participate in Selective Hydration, because they naturally divide your component tree into units that are independent from one another:

```jsx
function Page() {
 return (
 <>
 <MessageComposer />

 <Suspense fallback="Loading chats...">
 <Chats />
 </Suspense>
 </>
 )
}
```

Here, `MessageComposer` can be fully hydrated during the initial render of the page, even before `Chats` is mounted and starts to fetch its data.

So by breaking up your component tree into discrete units, Suspense allows React to hydrate your app's server-rendered HTML in chunks, enabling parts of your app to become interactive as fast as possible.

But what about pages that don't use Suspense?

Take this tabs example:

```jsx
function Page() {
 const [activeTab, setActiveTab] = useState('home');

 return (
 <>
 <TabButton onClick={() => setActiveTab('home')}>
 Home
 </TabButton>
 <TabButton onClick={() => setActiveTab('video')}>
 Video
 </TabButton>

 {activeTab === 'home' && (
 <Home />
 )}
 {activeTab === 'video' && (
 <Video />
 )}
 </>
 )
}
```

Here, React must hydrate the entire page all at once. If `Home` or `Video` are slower to render, they could make the tab buttons feel unresponsive during hydration.

Adding Suspense around the active tab would solve this:

```jsx {13,20}
function Page() {
 const [activeTab, setActiveTab] = useState('home');

 return (
 <>
 <TabButton onClick={() => setActiveTab('home')}>
 Home
 </TabButton>
 <TabButton onClick={() => setActiveTab('video')}>
 Video
 </TabButton>

 <Suspense fallback={<Placeholder />}>
 {activeTab === 'home' && (
 <Home />
 )}
 {activeTab === 'video' && (
 <Video />
 )}
 </Suspense>
 </>
 )
}
```

...but it would also change the UI, since the `Placeholder` fallback would be displayed on the initial render.

Instead, we can use Activity. Since Activity boundaries show and hide their children, they already naturally divide the component tree into independent units. And just like Suspense, this feature allows them to participate in Selective Hydration.

Let's update our example to use Activity boundaries around the active tab:

```jsx {13-18}
function Page() {
 const [activeTab, setActiveTab] = useState('home');

 return (
 <>
 <TabButton onClick={() => setActiveTab('home')}>
 Home
 </TabButton>
 <TabButton onClick={() => setActiveTab('video')}>
 Video
 </TabButton>

 <Activity mode={activeTab === "home" ? "visible" : "hidden"}>
 <Home />
 </Activity>
 <Activity mode={activeTab === "video" ? "visible" : "hidden"}>
 <Video />
 </Activity>
 </>
 )
}
```

Now our initial server-rendered HTML looks the same as it did in the original version, but thanks to Activity, React can hydrate the tab buttons first, before it even mounts `Home` or `Video`.

---

Thus, in addition to hiding and showing content, Activity boundaries help improve your app's performance during hydration by letting React know which parts of your page can become interactive in isolation.

And even if your page doesn't ever hide part of its content, you can still add always-visible Activity boundaries to improve hydration performance:

```jsx
function Page() {
 return (
 <>
 <Post />

 <Activity>
 <Comments />
 </Activity>
 </>
 );
}
```

---

## Troubleshooting {/*troubleshooting*/}

### My hidden components have unwanted side effects {/*my-hidden-components-have-unwanted-side-effects*/}

An Activity boundary hides its content by setting `display: none` on its children and cleaning up any of their Effects. So, most well-behaved React components that properly clean up their side effects will already be robust to being hidden by Activity.

But there _are_ some situations where a hidden component behaves differently than an unmounted one. Most notably, since a hidden component's DOM is not destroyed, any side effects from that DOM will persist, even after the component is hidden.

As an example, consider a `<video>` tag. Typically it doesn't require any cleanup, because even if you're playing a video, unmounting the tag stops the video and audio from playing in the browser. Try playing the video and then pressing Home in this demo:

<Sandpack>

```js src/App.js active
import { useState } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Video from './Video.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('video');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'video'}
 onClick={() => setActiveTab('video')}
 >
 Video
 </TabButton>

 <hr />

 {activeTab === 'home' && <Home />}
 {activeTab === 'video' && <Video />}
 </>
 );
}
```

```js src/TabButton.js hidden
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Video.js
export default function Video() {
 return (
 <video
 // 'Big Buck Bunny' licensed under CC 3.0 by the Blender foundation. Hosted by archive.org
 src="https://archive.org/download/BigBuckBunny_124/Content/big_buck_bunny_720p_surround.mp4"
 controls
 playsInline
 />

 );
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
video { width: 300px; margin-top: 10px; aspect-ratio: 16/9; }
```

</Sandpack>

The video stops playing as expected.

Now, let's say we wanted to preserve the timecode where the user last watched, so that when they tab back to the video, it doesn't start over from the beginning again.

This is a great use case for Activity!

Let's update `App` to hide the inactive tab with a hidden Activity boundary instead of unmounting it, and see how the demo behaves this time:

<Sandpack>

```js src/App.js active
import { Activity, useState } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Video from './Video.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('video');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'video'}
 onClick={() => setActiveTab('video')}
 >
 Video
 </TabButton>

 <hr />

 <Activity mode={activeTab === 'home' ? 'visible' : 'hidden'}>
 <Home />
 </Activity>
 <Activity mode={activeTab === 'video' ? 'visible' : 'hidden'}>
 <Video />
 </Activity>
 </>
 );
}
```

```js src/TabButton.js hidden
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Video.js
export default function Video() {
 return (
 <video
 controls
 playsInline
 // 'Big Buck Bunny' licensed under CC 3.0 by the Blender foundation. Hosted by archive.org
 src="https://archive.org/download/BigBuckBunny_124/Content/big_buck_bunny_720p_surround.mp4"
 />

 );
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
video { width: 300px; margin-top: 10px; aspect-ratio: 16/9; }
```

</Sandpack>

Whoops! The video and audio continue to play even after it's been hidden, because the tab's `<video>` element is still in the DOM.

To fix this, we can add an Effect with a cleanup function that pauses the video:

```jsx {2,4-10,14}
export default function VideoTab() {
 const ref = useRef();

 useLayoutEffect(() => {
 const videoRef = ref.current;

 return () => {
 videoRef.pause()
 }
 }, []);

 return (
 <video
 ref={ref}
 controls
 playsInline
 src="..."
 />

 );
}
```

We call `useLayoutEffect` instead of `useEffect` because conceptually the clean-up code is tied to the component's UI being visually hidden. If we used a regular effect, the code could be delayed by (say) a re-suspending Suspense boundary or a View Transition.

Let's see the new behavior. Try playing the video, switching to the Home tab, then back to the Video tab:

<Sandpack>

```js src/App.js active
import { Activity, useState } from 'react';
import TabButton from './TabButton.js';
import Home from './Home.js';
import Video from './Video.js';

export default function App() {
 const [activeTab, setActiveTab] = useState('video');

 return (
 <>
 <TabButton
 isActive={activeTab === 'home'}
 onClick={() => setActiveTab('home')}
 >
 Home
 </TabButton>
 <TabButton
 isActive={activeTab === 'video'}
 onClick={() => setActiveTab('video')}
 >
 Video
 </TabButton>

 <hr />

 <Activity mode={activeTab === 'home' ? 'visible' : 'hidden'}>
 <Home />
 </Activity>
 <Activity mode={activeTab === 'video' ? 'visible' : 'hidden'}>
 <Video />
 </Activity>
 </>
 );
}
```

```js src/TabButton.js hidden
export default function TabButton({ onClick, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }

 return (
 <button onClick={onClick}>
 {children}
 </button>
 );
}
```

```js src/Home.js
export default function Home() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/Video.js
import { useRef, useLayoutEffect } from 'react';

export default function Video() {
 const ref = useRef();

 useLayoutEffect(() => {
 const videoRef = ref.current

 return () => {
 videoRef.pause()
 };
 }, [])

 return (
 <video
 ref={ref}
 controls
 playsInline
 // 'Big Buck Bunny' licensed under CC 3.0 by the Blender foundation. Hosted by archive.org
 src="https://archive.org/download/BigBuckBunny_124/Content/big_buck_bunny_720p_surround.mp4"
 />

 );
}
```

```css
body { height: 275px; }
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
video { width: 300px; margin-top: 10px; aspect-ratio: 16/9; }
```

</Sandpack>

It works great! Our cleanup function ensures that the video stops playing if it's ever hidden by an Activity boundary, and even better, because the `<video>` tag is never destroyed, the timecode is preserved, and the video itself doesn't need to be initialized or downloaded again when the user switches back to keep watching it.

This is a great example of using Activity to preserve ephemeral DOM state for parts of the UI that become hidden, but the user is likely to interact with again soon.

---

Our example illustrates that for certain tags like `<video>`, unmounting and hiding have different behavior. If a component renders DOM that has a side effect, and you want to prevent that side effect when an Activity boundary hides it, add an Effect with a return function to clean it up.

The most common cases of this will be from the following tags:

 - `<video>`
 - `<audio>`
 - `<iframe>`

Typically, though, most of your React components should already be robust to being hidden by an Activity boundary. And conceptually, you should think of "hidden" Activities as being unmounted.

To eagerly discover other Effects that don't have proper cleanup, which is important not only for Activity boundaries but for many other behaviors in React, we recommend using [`<StrictMode>`](/reference/react/StrictMode).

---

### My hidden components have Effects that aren't running {/*my-hidden-components-have-effects-that-arent-running*/}

When an `<Activity>` is "hidden", all its children's Effects are cleaned up. Conceptually, the children are unmounted, but React saves their state for later. This is a feature of Activity because it means subscriptions won't be active for hidden parts of the UI, reducing the amount of work needed for hidden content.

If you're relying on an Effect mounting to clean up a component's side effects, refactor the Effect to do the work in the returned cleanup function instead.

To eagerly find problematic Effects, we recommend adding [`<StrictMode>`](/reference/react/StrictMode) which will eagerly perform Activity unmounts and mounts to catch any unexpected side-effects.

---
title: Children
---

<Pitfall>

Using `Children` is uncommon and can lead to fragile code. [See common alternatives.](#alternatives)

</Pitfall>

<Intro>

`Children` lets you manipulate and transform the JSX you received as the [`children` prop.](/learn/passing-props-to-a-component#passing-jsx-as-children)

```js
const mappedChildren = Children.map(children, child =>
 <div className="Row">
 {child}
 </div>
);

```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `Children.count(children)` {/*children-count*/}

Call `Children.count(children)` to count the number of children in the `children` data structure.

```js src/RowList.js active
import { Children } from 'react';

function RowList({ children }) {
 return (
 <>
 <h1>Total rows: {Children.count(children)}</h1>
 ...
 </>
 );
}
```

[See more examples below.](#counting-children)

#### Parameters {/*children-count-parameters*/}

* `children`: The value of the [`children` prop](/learn/passing-props-to-a-component#passing-jsx-as-children) received by your component.

#### Returns {/*children-count-returns*/}

The number of nodes inside these `children`.

#### Caveats {/*children-count-caveats*/}

- Empty nodes (`null`, `undefined`, and Booleans), strings, numbers, and [React elements](/reference/react/createElement) count as individual nodes. Arrays don't count as individual nodes, but their children do. **The traversal does not go deeper than React elements:** they don't get rendered, and their children aren't traversed. [Fragments](/reference/react/Fragment) don't get traversed.

---

### `Children.forEach(children, fn, thisArg?)` {/*children-foreach*/}

Call `Children.forEach(children, fn, thisArg?)` to run some code for each child in the `children` data structure.

```js src/RowList.js active
import { Children } from 'react';

function SeparatorList({ children }) {
 const result = [];
 Children.forEach(children, (child, index) => {
 result.push(child);
 result.push(<hr key={index} />);
 });
 // ...
```

[See more examples below.](#running-some-code-for-each-child)

#### Parameters {/*children-foreach-parameters*/}

* `children`: The value of the [`children` prop](/learn/passing-props-to-a-component#passing-jsx-as-children) received by your component.
* `fn`: The function you want to run for each child, similar to the [array `forEach` method](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/forEach) callback. It will be called with the child as the first argument and its index as the second argument. The index starts at `0` and increments on each call.
* **optional** `thisArg`: The [`this` value](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/this) with which the `fn` function should be called. If omitted, it's `undefined`.

#### Returns {/*children-foreach-returns*/}

`Children.forEach` returns `undefined`.

#### Caveats {/*children-foreach-caveats*/}

- Empty nodes (`null`, `undefined`, and Booleans), strings, numbers, and [React elements](/reference/react/createElement) count as individual nodes. Arrays don't count as individual nodes, but their children do. **The traversal does not go deeper than React elements:** they don't get rendered, and their children aren't traversed. [Fragments](/reference/react/Fragment) don't get traversed.

---

### `Children.map(children, fn, thisArg?)` {/*children-map*/}

Call `Children.map(children, fn, thisArg?)` to map or transform each child in the `children` data structure.

```js src/RowList.js active
import { Children } from 'react';

function RowList({ children }) {
 return (
 <div className="RowList">
 {Children.map(children, child =>
 <div className="Row">
 {child}
 </div>
 )}
 </div>
 );
}
```

[See more examples below.](#transforming-children)

#### Parameters {/*children-map-parameters*/}

* `children`: The value of the [`children` prop](/learn/passing-props-to-a-component#passing-jsx-as-children) received by your component.
* `fn`: The mapping function, similar to the [array `map` method](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/map) callback. It will be called with the child as the first argument and its index as the second argument. The index starts at `0` and increments on each call. You need to return a React node from this function. This may be an empty node (`null`, `undefined`, or a Boolean), a string, a number, a React element, or an array of other React nodes.
* **optional** `thisArg`: The [`this` value](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/this) with which the `fn` function should be called. If omitted, it's `undefined`.

#### Returns {/*children-map-returns*/}

If `children` is `null` or `undefined`, returns the same value.

Otherwise, returns a flat array consisting of the nodes you've returned from the `fn` function. The returned array will contain all nodes you returned except for `null` and `undefined`.

#### Caveats {/*children-map-caveats*/}

- Empty nodes (`null`, `undefined`, and Booleans), strings, numbers, and [React elements](/reference/react/createElement) count as individual nodes. Arrays don't count as individual nodes, but their children do. **The traversal does not go deeper than React elements:** they don't get rendered, and their children aren't traversed. [Fragments](/reference/react/Fragment) don't get traversed.

- If you return an element or an array of elements with keys from `fn`, **the returned elements' keys will be automatically combined with the key of the corresponding original item from `children`.** When you return multiple elements from `fn` in an array, their keys only need to be unique locally amongst each other.

---

### `Children.only(children)` {/*children-only*/}

Call `Children.only(children)` to assert that `children` represent a single React element.

```js
function Box({ children }) {
 const element = Children.only(children);
 // ...
```

#### Parameters {/*children-only-parameters*/}

* `children`: The value of the [`children` prop](/learn/passing-props-to-a-component#passing-jsx-as-children) received by your component.

#### Returns {/*children-only-returns*/}

If `children` [is a valid element,](/reference/react/isValidElement) returns that element.

Otherwise, throws an error.

#### Caveats {/*children-only-caveats*/}

- This method always **throws if you pass an array (such as the return value of `Children.map`) as `children`.** In other words, it enforces that `children` is a single React element, not that it's an array with a single element.

---

### `Children.toArray(children)` {/*children-toarray*/}

Call `Children.toArray(children)` to create an array out of the `children` data structure.

```js src/ReversedList.js active
import { Children } from 'react';

export default function ReversedList({ children }) {
 const result = Children.toArray(children);
 result.reverse();
 // ...
```

#### Parameters {/*children-toarray-parameters*/}

* `children`: The value of the [`children` prop](/learn/passing-props-to-a-component#passing-jsx-as-children) received by your component.

#### Returns {/*children-toarray-returns*/}

Returns a flat array of elements in `children`.

#### Caveats {/*children-toarray-caveats*/}

- Empty nodes (`null`, `undefined`, and Booleans) will be omitted in the returned array. **The returned elements' keys will be calculated from the original elements' keys and their level of nesting and position.** This ensures that flattening the array does not introduce changes in behavior.

---

## Usage {/*usage*/}

### Transforming children {/*transforming-children*/}

To transform the children JSX that your component [receives as the `children` prop,](/learn/passing-props-to-a-component#passing-jsx-as-children) call `Children.map`:

```js {6,10}
import { Children } from 'react';

function RowList({ children }) {
 return (
 <div className="RowList">
 {Children.map(children, child =>
 <div className="Row">
 {child}
 </div>
 )}
 </div>
 );
}
```

In the example above, the `RowList` wraps every child it receives into a `<div className="Row">` container. For example, let's say the parent component passes three `<p>` tags as the `children` prop to `RowList`:

```js
<RowList>
 <p>This is the first item.</p>
 <p>This is the second item.</p>
 <p>This is the third item.</p>
</RowList>
```

Then, with the `RowList` implementation above, the final rendered result will look like this:

```js
<div className="RowList">
 <div className="Row">
 <p>This is the first item.</p>
 </div>
 <div className="Row">
 <p>This is the second item.</p>
 </div>
 <div className="Row">
 <p>This is the third item.</p>
 </div>
</div>
```

`Children.map` is similar to [to transforming arrays with `map()`.](/learn/rendering-lists) The difference is that the `children` data structure is considered *opaque.* This means that even if it's sometimes an array, you should not assume it's an array or any other particular data type. This is why you should use `Children.map` if you need to transform it.

<Sandpack>

```js
import RowList from './RowList.js';

export default function App() {
 return (
 <RowList>
 <p>This is the first item.</p>
 <p>This is the second item.</p>
 <p>This is the third item.</p>
 </RowList>
 );
}
```

```js src/RowList.js active
import { Children } from 'react';

export default function RowList({ children }) {
 return (
 <div className="RowList">
 {Children.map(children, child =>
 <div className="Row">
 {child}
 </div>
 )}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}
```

</Sandpack>

<DeepDive>

#### Why is the children prop not always an array? {/*why-is-the-children-prop-not-always-an-array*/}

In React, the `children` prop is considered an *opaque* data structure. This means that you shouldn't rely on how it is structured. To transform, filter, or count children, you should use the `Children` methods.

In practice, the `children` data structure is often represented as an array internally. However, if there is only a single child, then React won't create an extra array since this would lead to unnecessary memory overhead. As long as you use the `Children` methods instead of directly introspecting the `children` prop, your code will not break even if React changes how the data structure is actually implemented.

Even when `children` is an array, `Children.map` has useful special behavior. For example, `Children.map` combines the [keys](/learn/rendering-lists#keeping-list-items-in-order-with-key) on the returned elements with the keys on the `children` you've passed to it. This ensures the original JSX children don't "lose" keys even if they get wrapped like in the example above.

</DeepDive>

<Pitfall>

The `children` data structure **does not include rendered output** of the components you pass as JSX. In the example below, the `children` received by the `RowList` only contains two items rather than three:

1. `<p>This is the first item.</p>`
2. `<MoreRows />`

This is why only two row wrappers are generated in this example:

<Sandpack>

```js
import RowList from './RowList.js';

export default function App() {
 return (
 <RowList>
 <p>This is the first item.</p>
 <MoreRows />
 </RowList>
 );
}

function MoreRows() {
 return (
 <>
 <p>This is the second item.</p>
 <p>This is the third item.</p>
 </>
 );
}
```

```js src/RowList.js
import { Children } from 'react';

export default function RowList({ children }) {
 return (
 <div className="RowList">
 {Children.map(children, child =>
 <div className="Row">
 {child}
 </div>
 )}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}
```

</Sandpack>

**There is no way to get the rendered output of an inner component** like `<MoreRows />` when manipulating `children`. This is why [it's usually better to use one of the alternative solutions.](#alternatives)

</Pitfall>

---

### Running some code for each child {/*running-some-code-for-each-child*/}

Call `Children.forEach` to iterate over each child in the `children` data structure. It does not return any value and is similar to the [array `forEach` method.](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/forEach) You can use it to run custom logic like constructing your own array.

<Sandpack>

```js
import SeparatorList from './SeparatorList.js';

export default function App() {
 return (
 <SeparatorList>
 <p>This is the first item.</p>
 <p>This is the second item.</p>
 <p>This is the third item.</p>
 </SeparatorList>
 );
}
```

```js src/SeparatorList.js active
import { Children } from 'react';

export default function SeparatorList({ children }) {
 const result = [];
 Children.forEach(children, (child, index) => {
 result.push(child);
 result.push(<hr key={index} />);
 });
 result.pop(); // Remove the last separator
 return result;
}
```

</Sandpack>

<Pitfall>

As mentioned earlier, there is no way to get the rendered output of an inner component when manipulating `children`. This is why [it's usually better to use one of the alternative solutions.](#alternatives)

</Pitfall>

---

### Counting children {/*counting-children*/}

Call `Children.count(children)` to calculate the number of children.

<Sandpack>

```js
import RowList from './RowList.js';

export default function App() {
 return (
 <RowList>
 <p>This is the first item.</p>
 <p>This is the second item.</p>
 <p>This is the third item.</p>
 </RowList>
 );
}
```

```js src/RowList.js active
import { Children } from 'react';

export default function RowList({ children }) {
 return (
 <div className="RowList">
 <h1 className="RowListHeader">
 Total rows: {Children.count(children)}
 </h1>
 {Children.map(children, child =>
 <div className="Row">
 {child}
 </div>
 )}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.RowListHeader {
 padding-top: 5px;
 font-size: 25px;
 font-weight: bold;
 text-align: center;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}
```

</Sandpack>

<Pitfall>

As mentioned earlier, there is no way to get the rendered output of an inner component when manipulating `children`. This is why [it's usually better to use one of the alternative solutions.](#alternatives)

</Pitfall>

---

### Converting children to an array {/*converting-children-to-an-array*/}

Call `Children.toArray(children)` to turn the `children` data structure into a regular JavaScript array. This lets you manipulate the array with built-in array methods like [`filter`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/filter), [`sort`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/sort), or [`reverse`.](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/reverse)

<Sandpack>

```js
import ReversedList from './ReversedList.js';

export default function App() {
 return (
 <ReversedList>
 <p>This is the first item.</p>
 <p>This is the second item.</p>
 <p>This is the third item.</p>
 </ReversedList>
 );
}
```

```js src/ReversedList.js active
import { Children } from 'react';

export default function ReversedList({ children }) {
 const result = Children.toArray(children);
 result.reverse();
 return result;
}
```

</Sandpack>

<Pitfall>

As mentioned earlier, there is no way to get the rendered output of an inner component when manipulating `children`. This is why [it's usually better to use one of the alternative solutions.](#alternatives)

</Pitfall>

---

## Alternatives {/*alternatives*/}

<Note>

This section describes alternatives to the `Children` API (with capital `C`) that's imported like this:

```js
import { Children } from 'react';
```

Don't confuse it with [using the `children` prop](/learn/passing-props-to-a-component#passing-jsx-as-children) (lowercase `c`), which is good and encouraged.

</Note>

### Exposing multiple components {/*exposing-multiple-components*/}

Manipulating children with the `Children` methods often leads to fragile code. When you pass children to a component in JSX, you don't usually expect the component to manipulate or transform the individual children.

When you can, try to avoid using the `Children` methods. For example, if you want every child of `RowList` to be wrapped in `<div className="Row">`, export a `Row` component, and manually wrap every row into it like this:

<Sandpack>

```js
import { RowList, Row } from './RowList.js';

export default function App() {
 return (
 <RowList>
 <Row>
 <p>This is the first item.</p>
 </Row>
 <Row>
 <p>This is the second item.</p>
 </Row>
 <Row>
 <p>This is the third item.</p>
 </Row>
 </RowList>
 );
}
```

```js src/RowList.js
export function RowList({ children }) {
 return (
 <div className="RowList">
 {children}
 </div>
 );
}

export function Row({ children }) {
 return (
 <div className="Row">
 {children}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}
```

</Sandpack>

Unlike using `Children.map`, this approach does not wrap every child automatically. **However, this approach has a significant benefit compared to the [earlier example with `Children.map`](#transforming-children) because it works even if you keep extracting more components.** For example, it still works if you extract your own `MoreRows` component:

<Sandpack>

```js
import { RowList, Row } from './RowList.js';

export default function App() {
 return (
 <RowList>
 <Row>
 <p>This is the first item.</p>
 </Row>
 <MoreRows />
 </RowList>
 );
}

function MoreRows() {
 return (
 <>
 <Row>
 <p>This is the second item.</p>
 </Row>
 <Row>
 <p>This is the third item.</p>
 </Row>
 </>
 );
}
```

```js src/RowList.js
export function RowList({ children }) {
 return (
 <div className="RowList">
 {children}
 </div>
 );
}

export function Row({ children }) {
 return (
 <div className="Row">
 {children}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}
```

</Sandpack>

This wouldn't work with `Children.map` because it would "see" `<MoreRows />` as a single child (and a single row).

---

### Accepting an array of objects as a prop {/*accepting-an-array-of-objects-as-a-prop*/}

You can also explicitly pass an array as a prop. For example, this `RowList` accepts a `rows` array as a prop:

<Sandpack>

```js
import { RowList, Row } from './RowList.js';

export default function App() {
 return (
 <RowList rows={[
 { id: 'first', content: <p>This is the first item.</p> },
 { id: 'second', content: <p>This is the second item.</p> },
 { id: 'third', content: <p>This is the third item.</p> }
 ]} />
 );
}
```

```js src/RowList.js
export function RowList({ rows }) {
 return (
 <div className="RowList">
 {rows.map(row => (
 <div className="Row" key={row.id}>
 {row.content}
 </div>
 ))}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}
```

</Sandpack>

Since `rows` is a regular JavaScript array, the `RowList` component can use built-in array methods like [`map`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Array/map) on it.

This pattern is especially useful when you want to be able to pass more information as structured data together with children. In the below example, the `TabSwitcher` component receives an array of objects as the `tabs` prop:

<Sandpack>

```js
import TabSwitcher from './TabSwitcher.js';

export default function App() {
 return (
 <TabSwitcher tabs={[
 {
 id: 'first',
 header: 'First',
 content: <p>This is the first item.</p>
 },
 {
 id: 'second',
 header: 'Second',
 content: <p>This is the second item.</p>
 },
 {
 id: 'third',
 header: 'Third',
 content: <p>This is the third item.</p>
 }
 ]} />
 );
}
```

```js src/TabSwitcher.js
import { useState } from 'react';

export default function TabSwitcher({ tabs }) {
 const [selectedId, setSelectedId] = useState(tabs[0].id);
 const selectedTab = tabs.find(tab => tab.id === selectedId);
 return (
 <>
 {tabs.map(tab => (
 <button
 key={tab.id}
 onClick={() => setSelectedId(tab.id)}
 >
 {tab.header}
 </button>
 ))}
 <hr />
 <div key={selectedId}>
 <h3>{selectedTab.header}</h3>
 {selectedTab.content}
 </div>
 </>
 );
}
```

</Sandpack>

Unlike passing the children as JSX, this approach lets you associate some extra data like `header` with each item. Because you are working with the `tabs` directly, and it is an array, you do not need the `Children` methods.

---

### Calling a render prop to customize rendering {/*calling-a-render-prop-to-customize-rendering*/}

Instead of producing JSX for every single item, you can also pass a function that returns JSX, and call that function when necessary. In this example, the `App` component passes a `renderContent` function to the `TabSwitcher` component. The `TabSwitcher` component calls `renderContent` only for the selected tab:

<Sandpack>

```js
import TabSwitcher from './TabSwitcher.js';

export default function App() {
 return (
 <TabSwitcher
 tabIds={['first', 'second', 'third']}
 getHeader={tabId => {
 return tabId[0].toUpperCase() + tabId.slice(1);
 }}
 renderContent={tabId => {
 return <p>This is the {tabId} item.</p>;
 }}
 />
 );
}
```

```js src/TabSwitcher.js
import { useState } from 'react';

export default function TabSwitcher({ tabIds, getHeader, renderContent }) {
 const [selectedId, setSelectedId] = useState(tabIds[0]);
 return (
 <>
 {tabIds.map((tabId) => (
 <button
 key={tabId}
 onClick={() => setSelectedId(tabId)}
 >
 {getHeader(tabId)}
 </button>
 ))}
 <hr />
 <div key={selectedId}>
 <h3>{getHeader(selectedId)}</h3>
 {renderContent(selectedId)}
 </div>
 </>
 );
}
```

</Sandpack>

A prop like `renderContent` is called a *render prop* because it is a prop that specifies how to render a piece of the user interface. However, there is nothing special about it: it is a regular prop which happens to be a function.

Render props are functions, so you can pass information to them. For example, this `RowList` component passes the `id` and the `index` of each row to the `renderRow` render prop, which uses `index` to highlight even rows:

<Sandpack>

```js
import { RowList, Row } from './RowList.js';

export default function App() {
 return (
 <RowList
 rowIds={['first', 'second', 'third']}
 renderRow={(id, index) => {
 return (
 <Row isHighlighted={index % 2 === 0}>
 <p>This is the {id} item.</p>
 </Row>
 );
 }}
 />
 );
}
```

```js src/RowList.js
import { Fragment } from 'react';

export function RowList({ rowIds, renderRow }) {
 return (
 <div className="RowList">
 <h1 className="RowListHeader">
 Total rows: {rowIds.length}
 </h1>
 {rowIds.map((rowId, index) =>
 <Fragment key={rowId}>
 {renderRow(rowId, index)}
 </Fragment>
 )}
 </div>
 );
}

export function Row({ children, isHighlighted }) {
 return (
 <div className={[
 'Row',
 isHighlighted ? 'RowHighlighted' : ''
 ].join(' ')}>
 {children}
 </div>
 );
}
```

```css
.RowList {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.RowListHeader {
 padding-top: 5px;
 font-size: 25px;
 font-weight: bold;
 text-align: center;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}

.RowHighlighted {
 background: #ffa;
}
```

</Sandpack>

This is another example of how parent and child components can cooperate without manipulating the children.

---

## Troubleshooting {/*troubleshooting*/}

### I pass a custom component, but the `Children` methods don't show its render result {/*i-pass-a-custom-component-but-the-children-methods-dont-show-its-render-result*/}

Suppose you pass two children to `RowList` like this:

```js
<RowList>
 <p>First item</p>
 <MoreRows />
</RowList>
```

If you do `Children.count(children)` inside `RowList`, you will get `2`. Even if `MoreRows` renders 10 different items, or if it returns `null`, `Children.count(children)` will still be `2`. From the `RowList`'s perspective, it only "sees" the JSX it has received. It does not "see" the internals of the `MoreRows` component.

The limitation makes it hard to extract a component. This is why [alternatives](#alternatives) are preferred to using `Children`.

---
title: Component
---

<Pitfall>

We recommend defining components as functions instead of classes. [See how to migrate.](#alternatives)

</Pitfall>

<Intro>

`Component` is the base class for the React components defined as [JavaScript classes.](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes) Class components are still supported by React, but we don't recommend using them in new code.

```js
class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `Component` {/*component*/}

To define a React component as a class, extend the built-in `Component` class and define a [`render` method:](#render)

```js
import { Component } from 'react';

class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

Only the `render` method is required, other methods are optional.

[See more examples below.](#usage)

---

### `context` {/*context*/}

The [context](/learn/passing-data-deeply-with-context) of a class component is available as `this.context`. It is only available if you specify *which* context you want to receive using [`static contextType`](#static-contexttype).

A class component can only read one context at a time.

```js {2,5}
class Button extends Component {
 static contextType = ThemeContext;

 render() {
 const theme = this.context;
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {this.props.children}
 </button>
 );
 }
}

```

<Note>

Reading `this.context` in class components is equivalent to [`useContext`](/reference/react/useContext) in function components.

[See how to migrate.](#migrating-a-component-with-context-from-a-class-to-a-function)

</Note>

---

### `props` {/*props*/}

The props passed to a class component are available as `this.props`.

```js {3}
class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}

<Greeting name="Taylor" />
```

<Note>

Reading `this.props` in class components is equivalent to [declaring props](/learn/passing-props-to-a-component#step-2-read-props-inside-the-child-component) in function components.

[See how to migrate.](#migrating-a-simple-component-from-a-class-to-a-function)

</Note>

---

### `state` {/*state*/}

The state of a class component is available as `this.state`. The `state` field must be an object. Do not mutate the state directly. If you wish to change the state, call `setState` with the new state.

```js {2-4,7-9,18}
class Counter extends Component {
 state = {
 age: 42,
 };

 handleAgeChange = () => {
 this.setState({
 age: this.state.age + 1
 });
 };

 render() {
 return (
 <>
 <button onClick={this.handleAgeChange}>
 Increment age
 </button>
 <p>You are {this.state.age}.</p>
 </>
 );
 }
}
```

<Note>

Defining `state` in class components is equivalent to calling [`useState`](/reference/react/useState) in function components.

[See how to migrate.](#migrating-a-component-with-state-from-a-class-to-a-function)

</Note>

---

### `constructor(props)` {/*constructor*/}

The [constructor](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes/constructor) runs before your class component *mounts* (gets added to the screen). Typically, a constructor is only used for two purposes in React. It lets you declare state and [bind](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_objects/Function/bind) your class methods to the class instance:

```js {2-6}
class Counter extends Component {
 constructor(props) {
 super(props);
 this.state = { counter: 0 };
 this.handleClick = this.handleClick.bind(this);
 }

 handleClick() {
 // ...
 }
```

If you use modern JavaScript syntax, constructors are rarely needed. Instead, you can rewrite this code above using the [public class field syntax](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Classes/Public_class_fields) which is supported both by modern browsers and tools like [Babel:](https://babeljs.io/)

```js {2,4}
class Counter extends Component {
 state = { counter: 0 };

 handleClick = () => {
 // ...
 }
```

A constructor should not contain any side effects or subscriptions.

#### Parameters {/*constructor-parameters*/}

* `props`: The component's initial props.

#### Returns {/*constructor-returns*/}

`constructor` should not return anything.

#### Caveats {/*constructor-caveats*/}

* Do not run any side effects or subscriptions in the constructor. Instead, use [`componentDidMount`](#componentdidmount) for that.

* Inside a constructor, you need to call `super(props)` before any other statement. If you don't do that, `this.props` will be `undefined` while the constructor runs, which can be confusing and cause bugs.

* Constructor is the only place where you can assign [`this.state`](#state) directly. In all other methods, you need to use [`this.setState()`](#setstate) instead. Do not call `setState` in the constructor.

* When you use [server rendering,](/reference/react-dom/server) the constructor will run on the server too, followed by the [`render`](#render) method. However, lifecycle methods like `componentDidMount` or `componentWillUnmount` will not run on the server.

* When [Strict Mode](/reference/react/StrictMode) is on, React will call `constructor` twice in development and then throw away one of the instances. This helps you notice the accidental side effects that need to be moved out of the `constructor`.

<Note>

There is no exact equivalent for `constructor` in function components. To declare state in a function component, call [`useState`.](/reference/react/useState) To avoid recalculating the initial state, [pass a function to `useState`.](/reference/react/useState#avoiding-recreating-the-initial-state)

</Note>

---

### `componentDidCatch(error, info)` {/*componentdidcatch*/}

If you define `componentDidCatch`, React will call it when some child component (including distant children) throws an error during rendering. This lets you log that error to an error reporting service in production.

Typically, it is used together with [`static getDerivedStateFromError`](#static-getderivedstatefromerror) which lets you update state in response to an error and display an error message to the user. A component with these methods is called an *Error Boundary*.

[See an example.](#catching-rendering-errors-with-an-error-boundary)

#### Parameters {/*componentdidcatch-parameters*/}

* `error`: The error that was thrown. In practice, it will usually be an instance of [`Error`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error) but this is not guaranteed because JavaScript allows to [`throw`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/throw) any value, including strings or even `null`.

* `info`: An object containing additional information about the error. Its `componentStack` field contains a stack trace with the component that threw, as well as the names and source locations of all its parent components. In production, the component names will be minified. If you set up production error reporting, you can decode the component stack using sourcemaps the same way as you would do for regular JavaScript error stacks.

#### Returns {/*componentdidcatch-returns*/}

`componentDidCatch` should not return anything.

#### Caveats {/*componentdidcatch-caveats*/}

* In the past, it was common to call `setState` inside `componentDidCatch` in order to update the UI and display the fallback error message. This is deprecated in favor of defining [`static getDerivedStateFromError`.](#static-getderivedstatefromerror)

* Production and development builds of React slightly differ in the way `componentDidCatch` handles errors. In development, the errors will bubble up to `window`, which means that any `window.onerror` or `window.addEventListener('error', callback)` will intercept the errors that have been caught by `componentDidCatch`. In production, instead, the errors will not bubble up, which means any ancestor error handler will only receive errors not explicitly caught by `componentDidCatch`.

<Note>

There is no direct equivalent for `componentDidCatch` in function components yet. If you'd like to avoid creating class components, write a single `ErrorBoundary` component like above and use it throughout your app. Alternatively, you can use the [`react-error-boundary`](https://github.com/bvaughn/react-error-boundary) package which does that for you.

</Note>

---

### `componentDidMount()` {/*componentdidmount*/}

If you define the `componentDidMount` method, React will call it when your component is added *(mounted)* to the screen. This is a common place to start data fetching, set up subscriptions, or manipulate the DOM nodes.

If you implement `componentDidMount`, you usually need to implement other lifecycle methods to avoid bugs. For example, if `componentDidMount` reads some state or props, you also have to implement [`componentDidUpdate`](#componentdidupdate) to handle their changes, and [`componentWillUnmount`](#componentwillunmount) to clean up whatever `componentDidMount` was doing.

```js {6-8}
class ChatRoom extends Component {
 state = {
 serverUrl: 'https://localhost:1234'
 };

 componentDidMount() {
 this.setupConnection();
 }

 componentDidUpdate(prevProps, prevState) {
 if (
 this.props.roomId !== prevProps.roomId ||
 this.state.serverUrl !== prevState.serverUrl
 ) {
 this.destroyConnection();
 this.setupConnection();
 }
 }

 componentWillUnmount() {
 this.destroyConnection();
 }

 // ...
}
```

[See more examples.](#adding-lifecycle-methods-to-a-class-component)

#### Parameters {/*componentdidmount-parameters*/}

`componentDidMount` does not take any parameters.

#### Returns {/*componentdidmount-returns*/}

`componentDidMount` should not return anything.

#### Caveats {/*componentdidmount-caveats*/}

- When [Strict Mode](/reference/react/StrictMode) is on, in development React will call `componentDidMount`, then immediately call [`componentWillUnmount`,](#componentwillunmount) and then call `componentDidMount` again. This helps you notice if you forgot to implement `componentWillUnmount` or if its logic doesn't fully "mirror" what `componentDidMount` does.

- Although you may call [`setState`](#setstate) immediately in `componentDidMount`, it's best to avoid that when you can. It will trigger an extra rendering, but it will happen before the browser updates the screen. This guarantees that even though the [`render`](#render) will be called twice in this case, the user won't see the intermediate state. Use this pattern with caution because it often causes performance issues. In most cases, you should be able to assign the initial state in the [`constructor`](#constructor) instead. It can, however, be necessary for cases like modals and tooltips when you need to measure a DOM node before rendering something that depends on its size or position.

<Note>

For many use cases, defining `componentDidMount`, `componentDidUpdate`, and `componentWillUnmount` together in class components is equivalent to calling [`useEffect`](/reference/react/useEffect) in function components. In the rare cases where it's important for the code to run before browser paint, [`useLayoutEffect`](/reference/react/useLayoutEffect) is a closer match.

[See how to migrate.](#migrating-a-component-with-lifecycle-methods-from-a-class-to-a-function)

</Note>

---

### `componentDidUpdate(prevProps, prevState, snapshot?)` {/*componentdidupdate*/}

If you define the `componentDidUpdate` method, React will call it immediately after your component has been re-rendered with updated props or state. This method is not called for the initial render.

You can use it to manipulate the DOM after an update. This is also a common place to do network requests as long as you compare the current props to previous props (e.g. a network request may not be necessary if the props have not changed). Typically, you'd use it together with [`componentDidMount`](#componentdidmount) and [`componentWillUnmount`:](#componentwillunmount)

```js {10-18}
class ChatRoom extends Component {
 state = {
 serverUrl: 'https://localhost:1234'
 };

 componentDidMount() {
 this.setupConnection();
 }

 componentDidUpdate(prevProps, prevState) {
 if (
 this.props.roomId !== prevProps.roomId ||
 this.state.serverUrl !== prevState.serverUrl
 ) {
 this.destroyConnection();
 this.setupConnection();
 }
 }

 componentWillUnmount() {
 this.destroyConnection();
 }

 // ...
}
```

[See more examples.](#adding-lifecycle-methods-to-a-class-component)

#### Parameters {/*componentdidupdate-parameters*/}

* `prevProps`: Props before the update. Compare `prevProps` to [`this.props`](#props) to determine what changed.

* `prevState`: State before the update. Compare `prevState` to [`this.state`](#state) to determine what changed.

* `snapshot`: If you implemented [`getSnapshotBeforeUpdate`](#getsnapshotbeforeupdate), `snapshot` will contain the value you returned from that method. Otherwise, it will be `undefined`.

#### Returns {/*componentdidupdate-returns*/}

`componentDidUpdate` should not return anything.

#### Caveats {/*componentdidupdate-caveats*/}

- `componentDidUpdate` will not get called if [`shouldComponentUpdate`](#shouldcomponentupdate) is defined and returns `false`.

- The logic inside `componentDidUpdate` should usually be wrapped in conditions comparing `this.props` with `prevProps`, and `this.state` with `prevState`. Otherwise, there's a risk of creating infinite loops.

- Although you may call [`setState`](#setstate) immediately in `componentDidUpdate`, it's best to avoid that when you can. It will trigger an extra rendering, but it will happen before the browser updates the screen. This guarantees that even though the [`render`](#render) will be called twice in this case, the user won't see the intermediate state. This pattern often causes performance issues, but it may be necessary for rare cases like modals and tooltips when you need to measure a DOM node before rendering something that depends on its size or position.

<Note>

For many use cases, defining `componentDidMount`, `componentDidUpdate`, and `componentWillUnmount` together in class components is equivalent to calling [`useEffect`](/reference/react/useEffect) in function components. In the rare cases where it's important for the code to run before browser paint, [`useLayoutEffect`](/reference/react/useLayoutEffect) is a closer match.

[See how to migrate.](#migrating-a-component-with-lifecycle-methods-from-a-class-to-a-function)

</Note>
---

### `componentWillMount()` {/*componentwillmount*/}

<Deprecated>

This API has been renamed from `componentWillMount` to [`UNSAFE_componentWillMount`.](#unsafe_componentwillmount) The old name has been deprecated. In a future major version of React, only the new name will work.

Run the [`rename-unsafe-lifecycles` codemod](https://github.com/reactjs/react-codemod#rename-unsafe-lifecycles) to automatically update your components.

</Deprecated>

---

### `componentWillReceiveProps(nextProps)` {/*componentwillreceiveprops*/}

<Deprecated>

This API has been renamed from `componentWillReceiveProps` to [`UNSAFE_componentWillReceiveProps`.](#unsafe_componentwillreceiveprops) The old name has been deprecated. In a future major version of React, only the new name will work.

Run the [`rename-unsafe-lifecycles` codemod](https://github.com/reactjs/react-codemod#rename-unsafe-lifecycles) to automatically update your components.

</Deprecated>

---

### `componentWillUpdate(nextProps, nextState)` {/*componentwillupdate*/}

<Deprecated>

This API has been renamed from `componentWillUpdate` to [`UNSAFE_componentWillUpdate`.](#unsafe_componentwillupdate) The old name has been deprecated. In a future major version of React, only the new name will work.

Run the [`rename-unsafe-lifecycles` codemod](https://github.com/reactjs/react-codemod#rename-unsafe-lifecycles) to automatically update your components.

</Deprecated>

---

### `componentWillUnmount()` {/*componentwillunmount*/}

If you define the `componentWillUnmount` method, React will call it before your component is removed *(unmounted)* from the screen. This is a common place to cancel data fetching or remove subscriptions.

The logic inside `componentWillUnmount` should "mirror" the logic inside [`componentDidMount`.](#componentdidmount) For example, if `componentDidMount` sets up a subscription, `componentWillUnmount` should clean up that subscription. If the cleanup logic in your `componentWillUnmount` reads some props or state, you will usually also need to implement [`componentDidUpdate`](#componentdidupdate) to clean up resources (such as subscriptions) corresponding to the old props and state.

```js {20-22}
class ChatRoom extends Component {
 state = {
 serverUrl: 'https://localhost:1234'
 };

 componentDidMount() {
 this.setupConnection();
 }

 componentDidUpdate(prevProps, prevState) {
 if (
 this.props.roomId !== prevProps.roomId ||
 this.state.serverUrl !== prevState.serverUrl
 ) {
 this.destroyConnection();
 this.setupConnection();
 }
 }

 componentWillUnmount() {
 this.destroyConnection();
 }

 // ...
}
```

[See more examples.](#adding-lifecycle-methods-to-a-class-component)

#### Parameters {/*componentwillunmount-parameters*/}

`componentWillUnmount` does not take any parameters.

#### Returns {/*componentwillunmount-returns*/}

`componentWillUnmount` should not return anything.

#### Caveats {/*componentwillunmount-caveats*/}

- When [Strict Mode](/reference/react/StrictMode) is on, in development React will call [`componentDidMount`,](#componentdidmount) then immediately call `componentWillUnmount`, and then call `componentDidMount` again. This helps you notice if you forgot to implement `componentWillUnmount` or if its logic doesn't fully "mirror" what `componentDidMount` does.

<Note>

For many use cases, defining `componentDidMount`, `componentDidUpdate`, and `componentWillUnmount` together in class components is equivalent to calling [`useEffect`](/reference/react/useEffect) in function components. In the rare cases where it's important for the code to run before browser paint, [`useLayoutEffect`](/reference/react/useLayoutEffect) is a closer match.

[See how to migrate.](#migrating-a-component-with-lifecycle-methods-from-a-class-to-a-function)

</Note>

---

### `forceUpdate(callback?)` {/*forceupdate*/}

Forces a component to re-render.

Usually, this is not necessary. If your component's [`render`](#render) method only reads from [`this.props`](#props), [`this.state`](#state), or [`this.context`,](#context) it will re-render automatically when you call [`setState`](#setstate) inside your component or one of its parents. However, if your component's `render` method reads directly from an external data source, you have to tell React to update the user interface when that data source changes. That's what `forceUpdate` lets you do.

Try to avoid all uses of `forceUpdate` and only read from `this.props` and `this.state` in `render`.

#### Parameters {/*forceupdate-parameters*/}

* **optional** `callback` If specified, React will call the `callback` you've provided after the update is committed.

#### Returns {/*forceupdate-returns*/}

`forceUpdate` does not return anything.

#### Caveats {/*forceupdate-caveats*/}

- If you call `forceUpdate`, React will re-render without calling [`shouldComponentUpdate`.](#shouldcomponentupdate)

<Note>

Reading an external data source and forcing class components to re-render in response to its changes with `forceUpdate` has been superseded by [`useSyncExternalStore`](/reference/react/useSyncExternalStore) in function components.

</Note>

---

### `getSnapshotBeforeUpdate(prevProps, prevState)` {/*getsnapshotbeforeupdate*/}

If you implement `getSnapshotBeforeUpdate`, React will call it immediately before React updates the DOM. It enables your component to capture some information from the DOM (e.g. scroll position) before it is potentially changed. Any value returned by this lifecycle method will be passed as a parameter to [`componentDidUpdate`.](#componentdidupdate)

For example, you can use it in a UI like a chat thread that needs to preserve its scroll position during updates:

```js {7-15,17}
class ScrollingList extends React.Component {
 constructor(props) {
 super(props);
 this.listRef = React.createRef();
 }

 getSnapshotBeforeUpdate(prevProps, prevState) {
 // Are we adding new items to the list?
 // Capture the scroll position so we can adjust scroll later.
 if (prevProps.list.length < this.props.list.length) {
 const list = this.listRef.current;
 return list.scrollHeight - list.scrollTop;
 }
 return null;
 }

 componentDidUpdate(prevProps, prevState, snapshot) {
 // If we have a snapshot value, we've just added new items.
 // Adjust scroll so these new items don't push the old ones out of view.
 // (snapshot here is the value returned from getSnapshotBeforeUpdate)
 if (snapshot !== null) {
 const list = this.listRef.current;
 list.scrollTop = list.scrollHeight - snapshot;
 }
 }

 render() {
 return (
 <div ref={this.listRef}>{/* ...contents... */}</div>
 );
 }
}
```

In the above example, it is important to read the `scrollHeight` property directly in `getSnapshotBeforeUpdate`. It is not safe to read it in [`render`](#render), [`UNSAFE_componentWillReceiveProps`](#unsafe_componentwillreceiveprops), or [`UNSAFE_componentWillUpdate`](#unsafe_componentwillupdate) because there is a potential time gap between these methods getting called and React updating the DOM.

#### Parameters {/*getsnapshotbeforeupdate-parameters*/}

* `prevProps`: Props before the update. Compare `prevProps` to [`this.props`](#props) to determine what changed.

* `prevState`: State before the update. Compare `prevState` to [`this.state`](#state) to determine what changed.

#### Returns {/*getsnapshotbeforeupdate-returns*/}

You should return a snapshot value of any type that you'd like, or `null`. The value you returned will be passed as the third argument to [`componentDidUpdate`.](#componentdidupdate)

#### Caveats {/*getsnapshotbeforeupdate-caveats*/}

- `getSnapshotBeforeUpdate` will not get called if [`shouldComponentUpdate`](#shouldcomponentupdate) is defined and returns `false`.

<Note>

At the moment, there is no equivalent to `getSnapshotBeforeUpdate` for function components. This use case is very uncommon, but if you have the need for it, for now you'll have to write a class component.

</Note>

---

### `render()` {/*render*/}

The `render` method is the only required method in a class component.

The `render` method should specify what you want to appear on the screen, for example:

```js {4-6}
import { Component } from 'react';

class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

React may call `render` at any moment, so you shouldn't assume that it runs at a particular time. Usually, the `render` method should return a piece of [JSX](/learn/writing-markup-with-jsx), but a few [other return types](#render-returns) (like strings) are supported. To calculate the returned JSX, the `render` method can read [`this.props`](#props), [`this.state`](#state), and [`this.context`](#context).

You should write the `render` method as a pure function, meaning that it should return the same result if props, state, and context are the same. It also shouldn't contain side effects (like setting up subscriptions) or interact with the browser APIs. Side effects should happen either in event handlers or methods like [`componentDidMount`.](#componentdidmount)

#### Parameters {/*render-parameters*/}

`render` does not take any parameters.

#### Returns {/*render-returns*/}

`render` can return any valid React node. This includes React elements such as `<div />`, strings, numbers, [portals](/reference/react-dom/createPortal), empty nodes (`null`, `undefined`, `true`, and `false`), and arrays of React nodes.

#### Caveats {/*render-caveats*/}

- `render` should be written as a pure function of props, state, and context. It should not have side effects.

- `render` will not get called if [`shouldComponentUpdate`](#shouldcomponentupdate) is defined and returns `false`.

- When [Strict Mode](/reference/react/StrictMode) is on, React will call `render` twice in development and then throw away one of the results. This helps you notice the accidental side effects that need to be moved out of the `render` method.

- There is no one-to-one correspondence between the `render` call and the subsequent `componentDidMount` or `componentDidUpdate` call. Some of the `render` call results may be discarded by React when it's beneficial.

---

### `setState(nextState, callback?)` {/*setstate*/}

Call `setState` to update the state of your React component.

```js {8-10}
class Form extends Component {
 state = {
 name: 'Taylor',
 };

 handleNameChange = (e) => {
 const newName = e.target.value;
 this.setState({
 name: newName
 });
 }

 render() {
 return (
 <>
 <input value={this.state.name} onChange={this.handleNameChange} />
 <p>Hello, {this.state.name}.</p>
 </>
 );
 }
}
```

`setState` enqueues changes to the component state. It tells React that this component and its children need to re-render with the new state. This is the main way you'll update the user interface in response to interactions.

<Pitfall>

Calling `setState` **does not** change the current state in the already executing code:

```js {6}
function handleClick() {
 console.log(this.state.name); // "Taylor"
 this.setState({
 name: 'Robin'
 });
 console.log(this.state.name); // Still "Taylor"!
}
```

It only affects what `this.state` will return starting from the *next* render.

</Pitfall>

You can also pass a function to `setState`. It lets you update state based on the previous state:

```js {2-6}
 handleIncreaseAge = () => {
 this.setState(prevState => {
 return {
 age: prevState.age + 1
 };
 });
 }
```

You don't have to do this, but it's handy if you want to update state multiple times during the same event.

#### Parameters {/*setstate-parameters*/}

* `nextState`: Either an object or a function.
 * If you pass an object as `nextState`, it will be shallowly merged into `this.state`.
 * If you pass a function as `nextState`, it will be treated as an _updater function_. It must be pure, should take the pending state and props as arguments, and should return the object to be shallowly merged into `this.state`. React will put your updater function in a queue and re-render your component. During the next render, React will calculate the next state by applying all of the queued updaters to the previous state.

* **optional** `callback`: If specified, React will call the `callback` you've provided after the update is committed.

#### Returns {/*setstate-returns*/}

`setState` does not return anything.

#### Caveats {/*setstate-caveats*/}

- Think of `setState` as a *request* rather than an immediate command to update the component. When multiple components update their state in response to an event, React will batch their updates and re-render them together in a single pass at the end of the event. In the rare case that you need to force a particular state update to be applied synchronously, you may wrap it in [`flushSync`,](/reference/react-dom/flushSync) but this may hurt performance.

- `setState` does not update `this.state` immediately. This makes reading `this.state` right after calling `setState` a potential pitfall. Instead, use [`componentDidUpdate`](#componentdidupdate) or the setState `callback` argument, either of which are guaranteed to fire after the update has been applied. If you need to set the state based on the previous state, you can pass a function to `nextState` as described above.

<Note>

Calling `setState` in class components is similar to calling a [`set` function](/reference/react/useState#setstate) in function components.

[See how to migrate.](#migrating-a-component-with-state-from-a-class-to-a-function)

</Note>

---

### `shouldComponentUpdate(nextProps, nextState, nextContext)` {/*shouldcomponentupdate*/}

If you define `shouldComponentUpdate`, React will call it to determine whether a re-render can be skipped.

If you are confident you want to write it by hand, you may compare `this.props` with `nextProps` and `this.state` with `nextState` and return `false` to tell React the update can be skipped.

```js {6-18}
class Rectangle extends Component {
 state = {
 isHovered: false
 };

 shouldComponentUpdate(nextProps, nextState) {
 if (
 nextProps.position.x === this.props.position.x &&
 nextProps.position.y === this.props.position.y &&
 nextProps.size.width === this.props.size.width &&
 nextProps.size.height === this.props.size.height &&
 nextState.isHovered === this.state.isHovered
 ) {
 // Nothing has changed, so a re-render is unnecessary
 return false;
 }
 return true;
 }

 // ...
}

```

React calls `shouldComponentUpdate` before rendering when new props or state are being received. Defaults to `true`. This method is not called for the initial render or when [`forceUpdate`](#forceupdate) is used.

#### Parameters {/*shouldcomponentupdate-parameters*/}

- `nextProps`: The next props that the component is about to render with. Compare `nextProps` to [`this.props`](#props) to determine what changed.
- `nextState`: The next state that the component is about to render with. Compare `nextState` to [`this.state`](#props) to determine what changed.
- `nextContext`: The next context that the component is about to render with. Compare `nextContext` to [`this.context`](#context) to determine what changed. Only available if you specify [`static contextType`](#static-contexttype).

#### Returns {/*shouldcomponentupdate-returns*/}

Return `true` if you want the component to re-render. That's the default behavior.

Return `false` to tell React that re-rendering can be skipped.

#### Caveats {/*shouldcomponentupdate-caveats*/}

- This method *only* exists as a performance optimization. If your component breaks without it, fix that first.

- Consider using [`PureComponent`](/reference/react/PureComponent) instead of writing `shouldComponentUpdate` by hand. `PureComponent` shallowly compares props and state, and reduces the chance that you'll skip a necessary update.

- We do not recommend doing deep equality checks or using `JSON.stringify` in `shouldComponentUpdate`. It makes performance unpredictable and dependent on the data structure of every prop and state. In the best case, you risk introducing multi-second stalls to your application, and in the worst case you risk crashing it.

- Returning `false` does not prevent child components from re-rendering when *their* state changes.

- Returning `false` does not *guarantee* that the component will not re-render. React will use the return value as a hint but it may still choose to re-render your component if it makes sense to do for other reasons.

<Note>

Optimizing class components with `shouldComponentUpdate` is similar to optimizing function components with [`memo`.](/reference/react/memo) Function components also offer more granular optimization with [`useMemo`.](/reference/react/useMemo)

</Note>

---

### `UNSAFE_componentWillMount()` {/*unsafe_componentwillmount*/}

If you define `UNSAFE_componentWillMount`, React will call it immediately after the [`constructor`.](#constructor) It only exists for historical reasons and should not be used in any new code. Instead, use one of the alternatives:

- To initialize state, declare [`state`](#state) as a class field or set `this.state` inside the [`constructor`.](#constructor)
- If you need to run a side effect or set up a subscription, move that logic to [`componentDidMount`](#componentdidmount) instead.

[See examples of migrating away from unsafe lifecycles.](https://legacy.reactjs.org/blog/2018/03/27/update-on-async-rendering.html#examples)

#### Parameters {/*unsafe_componentwillmount-parameters*/}

`UNSAFE_componentWillMount` does not take any parameters.

#### Returns {/*unsafe_componentwillmount-returns*/}

`UNSAFE_componentWillMount` should not return anything.

#### Caveats {/*unsafe_componentwillmount-caveats*/}

- `UNSAFE_componentWillMount` will not get called if the component implements [`static getDerivedStateFromProps`](#static-getderivedstatefromprops) or [`getSnapshotBeforeUpdate`.](#getsnapshotbeforeupdate)

- Despite its naming, `UNSAFE_componentWillMount` does not guarantee that the component *will* get mounted if your app uses modern React features like [`Suspense`.](/reference/react/Suspense) If a render attempt is suspended (for example, because the code for some child component has not loaded yet), React will throw the in-progress tree away and attempt to construct the component from scratch during the next attempt. This is why this method is "unsafe". Code that relies on mounting (like adding a subscription) should go into [`componentDidMount`.](#componentdidmount)

- `UNSAFE_componentWillMount` is the only lifecycle method that runs during [server rendering.](/reference/react-dom/server) For all practical purposes, it is identical to [`constructor`,](#constructor) so you should use the `constructor` for this type of logic instead.

<Note>

Calling [`setState`](#setstate) inside `UNSAFE_componentWillMount` in a class component to initialize state is equivalent to passing that state as the initial state to [`useState`](/reference/react/useState) in a function component.

</Note>

---

### `UNSAFE_componentWillReceiveProps(nextProps, nextContext)` {/*unsafe_componentwillreceiveprops*/}

If you define `UNSAFE_componentWillReceiveProps`, React will call it when the component receives new props. It only exists for historical reasons and should not be used in any new code. Instead, use one of the alternatives:

- If you need to **run a side effect** (for example, fetch data, run an animation, or reinitialize a subscription) in response to prop changes, move that logic to [`componentDidUpdate`](#componentdidupdate) instead.
- If you need to **avoid re-computing some data only when a prop changes,** use a [memoization helper](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#what-about-memoization) instead.
- If you need to **"reset" some state when a prop changes,** consider either making a component [fully controlled](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#recommendation-fully-controlled-component) or [fully uncontrolled with a key](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#recommendation-fully-uncontrolled-component-with-a-key) instead.
- If you need to **"adjust" some state when a prop changes,** check whether you can compute all the necessary information from props alone during rendering. If you can't, use [`static getDerivedStateFromProps`](/reference/react/Component#static-getderivedstatefromprops) instead.

[See examples of migrating away from unsafe lifecycles.](https://legacy.reactjs.org/blog/2018/03/27/update-on-async-rendering.html#updating-state-based-on-props)

#### Parameters {/*unsafe_componentwillreceiveprops-parameters*/}

- `nextProps`: The next props that the component is about to receive from its parent component. Compare `nextProps` to [`this.props`](#props) to determine what changed.
- `nextContext`: The next context that the component is about to receive from the closest provider. Compare `nextContext` to [`this.context`](#context) to determine what changed. Only available if you specify [`static contextType`](#static-contexttype).

#### Returns {/*unsafe_componentwillreceiveprops-returns*/}

`UNSAFE_componentWillReceiveProps` should not return anything.

#### Caveats {/*unsafe_componentwillreceiveprops-caveats*/}

- `UNSAFE_componentWillReceiveProps` will not get called if the component implements [`static getDerivedStateFromProps`](#static-getderivedstatefromprops) or [`getSnapshotBeforeUpdate`.](#getsnapshotbeforeupdate)

- Despite its naming, `UNSAFE_componentWillReceiveProps` does not guarantee that the component *will* receive those props if your app uses modern React features like [`Suspense`.](/reference/react/Suspense) If a render attempt is suspended (for example, because the code for some child component has not loaded yet), React will throw the in-progress tree away and attempt to construct the component from scratch during the next attempt. By the time of the next render attempt, the props might be different. This is why this method is "unsafe". Code that should run only for committed updates (like resetting a subscription) should go into [`componentDidUpdate`.](#componentdidupdate)

- `UNSAFE_componentWillReceiveProps` does not mean that the component has received *different* props than the last time. You need to compare `nextProps` and `this.props` yourself to check if something changed.

- React doesn't call `UNSAFE_componentWillReceiveProps` with initial props during mounting. It only calls this method if some of component's props are going to be updated. For example, calling [`setState`](#setstate) doesn't generally trigger `UNSAFE_componentWillReceiveProps` inside the same component.

<Note>

Calling [`setState`](#setstate) inside `UNSAFE_componentWillReceiveProps` in a class component to "adjust" state is equivalent to [calling the `set` function from `useState` during rendering](/reference/react/useState#storing-information-from-previous-renders) in a function component.

</Note>

---

### `UNSAFE_componentWillUpdate(nextProps, nextState)` {/*unsafe_componentwillupdate*/}

If you define `UNSAFE_componentWillUpdate`, React will call it before rendering with the new props or state. It only exists for historical reasons and should not be used in any new code. Instead, use one of the alternatives:

- If you need to run a side effect (for example, fetch data, run an animation, or reinitialize a subscription) in response to prop or state changes, move that logic to [`componentDidUpdate`](#componentdidupdate) instead.
- If you need to read some information from the DOM (for example, to save the current scroll position) so that you can use it in [`componentDidUpdate`](#componentdidupdate) later, read it inside [`getSnapshotBeforeUpdate`](#getsnapshotbeforeupdate) instead.

[See examples of migrating away from unsafe lifecycles.](https://legacy.reactjs.org/blog/2018/03/27/update-on-async-rendering.html#examples)

#### Parameters {/*unsafe_componentwillupdate-parameters*/}

- `nextProps`: The next props that the component is about to render with. Compare `nextProps` to [`this.props`](#props) to determine what changed.
- `nextState`: The next state that the component is about to render with. Compare `nextState` to [`this.state`](#state) to determine what changed.

#### Returns {/*unsafe_componentwillupdate-returns*/}

`UNSAFE_componentWillUpdate` should not return anything.

#### Caveats {/*unsafe_componentwillupdate-caveats*/}

- `UNSAFE_componentWillUpdate` will not get called if [`shouldComponentUpdate`](#shouldcomponentupdate) is defined and returns `false`.

- `UNSAFE_componentWillUpdate` will not get called if the component implements [`static getDerivedStateFromProps`](#static-getderivedstatefromprops) or [`getSnapshotBeforeUpdate`.](#getsnapshotbeforeupdate)

- It's not supported to call [`setState`](#setstate) (or any method that leads to `setState` being called, like dispatching a Redux action) during `componentWillUpdate`.

- Despite its naming, `UNSAFE_componentWillUpdate` does not guarantee that the component *will* update if your app uses modern React features like [`Suspense`.](/reference/react/Suspense) If a render attempt is suspended (for example, because the code for some child component has not loaded yet), React will throw the in-progress tree away and attempt to construct the component from scratch during the next attempt. By the time of the next render attempt, the props and state might be different. This is why this method is "unsafe". Code that should run only for committed updates (like resetting a subscription) should go into [`componentDidUpdate`.](#componentdidupdate)

- `UNSAFE_componentWillUpdate` does not mean that the component has received *different* props or state than the last time. You need to compare `nextProps` with `this.props` and `nextState` with `this.state` yourself to check if something changed.

- React doesn't call `UNSAFE_componentWillUpdate` with initial props and state during mounting.

<Note>

There is no direct equivalent to `UNSAFE_componentWillUpdate` in function components.

</Note>

---

### `static contextType` {/*static-contexttype*/}

If you want to read [`this.context`](#context-instance-field) from your class component, you must specify which context it needs to read. The context you specify as the `static contextType` must be a value previously created by [`createContext`.](/reference/react/createContext)

```js {2}
class Button extends Component {
 static contextType = ThemeContext;

 render() {
 const theme = this.context;
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {this.props.children}
 </button>
 );
 }
}
```

<Note>

Reading `this.context` in class components is equivalent to [`useContext`](/reference/react/useContext) in function components.

[See how to migrate.](#migrating-a-component-with-context-from-a-class-to-a-function)

</Note>

---

### `static defaultProps` {/*static-defaultprops*/}

You can define `static defaultProps` to set the default props for the class. They will be used for `undefined` and missing props, but not for `null` props.

For example, here is how you define that the `color` prop should default to `'blue'`:

```js {2-4}
class Button extends Component {
 static defaultProps = {
 color: 'blue'
 };

 render() {
 return <button className={this.props.color}>click me</button>;
 }
}
```

If the `color` prop is not provided or is `undefined`, it will be set by default to `'blue'`:

```js
<>
 {/* this.props.color is "blue" */}
 <Button />

 {/* this.props.color is "blue" */}
 <Button color={undefined} />

 {/* this.props.color is null */}
 <Button color={null} />

 {/* this.props.color is "red" */}
 <Button color="red" />
</>
```

<Note>

Defining `defaultProps` in class components is similar to using [default values](/learn/passing-props-to-a-component#specifying-a-default-value-for-a-prop) in function components.

</Note>

---

### `static getDerivedStateFromError(error)` {/*static-getderivedstatefromerror*/}

If you define `static getDerivedStateFromError`, React will call it when a child component (including distant children) throws an error during rendering. This lets you display an error message instead of clearing the UI.

Typically, it is used together with [`componentDidCatch`](#componentdidcatch) which lets you send the error report to some analytics service. A component with these methods is called an *Error Boundary*.

[See an example.](#catching-rendering-errors-with-an-error-boundary)

#### Parameters {/*static-getderivedstatefromerror-parameters*/}

* `error`: The error that was thrown. In practice, it will usually be an instance of [`Error`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Error) but this is not guaranteed because JavaScript allows to [`throw`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/throw) any value, including strings or even `null`.

#### Returns {/*static-getderivedstatefromerror-returns*/}

`static getDerivedStateFromError` should return the state telling the component to display the error message.

#### Caveats {/*static-getderivedstatefromerror-caveats*/}

* `static getDerivedStateFromError` should be a pure function. If you want to perform a side effect (for example, to call an analytics service), you need to also implement [`componentDidCatch`.](#componentdidcatch)

<Note>

There is no direct equivalent for `static getDerivedStateFromError` in function components yet. If you'd like to avoid creating class components, write a single `ErrorBoundary` component like above and use it throughout your app. Alternatively, use the [`react-error-boundary`](https://github.com/bvaughn/react-error-boundary) package which does that.

</Note>

---

### `static getDerivedStateFromProps(props, state)` {/*static-getderivedstatefromprops*/}

If you define `static getDerivedStateFromProps`, React will call it right before calling [`render`,](#render) both on the initial mount and on subsequent updates. It should return an object to update the state, or `null` to update nothing.

This method exists for [rare use cases](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#when-to-use-derived-state) where the state depends on changes in props over time. For example, this `Form` component resets the `email` state when the `userID` prop changes:

```js {7-18}
class Form extends Component {
 state = {
 email: this.props.defaultEmail,
 prevUserID: this.props.userID
 };

 static getDerivedStateFromProps(props, state) {
 // Any time the current user changes,
 // Reset any parts of state that are tied to that user.
 // In this simple example, that's just the email.
 if (props.userID !== state.prevUserID) {
 return {
 prevUserID: props.userID,
 email: props.defaultEmail
 };
 }
 return null;
 }

 // ...
}
```

Note that this pattern requires you to keep a previous value of the prop (like `userID`) in state (like `prevUserID`).

<Pitfall>

Deriving state leads to verbose code and makes your components difficult to think about. [Make sure you're familiar with simpler alternatives:](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html)

- If you need to **perform a side effect** (for example, data fetching or an animation) in response to a change in props, use [`componentDidUpdate`](#componentdidupdate) method instead.
- If you want to **re-compute some data only when a prop changes,** [use a memoization helper instead.](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#what-about-memoization)
- If you want to **"reset" some state when a prop changes,** consider either making a component [fully controlled](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#recommendation-fully-controlled-component) or [fully uncontrolled with a key](https://legacy.reactjs.org/blog/2018/06/07/you-probably-dont-need-derived-state.html#recommendation-fully-uncontrolled-component-with-a-key) instead.

</Pitfall>

#### Parameters {/*static-getderivedstatefromprops-parameters*/}

- `props`: The next props that the component is about to render with.
- `state`: The next state that the component is about to render with.

#### Returns {/*static-getderivedstatefromprops-returns*/}

`static getDerivedStateFromProps` return an object to update the state, or `null` to update nothing.

#### Caveats {/*static-getderivedstatefromprops-caveats*/}

- This method is fired on *every* render, regardless of the cause. This is different from [`UNSAFE_componentWillReceiveProps`](#unsafe_componentwillreceiveprops), which only fires when the parent causes a re-render and not as a result of a local `setState`.

- This method doesn't have access to the component instance. If you'd like, you can reuse some code between `static getDerivedStateFromProps` and the other class methods by extracting pure functions of the component props and state outside the class definition.

<Note>

Implementing `static getDerivedStateFromProps` in a class component is equivalent to [calling the `set` function from `useState` during rendering](/reference/react/useState#storing-information-from-previous-renders) in a function component.

</Note>

---

## Usage {/*usage*/}

### Defining a class component {/*defining-a-class-component*/}

To define a React component as a class, extend the built-in `Component` class and define a [`render` method:](#render)

```js
import { Component } from 'react';

class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

React will call your [`render`](#render) method whenever it needs to figure out what to display on the screen. Usually, you will return some [JSX](/learn/writing-markup-with-jsx) from it. Your `render` method should be a [pure function:](https://en.wikipedia.org/wiki/Pure_function) it should only calculate the JSX.

Similarly to [function components,](/learn/your-first-component#defining-a-component) a class component can [receive information by props](/learn/your-first-component#defining-a-component) from its parent component. However, the syntax for reading props is different. For example, if the parent component renders `<Greeting name="Taylor" />`, then you can read the `name` prop from [`this.props`](#props), like `this.props.name`:

<Sandpack>

```js
import { Component } from 'react';

class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}

export default function App() {
 return (
 <>
 <Greeting name="Sara" />
 <Greeting name="Cahal" />
 <Greeting name="Edite" />
 </>
 );
}
```

</Sandpack>

Note that Hooks (functions starting with `use`, like [`useState`](/reference/react/useState)) are not supported inside class components.

<Pitfall>

We recommend defining components as functions instead of classes. [See how to migrate.](#migrating-a-simple-component-from-a-class-to-a-function)

</Pitfall>

---

### Adding state to a class component {/*adding-state-to-a-class-component*/}

To add [state](/learn/state-a-components-memory) to a class, assign an object to a property called [`state`](#state). To update state, call [`this.setState`](#setstate).

<Sandpack>

```js
import { Component } from 'react';

export default class Counter extends Component {
 state = {
 name: 'Taylor',
 age: 42,
 };

 handleNameChange = (e) => {
 this.setState({
 name: e.target.value
 });
 }

 handleAgeChange = () => {
 this.setState({
 age: this.state.age + 1
 });
 };

 render() {
 return (
 <>
 <input
 value={this.state.name}
 onChange={this.handleNameChange}
 />
 <button onClick={this.handleAgeChange}>
 Increment age
 </button>
 <p>Hello, {this.state.name}. You are {this.state.age}.</p>
 </>
 );
 }
}
```

```css
button { display: block; margin-top: 10px; }
```

</Sandpack>

<Pitfall>

We recommend defining components as functions instead of classes. [See how to migrate.](#migrating-a-component-with-state-from-a-class-to-a-function)

</Pitfall>

---

### Adding lifecycle methods to a class component {/*adding-lifecycle-methods-to-a-class-component*/}

There are a few special methods you can define on your class.

If you define the [`componentDidMount`](#componentdidmount) method, React will call it when your component is added *(mounted)* to the screen. React will call [`componentDidUpdate`](#componentdidupdate) after your component re-renders due to changed props or state. React will call [`componentWillUnmount`](#componentwillunmount) after your component has been removed *(unmounted)* from the screen.

If you implement `componentDidMount`, you usually need to implement all three lifecycles to avoid bugs. For example, if `componentDidMount` reads some state or props, you also have to implement `componentDidUpdate` to handle their changes, and `componentWillUnmount` to clean up whatever `componentDidMount` was doing.

For example, this `ChatRoom` component keeps a chat connection synchronized with props and state:

<Sandpack>

```js src/App.js
import { useState } from 'react';
import ChatRoom from './ChatRoom.js';

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/ChatRoom.js active
import { Component } from 'react';
import { createConnection } from './chat.js';

export default class ChatRoom extends Component {
 state = {
 serverUrl: 'https://localhost:1234'
 };

 componentDidMount() {
 this.setupConnection();
 }

 componentDidUpdate(prevProps, prevState) {
 if (
 this.props.roomId !== prevProps.roomId ||
 this.state.serverUrl !== prevState.serverUrl
 ) {
 this.destroyConnection();
 this.setupConnection();
 }
 }

 componentWillUnmount() {
 this.destroyConnection();
 }

 setupConnection() {
 this.connection = createConnection(
 this.state.serverUrl,
 this.props.roomId
 );
 this.connection.connect();
 }

 destroyConnection() {
 this.connection.disconnect();
 this.connection = null;
 }

 render() {
 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={this.state.serverUrl}
 onChange={e => {
 this.setState({
 serverUrl: e.target.value
 });
 }}
 />
 </label>
 <h1>Welcome to the {this.props.roomId} room!</h1>
 </>
 );
 }
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

Note that in development when [Strict Mode](/reference/react/StrictMode) is on, React will call `componentDidMount`, immediately call `componentWillUnmount`, and then call `componentDidMount` again. This helps you notice if you forgot to implement `componentWillUnmount` or if its logic doesn't fully "mirror" what `componentDidMount` does.

<Pitfall>

We recommend defining components as functions instead of classes. [See how to migrate.](#migrating-a-component-with-lifecycle-methods-from-a-class-to-a-function)

</Pitfall>

---

### Catching rendering errors with an Error Boundary {/*catching-rendering-errors-with-an-error-boundary*/}

By default, if your application throws an error during rendering, React will remove its UI from the screen. To prevent this, you can wrap a part of your UI into an *Error Boundary*. An Error Boundary is a special component that lets you display some fallback UI instead of the part that crashed--for example, an error message.

<Note>
Error boundaries do not catch errors for:

- Event handlers [(learn more)](/learn/responding-to-events)
- [Server side rendering](/reference/react-dom/server)
- Errors thrown in the error boundary itself (rather than its children)
- Asynchronous code (e.g. `setTimeout` or `requestAnimationFrame` callbacks); an exception is the usage of the [`startTransition`](/reference/react/useTransition#starttransition) function returned by the [`useTransition`](/reference/react/useTransition) Hook. Errors thrown inside the transition function are caught by error boundaries [(learn more)](/reference/react/useTransition#displaying-an-error-to-users-with-error-boundary)

</Note>

To implement an Error Boundary component, you need to provide [`static getDerivedStateFromError`](#static-getderivedstatefromerror) which lets you update state in response to an error and display an error message to the user. You can also optionally implement [`componentDidCatch`](#componentdidcatch) to add some extra logic, for example, to log the error to an analytics service.

With [`captureOwnerStack`](/reference/react/captureOwnerStack) you can include the Owner Stack during development.

```js {9-12,14-27}
import * as React from 'react';

class ErrorBoundary extends React.Component {
 constructor(props) {
 super(props);
 this.state = { hasError: false };
 }

 static getDerivedStateFromError(error) {
 // Update state so the next render will show the fallback UI.
 return { hasError: true };
 }

 componentDidCatch(error, info) {
 logErrorToMyService(
 error,
 // Example "componentStack":
 // in ComponentThatThrows (created by App)
 // in ErrorBoundary (created by App)
 // in div (created by App)
 // in App
 info.componentStack,
 // Warning: `captureOwnerStack` is not available in production.
 React.captureOwnerStack(),
 );
 }

 render() {
 if (this.state.hasError) {
 // You can render any custom fallback UI
 return this.props.fallback;
 }

 return this.props.children;
 }
}
```

Then you can wrap a part of your component tree with it:

```js {1,3}
<ErrorBoundary fallback={<p>Something went wrong</p>}>
 <Profile />
</ErrorBoundary>
```

If `Profile` or its child component throws an error, `ErrorBoundary` will "catch" that error, display a fallback UI with the error message you've provided, and send a production error report to your error reporting service.

You don't need to wrap every component into a separate Error Boundary. When you think about the [granularity of Error Boundaries,](https://www.brandondail.com/posts/fault-tolerance-react) consider where it makes sense to display an error message. For example, in a messaging app, it makes sense to place an Error Boundary around the list of conversations. It also makes sense to place one around every individual message. However, it wouldn't make sense to place a boundary around every avatar.

<Note>

There is currently no way to write an Error Boundary as a function component. However, you don't have to write the Error Boundary class yourself. For example, you can use [`react-error-boundary`](https://github.com/bvaughn/react-error-boundary) instead.

</Note>

---

## Alternatives {/*alternatives*/}

### Migrating a simple component from a class to a function {/*migrating-a-simple-component-from-a-class-to-a-function*/}

Typically, you will [define components as functions](/learn/your-first-component#defining-a-component) instead.

For example, suppose you're converting this `Greeting` class component to a function:

<Sandpack>

```js
import { Component } from 'react';

class Greeting extends Component {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}

export default function App() {
 return (
 <>
 <Greeting name="Sara" />
 <Greeting name="Cahal" />
 <Greeting name="Edite" />
 </>
 );
}
```

</Sandpack>

Define a function called `Greeting`. This is where you will move the body of your `render` function.

```js
function Greeting() {
 // ... move the code from the render method here ...
}
```

Instead of `this.props.name`, define the `name` prop [using the destructuring syntax](/learn/passing-props-to-a-component) and read it directly:

```js
function Greeting({ name }) {
 return <h1>Hello, {name}!</h1>;
}
```

Here is a complete example:

<Sandpack>

```js
function Greeting({ name }) {
 return <h1>Hello, {name}!</h1>;
}

export default function App() {
 return (
 <>
 <Greeting name="Sara" />
 <Greeting name="Cahal" />
 <Greeting name="Edite" />
 </>
 );
}
```

</Sandpack>

---

### Migrating a component with state from a class to a function {/*migrating-a-component-with-state-from-a-class-to-a-function*/}

Suppose you're converting this `Counter` class component to a function:

<Sandpack>

```js
import { Component } from 'react';

export default class Counter extends Component {
 state = {
 name: 'Taylor',
 age: 42,
 };

 handleNameChange = (e) => {
 this.setState({
 name: e.target.value
 });
 }

 handleAgeChange = (e) => {
 this.setState({
 age: this.state.age + 1
 });
 };

 render() {
 return (
 <>
 <input
 value={this.state.name}
 onChange={this.handleNameChange}
 />
 <button onClick={this.handleAgeChange}>
 Increment age
 </button>
 <p>Hello, {this.state.name}. You are {this.state.age}.</p>
 </>
 );
 }
}
```

```css
button { display: block; margin-top: 10px; }
```

</Sandpack>

Start by declaring a function with the necessary [state variables:](/reference/react/useState#adding-state-to-a-component)

```js {4-5}
import { useState } from 'react';

function Counter() {
 const [name, setName] = useState('Taylor');
 const [age, setAge] = useState(42);
 // ...
```

Next, convert the event handlers:

```js {5-7,9-11}
function Counter() {
 const [name, setName] = useState('Taylor');
 const [age, setAge] = useState(42);

 function handleNameChange(e) {
 setName(e.target.value);
 }

 function handleAgeChange() {
 setAge(age + 1);
 }
 // ...
```

Finally, replace all references starting with `this` with the variables and functions you defined in your component. For example, replace `this.state.age` with `age`, and replace `this.handleNameChange` with `handleNameChange`.

Here is a fully converted component:

<Sandpack>

```js
import { useState } from 'react';

export default function Counter() {
 const [name, setName] = useState('Taylor');
 const [age, setAge] = useState(42);

 function handleNameChange(e) {
 setName(e.target.value);
 }

 function handleAgeChange() {
 setAge(age + 1);
 }

 return (
 <>
 <input
 value={name}
 onChange={handleNameChange}
 />
 <button onClick={handleAgeChange}>
 Increment age
 </button>
 <p>Hello, {name}. You are {age}.</p>
 </>
 )
}
```

```css
button { display: block; margin-top: 10px; }
```

</Sandpack>

---

### Migrating a component with lifecycle methods from a class to a function {/*migrating-a-component-with-lifecycle-methods-from-a-class-to-a-function*/}

Suppose you're converting this `ChatRoom` class component with lifecycle methods to a function:

<Sandpack>

```js src/App.js
import { useState } from 'react';
import ChatRoom from './ChatRoom.js';

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/ChatRoom.js active
import { Component } from 'react';
import { createConnection } from './chat.js';

export default class ChatRoom extends Component {
 state = {
 serverUrl: 'https://localhost:1234'
 };

 componentDidMount() {
 this.setupConnection();
 }

 componentDidUpdate(prevProps, prevState) {
 if (
 this.props.roomId !== prevProps.roomId ||
 this.state.serverUrl !== prevState.serverUrl
 ) {
 this.destroyConnection();
 this.setupConnection();
 }
 }

 componentWillUnmount() {
 this.destroyConnection();
 }

 setupConnection() {
 this.connection = createConnection(
 this.state.serverUrl,
 this.props.roomId
 );
 this.connection.connect();
 }

 destroyConnection() {
 this.connection.disconnect();
 this.connection = null;
 }

 render() {
 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={this.state.serverUrl}
 onChange={e => {
 this.setState({
 serverUrl: e.target.value
 });
 }}
 />
 </label>
 <h1>Welcome to the {this.props.roomId} room!</h1>
 </>
 );
 }
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

First, verify that your [`componentWillUnmount`](#componentwillunmount) does the opposite of [`componentDidMount`.](#componentdidmount) In the above example, that's true: it disconnects the connection that `componentDidMount` sets up. If such logic is missing, add it first.

Next, verify that your [`componentDidUpdate`](#componentdidupdate) method handles changes to any props and state you're using in `componentDidMount`. In the above example, `componentDidMount` calls `setupConnection` which reads `this.state.serverUrl` and `this.props.roomId`. This is why `componentDidUpdate` checks whether `this.state.serverUrl` and `this.props.roomId` have changed, and resets the connection if they did. If your `componentDidUpdate` logic is missing or doesn't handle changes to all relevant props and state, fix that first.

In the above example, the logic inside the lifecycle methods connects the component to a system outside of React (a chat server). To connect a component to an external system, [describe this logic as a single Effect:](/reference/react/useEffect#connecting-to-an-external-system)

```js {6-12}
import { useState, useEffect } from 'react';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [serverUrl, roomId]);

 // ...
}
```

This [`useEffect`](/reference/react/useEffect) call is equivalent to the logic in the lifecycle methods above. If your lifecycle methods do multiple unrelated things, [split them into multiple independent Effects.](/learn/removing-effect-dependencies#is-your-effect-doing-several-unrelated-things) Here is a complete example you can play with:

<Sandpack>

```js src/App.js
import { useState } from 'react';
import ChatRoom from './ChatRoom.js';

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/ChatRoom.js active
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

export default function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [roomId, serverUrl]);

 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={serverUrl}
 onChange={e => setServerUrl(e.target.value)}
 />
 </label>
 <h1>Welcome to the {roomId} room!</h1>
 </>
 );
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

<Note>

If your component does not synchronize with any external systems, [you might not need an Effect.](/learn/you-might-not-need-an-effect)

</Note>

---

### Migrating a component with context from a class to a function {/*migrating-a-component-with-context-from-a-class-to-a-function*/}

In this example, the `Panel` and `Button` class components read [context](/learn/passing-data-deeply-with-context) from [`this.context`:](#context)

<Sandpack>

```js
import { createContext, Component } from 'react';

const ThemeContext = createContext(null);

class Panel extends Component {
 static contextType = ThemeContext;

 render() {
 const theme = this.context;
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{this.props.title}</h1>
 {this.props.children}
 </section>
 );
 }
}

class Button extends Component {
 static contextType = ThemeContext;

 render() {
 const theme = this.context;
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {this.props.children}
 </button>
 );
 }
}

function Form() {
 return (
 <Panel title="Welcome">
 <Button>Sign up</Button>
 <Button>Log in</Button>
 </Panel>
 );
}

export default function MyApp() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 )
}
```

```css
.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

When you convert them to function components, replace `this.context` with [`useContext`](/reference/react/useContext) calls:

<Sandpack>

```js
import { createContext, useContext } from 'react';

const ThemeContext = createContext(null);

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {children}
 </button>
 );
}

function Form() {
 return (
 <Panel title="Welcome">
 <Button>Sign up</Button>
 <Button>Log in</Button>
 </Panel>
 );
}

export default function MyApp() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 )
}
```

```css
.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

---
title: <Fragment> (<>...</>)
---

<Intro>

`<Fragment>`, often used via `<>...</>` syntax, lets you group elements without a wrapper node.

<Canary>Fragments can also accept refs, which enable interacting with underlying DOM nodes without adding wrapper elements.</Canary>

```js
<>
 <OneChild />
 <AnotherChild />
</>
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `<Fragment>` {/*fragment*/}

Wrap elements in `<Fragment>` to group them together in situations where you need a single element. Grouping elements in `Fragment` has no effect on the resulting DOM; it is the same as if the elements were not grouped. The empty JSX tag `<></>` is shorthand for `<Fragment></Fragment>` in most cases.

#### Props {/*props*/}

- **optional** `key`: Fragments declared with the explicit `<Fragment>` syntax may have [keys.](/learn/rendering-lists#keeping-list-items-in-order-with-key)
- <CanaryBadge /> **optional** `ref`: A ref object (e.g. from [`useRef`](/reference/react/useRef)) or [callback function](/reference/react-dom/components/common#ref-callback). React provides a `FragmentInstance` as the ref value that implements methods for interacting with the DOM nodes wrapped by the Fragment.

#### Caveats {/*caveats*/}

* If you want to pass `key` to a Fragment, you can't use the `<>...</>` syntax. You have to explicitly import `Fragment` from `'react'` and render `<Fragment key={yourKey}>...</Fragment>`.

* React does not [reset state](/learn/preserving-and-resetting-state) when you go from rendering `<><Child /></>` to `[<Child />]` or back, or when you go from rendering `<><Child /></>` to `<Child />` and back. This only works a single level deep: for example, going from `<><><Child /></></>` to `<Child />` resets the state. See the precise semantics [here.](https://gist.github.com/clemmy/b3ef00f9507909429d8aa0d3ee4f986b)

* <CanaryBadge /> If you want to pass `ref` to a Fragment, you can't use the `<>...</>` syntax. You have to explicitly import `Fragment` from `'react'` and render `<Fragment ref={yourRef}>...</Fragment>`.

---

### <CanaryBadge /> `FragmentInstance` {/*fragmentinstance*/}

When you pass a `ref` to a Fragment, React provides a `FragmentInstance` object. It implements methods for interacting with the first-level DOM children wrapped by the Fragment.

* [`addEventListener`](#addeventlistener) and [`removeEventListener`](#removeeventlistener) manage event listeners across all first-level DOM children.
* [`dispatchEvent`](#dispatchevent) dispatches an event on the Fragment, which can bubble to the DOM parent.
* [`focus`](#focus), [`focusLast`](#focuslast), and [`blur`](#blur) manage focus across all nested children depth-first.
* [`observeUsing`](#observeusing) and [`unobserveUsing`](#unobserveusing) attach and detach `IntersectionObserver` or `ResizeObserver` instances.
* [`getClientRects`](#getclientrects) returns bounding rectangles of all first-level DOM children.
* [`getRootNode`](#getrootnode) returns the root node of the Fragment's parent.
* [`compareDocumentPosition`](#comparedocumentposition) compares the Fragment's position with another node.
* [`scrollIntoView`](#scrollintoview) scrolls the Fragment's children into view.

---

#### `addEventListener(type, listener, options?)` {/*addeventlistener*/}

Adds an event listener to all first-level DOM children of the Fragment.

```js
fragmentRef.current.addEventListener('click', handleClick);
```

##### Parameters {/*addeventlistener-parameters*/}

* `type`: A string representing the event type to listen for (e.g. `'click'`, `'focus'`).
* `listener`: The event handler function.
* **optional** `options`: An options object or boolean for capture, matching the [DOM `addEventListener` API.](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener)

##### Returns {/*addeventlistener-returns*/}

`addEventListener` does not return anything (`undefined`).

---

#### `removeEventListener(type, listener, options?)` {/*removeeventlistener*/}

Removes an event listener from all first-level DOM children of the Fragment.

```js
fragmentRef.current.removeEventListener('click', handleClick);
```

##### Parameters {/*removeeventlistener-parameters*/}

* `type`: The event type string.
* `listener`: The event handler function to remove.
* **optional** `options`: An options object or boolean, matching the [DOM `removeEventListener` API.](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener)

##### Returns {/*removeeventlistener-returns*/}

`removeEventListener` does not return anything (`undefined`).

---

#### `dispatchEvent(event)` {/*dispatchevent*/}

Dispatches an event on the Fragment. Added event listeners are called, and the event can bubble to the Fragment's DOM parent.

```js
fragmentRef.current.dispatchEvent(new Event('custom', { bubbles: true }));
```

##### Parameters {/*dispatchevent-parameters*/}

* `event`: An [`Event`](https://developer.mozilla.org/en-US/docs/Web/API/Event) object to dispatch. If `bubbles` is `true`, the event bubbles to the Fragment's parent DOM node.

##### Returns {/*dispatchevent-returns*/}

`true` if the event was not cancelled, `false` if `preventDefault()` was called.

---

#### `focus(options?)` {/*focus*/}

Focuses the first focusable DOM node in the Fragment. Unlike calling `element.focus()` on a DOM element, this method searches *all* nested children depth-first until it finds a focusable element—not just the element itself or its direct children.

```js
fragmentRef.current.focus();
```

##### Parameters {/*focus-parameters*/}

* **optional** `options`: A [`FocusOptions`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/focus#options) object (e.g. `{ preventScroll: true }`).

##### Returns {/*focus-returns*/}

`focus` does not return anything (`undefined`).

---

#### `focusLast(options?)` {/*focuslast*/}

Focuses the last focusable DOM node in the Fragment. Searches nested children depth-first, then iterates in reverse.

```js
fragmentRef.current.focusLast();
```

##### Parameters {/*focuslast-parameters*/}

* **optional** `options`: A [`FocusOptions`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/focus#options) object.

##### Returns {/*focuslast-returns*/}

`focusLast` does not return anything (`undefined`).

---

#### `blur()` {/*blur*/}

Removes focus from the active element if it is within the Fragment. If `document.activeElement` is not within the Fragment, `blur` does nothing.

```js
fragmentRef.current.blur();
```

##### Returns {/*blur-returns*/}

`blur` does not return anything (`undefined`).

---

#### `observeUsing(observer)` {/*observeusing*/}

Starts observing all first-level DOM children of the Fragment with the provided observer.

```js
const observer = new IntersectionObserver(callback, options);
fragmentRef.current.observeUsing(observer);
```

##### Parameters {/*observeusing-parameters*/}

* `observer`: An [`IntersectionObserver`](https://developer.mozilla.org/en-US/docs/Web/API/IntersectionObserver) or [`ResizeObserver`](https://developer.mozilla.org/en-US/docs/Web/API/ResizeObserver) instance.

##### Returns {/*observeusing-returns*/}

`observeUsing` does not return anything (`undefined`).

---

#### `unobserveUsing(observer)` {/*unobserveusing*/}

Stops observing the Fragment's DOM children with the specified observer.

```js
fragmentRef.current.unobserveUsing(observer);
```

##### Parameters {/*unobserveusing-parameters*/}

* `observer`: The same `IntersectionObserver` or `ResizeObserver` instance previously passed to [`observeUsing`](#observeusing).

##### Returns {/*unobserveusing-returns*/}

`unobserveUsing` does not return anything (`undefined`).

---

#### `getClientRects()` {/*getclientrects*/}

Returns a flat array of [`DOMRect`](https://developer.mozilla.org/en-US/docs/Web/API/DOMRect) objects representing the bounding rectangles of all first-level DOM children.

```js
const rects = fragmentRef.current.getClientRects();
```

##### Returns {/*getclientrects-returns*/}

An `Array<DOMRect>` containing the bounding rectangles of all children.

---

#### `getRootNode(options?)` {/*getrootnode*/}

Returns the root node containing the Fragment's parent DOM node, matching the behavior of [`Node.getRootNode()`](https://developer.mozilla.org/en-US/docs/Web/API/Node/getRootNode).

```js
const root = fragmentRef.current.getRootNode();
```

##### Parameters {/*getrootnode-parameters*/}

* **optional** `options`: An object with a `composed` boolean property, matching the [DOM `getRootNode` API.](https://developer.mozilla.org/en-US/docs/Web/API/Node/getRootNode#options)

##### Returns {/*getrootnode-returns*/}

A `Document`, `ShadowRoot`, or the `FragmentInstance` itself if there is no parent DOM node.

---

#### `compareDocumentPosition(otherNode)` {/*comparedocumentposition*/}

Compares the document position of the Fragment with another node, returning a bitmask matching the behavior of [`Node.compareDocumentPosition()`](https://developer.mozilla.org/en-US/docs/Web/API/Node/compareDocumentPosition).

```js
const position = fragmentRef.current.compareDocumentPosition(otherElement);
```

##### Parameters {/*comparedocumentposition-parameters*/}

* `otherNode`: The DOM node to compare against.

##### Returns {/*comparedocumentposition-returns*/}

A bitmask of [position flags](https://developer.mozilla.org/en-US/docs/Web/API/Node/compareDocumentPosition#return_value). Empty Fragments and Fragments with children rendered through a [portal](/reference/react-dom/createPortal) include `Node.DOCUMENT_POSITION_IMPLEMENTATION_SPECIFIC` in the result.

---

#### `scrollIntoView(alignToTop?)` {/*scrollintoview*/}

Scrolls the Fragment's children into view. When `alignToTop` is `true` or omitted, scrolls to align the first child with the top of the scrollable ancestor. When `alignToTop` is `false`, scrolls to align the last child with the bottom.

```js
fragmentRef.current.scrollIntoView();
```

##### Parameters {/*scrollintoview-parameters*/}

* **optional** `alignToTop`: A boolean. If `true` (the default), scrolls the first child to the top of the scrollable area. If `false`, scrolls the last child to the bottom. Unlike [`Element.scrollIntoView()`](https://developer.mozilla.org/en-US/docs/Web/API/Element/scrollIntoView), this method does not accept a `ScrollIntoViewOptions` object.

##### Returns {/*scrollintoview-returns*/}

`scrollIntoView` does not return anything (`undefined`).

##### Caveats {/*scrollintoview-caveats*/}

* `scrollIntoView` does not accept an options object. Passing one throws an error. Use the `alignToTop` boolean instead.
* When the Fragment has no children, `scrollIntoView` scrolls the nearest sibling or parent into view as a fallback.

---

#### `FragmentInstance` Caveats {/*fragmentinstance-caveats*/}

* Methods that target children (such as `addEventListener`, `observeUsing`, and `getClientRects`) operate on *first-level host (DOM) children* of the Fragment. They do not directly target children nested inside another DOM element.
* `focus` and `focusLast` search nested children depth-first for focusable elements, unlike event and observer methods which only target first-level host children.
* `observeUsing` does not work on text nodes. React logs a warning in development if the Fragment contains only text children.
* React does not apply event listeners added via `addEventListener` to hidden [`<Activity>`](/reference/react/Activity) trees. When an `Activity` boundary switches from hidden to visible, listeners are applied automatically.
* Each first-level DOM child of a Fragment with a `ref` gets a `reactFragments` property—a `Set<FragmentInstance>` containing all Fragment instances that own the element. This enables [caching a shared observer](#caching-global-intersection-observer) across multiple Fragments.

---

## Usage {/*usage*/}

### Returning multiple elements {/*returning-multiple-elements*/}

Use `Fragment`, or the equivalent `<>...</>` syntax, to group multiple elements together. You can use it to put multiple elements in any place where a single element can go. For example, a component can only return one element, but by using a Fragment you can group multiple elements together and then return them as a group:

```js {3,6}
function Post() {
 return (
 <>
 <PostTitle />
 <PostBody />
 </>
 );
}
```

Fragments are useful because grouping elements with a Fragment has no effect on layout or styles, unlike if you wrapped the elements in another container like a DOM element. If you inspect this example with the browser tools, you'll see that all `<h1>` and `<article>` DOM nodes appear as siblings without wrappers around them:

<Sandpack>

```js
export default function Blog() {
 return (
 <>
 <Post title="An update" body="It's been a while since I posted..." />
 <Post title="My new blog" body="I am starting a new blog!" />
 </>
 )
}

function Post({ title, body }) {
 return (
 <>
 <PostTitle title={title} />
 <PostBody body={body} />
 </>
 );
}

function PostTitle({ title }) {
 return <h1>{title}</h1>
}

function PostBody({ body }) {
 return (
 <article>
 <p>{body}</p>
 </article>
 );
}
```

</Sandpack>

<DeepDive>

#### How to write a Fragment without the special syntax? {/*how-to-write-a-fragment-without-the-special-syntax*/}

The example above is equivalent to importing `Fragment` from React:

```js {1,5,8}
import { Fragment } from 'react';

function Post() {
 return (
 <Fragment>
 <PostTitle />
 <PostBody />
 </Fragment>
 );
}
```

Usually you won't need this unless you need to [pass a `key` to your `Fragment`.](#rendering-a-list-of-fragments)

</DeepDive>

---

### Assigning multiple elements to a variable {/*assigning-multiple-elements-to-a-variable*/}

Like any other element, you can assign Fragment elements to variables, pass them as props, and so on:

```js
function CloseDialog() {
 const buttons = (
 <>
 <OKButton />
 <CancelButton />
 </>
 );
 return (
 <AlertDialog buttons={buttons}>
 Are you sure you want to leave this page?
 </AlertDialog>
 );
}
```

---

### Grouping elements with text {/*grouping-elements-with-text*/}

You can use `Fragment` to group text together with components:

```js
function DateRangePicker({ start, end }) {
 return (
 <>
 From
 <DatePicker date={start} />
 to
 <DatePicker date={end} />
 </>
 );
}
```

---

### Rendering a list of Fragments {/*rendering-a-list-of-fragments*/}

Here's a situation where you need to write `Fragment` explicitly instead of using the `<></>` syntax. When you [render multiple elements in a loop](/learn/rendering-lists), you need to assign a `key` to each element. If the elements within the loop are Fragments, you need to use the normal JSX element syntax in order to provide the `key` attribute:

```js {3,6}
function Blog() {
 return posts.map(post =>
 <Fragment key={post.id}>
 <PostTitle title={post.title} />
 <PostBody body={post.body} />
 </Fragment>
 );
}
```

You can inspect the DOM to verify that there are no wrapper elements around the Fragment children:

<Sandpack>

```js
import { Fragment } from 'react';

const posts = [
 { id: 1, title: 'An update', body: "It's been a while since I posted..." },
 { id: 2, title: 'My new blog', body: 'I am starting a new blog!' }
];

export default function Blog() {
 return posts.map(post =>
 <Fragment key={post.id}>
 <PostTitle title={post.title} />
 <PostBody body={post.body} />
 </Fragment>
 );
}

function PostTitle({ title }) {
 return <h1>{title}</h1>
}

function PostBody({ body }) {
 return (
 <article>
 <p>{body}</p>
 </article>
 );
}
```

</Sandpack>

---

### <CanaryBadge /> Adding event listeners without a wrapper element {/*adding-event-listeners-without-wrapper*/}

Fragment `ref`s let you add event listeners to a group of elements without adding a wrapper DOM node. Use a [ref callback](/reference/react-dom/components/common#ref-callback) to attach and clean up listeners:

<Sandpack>

```js
import { Fragment, useState, useRef, useEffect } from 'react';

function ClickableFragment({ children, onClick }) {
 const fragmentRef = useRef(null);
 useEffect(() => {
 const fragmentInstance = fragmentRef.current;
 if (fragmentInstance === null) {
 return;
 }
 fragmentInstance.addEventListener('click', onClick);
 return () => {
 fragmentInstance.removeEventListener(
 'click',
 onClick
 );
 };
 }, [onClick])
 return (
 <Fragment ref={fragmentRef}>
 {children}
 </Fragment>
 );
}

export default function App() {
 const [clicks, setClicks] = useState(0);

 return (
 <>
 <p>Total clicks: {clicks}</p>
 <ClickableFragment onClick={() => {
 setClicks(c => c + 1);
 }}>
 <button>Button A</button>
 <button>Button B</button>
 <button>Button C</button>
 </ClickableFragment>
 </>
 );
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

The `addEventListener` call applies the listener to every first-level DOM child of the Fragment. When children are dynamically added or removed, the `FragmentInstance` automatically adds or removes the listener.

<DeepDive>

#### Which children does a Fragment ref target? {/*which-children-does-a-fragment-ref-target*/}

A `FragmentInstance` targets the **first-level host (DOM) children** of the Fragment. Consider this tree:

```js
<Fragment ref={ref}>
 <div id="A" />
 <Wrapper>
 <div id="B">
 <div id="C" />
 </div>
 </Wrapper>
 <div id="D" />
</Fragment>
```

`Wrapper` is a React component, so the `FragmentInstance` looks through it to find DOM nodes. The targeted children are `A`, `B`, and `D`. `C` is not targeted because it is nested inside the DOM element `B`.

Methods like `addEventListener`, `observeUsing`, and `getClientRects` operate on these first-level DOM children. `focus` and `focusLast` are different—they search *all* nested children depth-first to find focusable elements.

</DeepDive>

---

### <CanaryBadge /> Managing focus across a group of elements {/*managing-focus-across-elements*/}

Fragment `ref`s provide `focus`, `focusLast`, and `blur` methods that operate across all DOM nodes within the Fragment:

<Sandpack>

```js
import { Fragment, useRef } from 'react';

function FormFields({ children }) {
 const fragmentRef = useRef(null);

 return (
 <>
 <div className="buttons">
 <button onClick={() => {
 fragmentRef.current.focus();
 }}>
 Focus first
 </button>
 <button onClick={() => {
 fragmentRef.current.focusLast();
 }}>
 Focus last
 </button>
 <button onClick={() => {
 fragmentRef.current.blur();
 }}>
 Blur
 </button>
 </div>
 <Fragment ref={fragmentRef}>
 {children}
 </Fragment>
 </>
 );
}

// Even though the inputs are deeply nested,
// focus() searches depth-first to find them.
export default function App() {
 return (
 <FormFields>
 <fieldset>
 <legend>Shipping</legend>
 <label>
 Street: <input name="street" />
 </label>
 <label>
 City: <input name="city" />
 </label>
 </fieldset>
 </FormFields>
 );
}
```

```css
.buttons {
 display: flex;
 gap: 8px;
 margin-bottom: 10px;
}

label {
 display: inline-block;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

Calling `focus()` focuses the `street` input—even though it is nested inside a `<fieldset>` and `<label>`. `focus()` searches depth-first through all nested children, not just direct children of the Fragment. `focusLast()` does the same in reverse, and `blur()` removes focus if the currently focused element is within the Fragment.

---

### <CanaryBadge /> Scrolling a group of elements into view {/*scrolling-group-into-view*/}

Use `scrollIntoView` to scroll a Fragment's children into view without a wrapper element. Pass `true` (or omit the argument) to scroll the first child to the top. Pass `false` to scroll the last child to the bottom:

<Sandpack>

```js
import { Fragment, useRef } from 'react';

function ScrollableSection({ children }) {
 const fragmentRef = useRef(null);

 return (
 <>
 <div className="buttons">
 <button onClick={() => {
 fragmentRef.current.scrollIntoView();
 }}>
 Scroll to top
 </button>
 <button onClick={() => {
 fragmentRef.current.scrollIntoView(false);
 }}>
 Scroll to bottom
 </button>
 </div>
 <div className="container">
 <Fragment ref={fragmentRef}>
 {children}
 </Fragment>
 </div>
 </>
 );
}

const items = [];
for (let i = 1; i <= 25; i++) {
 items.push('Item ' + i);
}

export default function App() {
 return (
 <ScrollableSection>
 <h3>Section Start</h3>
 {items.map((item) => (
 <p key={item}>{item}</p>
 ))}
 <h3>Section End</h3>
 </ScrollableSection>
 );
}
```

```css
.buttons {
 display: flex;
 gap: 8px;
 margin-bottom: 10px;
}

.container {
 height: 200px;
 overflow-y: auto;
 border: 2px solid #c4c4c4;
 border-radius: 4px;
 padding: 10px;
}

h3 {
 margin: 4px 0;
 /* Padding to handle offset of global sticky nav when scrolling for example */
 padding-top: 4em;
 color: #1a73e8;
}

p {
 margin: 4px 0;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### <CanaryBadge /> Observing visibility without a wrapper element {/*observing-visibility-without-wrapper*/}

Use `observeUsing` to attach an `IntersectionObserver` to all first-level DOM children of a Fragment. This lets you track visibility without requiring child components to expose `ref`s or adding a wrapper element:

<Sandpack>

```js
import {
 Fragment,
 useRef,
 useLayoutEffect,
 useState,
} from 'react';
import Card from './Card';

function VisibleGroup({ onVisibilityChange, children }) {
 const fragmentRef = useRef(null);

 useLayoutEffect(() => {
 const visibleElements = new Set();
 const observer = new IntersectionObserver(
 (entries) => {
 entries.forEach(e => {
 if (e.isIntersecting) {
 visibleElements.add(e.target);
 } else {
 visibleElements.delete(e.target);
 }
 });
 onVisibilityChange(visibleElements.size > 0);
 }
 );
 const fragmentInstance = fragmentRef.current;
 fragmentInstance.observeUsing(observer);
 return () => {
 fragmentInstance.unobserveUsing(observer);
 };
 }, [onVisibilityChange]);

 return (
 <Fragment ref={fragmentRef}>
 {children}
 </Fragment>
 );
}

export default function App() {
 const [isVisible, setIsVisible] = useState(true);

 return (
 <div className={isVisible ? 'page visible' : 'page'}>
 <div className="filler">Scroll down</div>
 <VisibleGroup onVisibilityChange={setIsVisible}>
 <Card title="First section" />
 <Card title="Second section" />
 </VisibleGroup>
 <div className="filler">Scroll up</div>
 </div>
 );
}
```

```css
.page {
 transition: background 0.3s;
}

.page.visible {
 background: #d4edda;
}

.filler {
 height: 500px;
 display: flex;
 align-items: center;
 justify-content: center;
 color: #aaa;
 font-size: 14px;
}

.card {
 padding: 16px;
 background: white;
 border: 1px solid #ddd;
 border-radius: 8px;
 margin: 8px 16px;
 box-shadow: 0 1px 3px rgba(0,0,0,0.08);
 font-weight: 600;
 font-size: 14px;
}
```

```js src/Card.js hidden
export default function Card({ title }) {
 return <div className="card">{title}</div>;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### <CanaryBadge /> Caching a global IntersectionObserver {/*caching-global-intersection-observer*/}

A common performance optimization for sites with many observers is to share a single IntersectionObserver per config and route its entries to the correct callbacks based on which element intersected. Fragment `ref`s support this same pattern through the `reactFragments` property.

Each first-level DOM child of a Fragment with a `ref` has a `reactFragments` property: a `Set` of `FragmentInstance` objects that contain that element. When the shared observer fires, you can use this property to look up which `FragmentInstance` owns the intersecting element and run the right callbacks.

<Sandpack>

```js src/App.js active
import { useState, useCallback } from 'react';
import ObservedGroup from './ObservedGroup';
import Card from './Card';

export default function App() {
 const [bgColor, setBgColor] = useState(null);

 const onGreen = useCallback((entry) => {
 if (entry.isIntersecting) {
 setBgColor('#d4edda');
 }
 }, []);

 const onBlue = useCallback((entry) => {
 if (entry.isIntersecting) {
 setBgColor('#cce5ff');
 }
 }, []);

 return (
 <div className="page" style={{
 background: bgColor || 'white',
 }}>
 <div className="filler">Scroll down</div>
 <ObservedGroup onIntersection={onGreen}>
 <Card title="Green section" className="green" />
 </ObservedGroup>
 <div className="filler" />
 <ObservedGroup onIntersection={onBlue}>
 <Card title="Blue section" className="blue" />
 </ObservedGroup>
 <div className="filler">Scroll up</div>
 </div>
 );
}
```

```js src/ObservedGroup.js
import {
 Fragment,
 useRef,
 useLayoutEffect,
} from 'react';

const callbackMap = new WeakMap();
const observerCache = new Map();

function getOptionsKey(options) {
 const root = options?.root ?? null;
 const rootMargin = options?.rootMargin ?? '0px';
 const threshold = options?.threshold ?? 0;
 return `${rootMargin}|${threshold}`;
}

function getSharedObserver(
 fragmentInstance,
 onIntersection,
 options,
) {
 // Register this callback for the
 // fragment instance.
 const existing =
 callbackMap.get(fragmentInstance);
 callbackMap.set(
 fragmentInstance,
 existing
 ? [...existing, onIntersection]
 : [onIntersection],
 );

 const key = getOptionsKey(options);
 if (observerCache.has(key)) {
 return observerCache.get(key);
 }

 const observer = new IntersectionObserver(
 (entries) => {
 for (const entry of entries) {
 // Look up which FragmentInstances own
 // this element.
 const fragmentInstances =
 entry.target.reactFragments;
 if (fragmentInstances) {
 for (const inst of fragmentInstances) {
 const callbacks =
 callbackMap.get(inst) || [];
 callbacks.forEach(cb => cb(entry));
 }
 }
 }
 },
 options,
 );

 observerCache.set(key, observer);
 return observer;
}

export default function ObservedGroup({
 onIntersection,
 options,
 children,
}) {
 const fragmentRef = useRef(null);

 useLayoutEffect(() => {
 const fragmentInstance = fragmentRef.current;
 const observer = getSharedObserver(
 fragmentInstance,
 onIntersection,
 options,
 );
 fragmentInstance.observeUsing(observer);
 return () => {
 fragmentInstance.unobserveUsing(observer);
 callbackMap.delete(fragmentInstance);
 };
 }, [onIntersection, options]);

 return (
 <Fragment ref={fragmentRef}>
 {children}
 </Fragment>
 );
}
```

```css
.page {
 transition: background 0.3s;
}

.filler {
 height: 500px;
 display: flex;
 align-items: center;
 justify-content: center;
 color: #aaa;
 font-size: 14px;
}

.card {
 padding: 16px;
 background: white;
 border: 1px solid #ddd;
 border-radius: 8px;
 margin: 0 16px;
 box-shadow: 0 1px 3px rgba(0,0,0,0.08);
 font-weight: 600;
 font-size: 14px;
}

.card.green {
 border-left: 3px solid #28a745;
}

.card.blue {
 border-left: 3px solid #007bff;
}
```

```js src/Card.js hidden
export default function Card({ title, className }) {
 return <div className={'card' + (className ? ' ' + className : '')}>{title}</div>;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

Multiple `ObservedGroup` components with the same options reuse a single `IntersectionObserver`. When either section scrolls into view, the shared observer fires and uses `reactFragments` to route the entry to the correct callback.

---
title: <Profiler>
---

<Intro>

`<Profiler>` lets you measure rendering performance of a React tree programmatically.

```js
<Profiler id="App" onRender={onRender}>
 <App />
</Profiler>
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `<Profiler>` {/*profiler*/}

Wrap a component tree in a `<Profiler>` to measure its rendering performance.

```js
<Profiler id="App" onRender={onRender}>
 <App />
</Profiler>
```

#### Props {/*props*/}

* `id`: A string identifying the part of the UI you are measuring.
* `onRender`: An [`onRender` callback](#onrender-callback) that React calls every time components within the profiled tree update. It receives information about what was rendered and how much time it took.

#### Caveats {/*caveats*/}

* Profiling adds some additional overhead, so **it is disabled in the production build by default.** To opt into production profiling, you need to enable a [special production build with profiling enabled.](/reference/dev-tools/react-performance-tracks#using-profiling-builds)

---

### `onRender` callback {/*onrender-callback*/}

React will call your `onRender` callback with information about what was rendered.

```js
function onRender(id, phase, actualDuration, baseDuration, startTime, commitTime) {
 // Aggregate or log render timings...
}
```

#### Parameters {/*onrender-parameters*/}

* `id`: The string `id` prop of the `<Profiler>` tree that has just committed. This lets you identify which part of the tree was committed if you are using multiple profilers.
* `phase`: `"mount"`, `"update"` or `"nested-update"`. This lets you know whether the tree has just been mounted for the first time or re-rendered due to a change in props, state, or Hooks.
* `actualDuration`: The number of milliseconds spent rendering the `<Profiler>` and its descendants for the current update. This indicates how well the subtree makes use of memoization (e.g. [`memo`](/reference/react/memo) and [`useMemo`](/reference/react/useMemo)). Ideally this value should decrease significantly after the initial mount as many of the descendants will only need to re-render if their specific props change.
* `baseDuration`: The number of milliseconds estimating how much time it would take to re-render the entire `<Profiler>` subtree without any optimizations. It is calculated by summing up the most recent render durations of each component in the tree. This value estimates a worst-case cost of rendering (e.g. the initial mount or a tree with no memoization). Compare `actualDuration` against it to see if memoization is working.
* `startTime`: A numeric timestamp for when React began rendering the current update.
* `commitTime`: A numeric timestamp for when React committed the current update. This value is shared between all profilers in a commit, enabling them to be grouped if desirable.

---

## Usage {/*usage*/}

### Measuring rendering performance programmatically {/*measuring-rendering-performance-programmatically*/}

Wrap the `<Profiler>` component around a React tree to measure its rendering performance.

```js {2,4}
<App>
 <Profiler id="Sidebar" onRender={onRender}>
 <Sidebar />
 </Profiler>
 <PageContent />
</App>
```

It requires two props: an `id` (string) and an `onRender` callback (function) which React calls any time a component within the tree "commits" an update.

<Pitfall>

Profiling adds some additional overhead, so **it is disabled in the production build by default.** To opt into production profiling, you need to enable a [special production build with profiling enabled.](/reference/dev-tools/react-performance-tracks#using-profiling-builds)

</Pitfall>

<Note>

`<Profiler>` lets you gather measurements programmatically. If you're looking for an interactive profiler, try the Profiler tab in [React Developer Tools](/learn/react-developer-tools). It exposes similar functionality as a browser extension.

Components wrapped in `<Profiler>` will also be marked in the [Component tracks](/reference/dev-tools/react-performance-tracks#components) of React Performance tracks even in profiling builds.
In development builds, all components are marked in the Components track regardless of whether they're wrapped in `<Profiler>`.

</Note>

---

### Measuring different parts of the application {/*measuring-different-parts-of-the-application*/}

You can use multiple `<Profiler>` components to measure different parts of your application:

```js {5,7}
<App>
 <Profiler id="Sidebar" onRender={onRender}>
 <Sidebar />
 </Profiler>
 <Profiler id="Content" onRender={onRender}>
 <Content />
 </Profiler>
</App>
```

You can also nest `<Profiler>` components:

```js {5,7,9,12}
<App>
 <Profiler id="Sidebar" onRender={onRender}>
 <Sidebar />
 </Profiler>
 <Profiler id="Content" onRender={onRender}>
 <Content>
 <Profiler id="Editor" onRender={onRender}>
 <Editor />
 </Profiler>
 <Preview />
 </Content>
 </Profiler>
</App>
```

Although `<Profiler>` is a lightweight component, it should be used only when necessary. Each use adds some CPU and memory overhead to an application.

---

---
title: PureComponent
---

<Pitfall>

We recommend defining components as functions instead of classes. [See how to migrate.](#alternatives)

</Pitfall>

<Intro>

`PureComponent` is similar to [`Component`](/reference/react/Component) but it skips re-renders for same props and state. Class components are still supported by React, but we don't recommend using them in new code.

```js
class Greeting extends PureComponent {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `PureComponent` {/*purecomponent*/}

To skip re-rendering a class component for same props and state, extend `PureComponent` instead of [`Component`:](/reference/react/Component)

```js
import { PureComponent } from 'react';

class Greeting extends PureComponent {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

`PureComponent` is a subclass of `Component` and supports [all the `Component` APIs.](/reference/react/Component#reference) Extending `PureComponent` is equivalent to defining a custom [`shouldComponentUpdate`](/reference/react/Component#shouldcomponentupdate) method that shallowly compares props and state.

[See more examples below.](#usage)

---

## Usage {/*usage*/}

### Skipping unnecessary re-renders for class components {/*skipping-unnecessary-re-renders-for-class-components*/}

React normally re-renders a component whenever its parent re-renders. As an optimization, you can create a component that React will not re-render when its parent re-renders so long as its new props and state are the same as the old props and state. [Class components](/reference/react/Component) can opt into this behavior by extending `PureComponent`:

```js {1}
class Greeting extends PureComponent {
 render() {
 return <h1>Hello, {this.props.name}!</h1>;
 }
}
```

A React component should always have [pure rendering logic.](/learn/keeping-components-pure) This means that it must return the same output if its props, state, and context haven't changed. By using `PureComponent`, you are telling React that your component complies with this requirement, so React doesn't need to re-render as long as its props and state haven't changed. However, your component will still re-render if a context that it's using changes.

In this example, notice that the `Greeting` component re-renders whenever `name` is changed (because that's one of its props), but not when `address` is changed (because it's not passed to `Greeting` as a prop):

<Sandpack>

```js
import { PureComponent, useState } from 'react';

class Greeting extends PureComponent {
 render() {
 console.log("Greeting was rendered at", new Date().toLocaleTimeString());
 return <h3>Hello{this.props.name && ', '}{this.props.name}!</h3>;
 }
}

export default function MyApp() {
 const [name, setName] = useState('');
 const [address, setAddress] = useState('');
 return (
 <>
 <label>
 Name{': '}
 <input value={name} onChange={e => setName(e.target.value)} />
 </label>
 <label>
 Address{': '}
 <input value={address} onChange={e => setAddress(e.target.value)} />
 </label>
 <Greeting name={name} />
 </>
 );
}
```

```css
label {
 display: block;
 margin-bottom: 16px;
}
```

</Sandpack>

<Pitfall>

We recommend defining components as functions instead of classes. [See how to migrate.](#alternatives)

</Pitfall>

---

## Alternatives {/*alternatives*/}

### Migrating from a `PureComponent` class component to a function {/*migrating-from-a-purecomponent-class-component-to-a-function*/}

We recommend using function components instead of [class components](/reference/react/Component) in new code. If you have some existing class components using `PureComponent`, here is how you can convert them. This is the original code:

<Sandpack>

```js
import { PureComponent, useState } from 'react';

class Greeting extends PureComponent {
 render() {
 console.log("Greeting was rendered at", new Date().toLocaleTimeString());
 return <h3>Hello{this.props.name && ', '}{this.props.name}!</h3>;
 }
}

export default function MyApp() {
 const [name, setName] = useState('');
 const [address, setAddress] = useState('');
 return (
 <>
 <label>
 Name{': '}
 <input value={name} onChange={e => setName(e.target.value)} />
 </label>
 <label>
 Address{': '}
 <input value={address} onChange={e => setAddress(e.target.value)} />
 </label>
 <Greeting name={name} />
 </>
 );
}
```

```css
label {
 display: block;
 margin-bottom: 16px;
}
```

</Sandpack>

When you [convert this component from a class to a function,](/reference/react/Component#alternatives) wrap it in [`memo`:](/reference/react/memo)

<Sandpack>

```js
import { memo, useState } from 'react';

const Greeting = memo(function Greeting({ name }) {
 console.log("Greeting was rendered at", new Date().toLocaleTimeString());
 return <h3>Hello{name && ', '}{name}!</h3>;
});

export default function MyApp() {
 const [name, setName] = useState('');
 const [address, setAddress] = useState('');
 return (
 <>
 <label>
 Name{': '}
 <input value={name} onChange={e => setName(e.target.value)} />
 </label>
 <label>
 Address{': '}
 <input value={address} onChange={e => setAddress(e.target.value)} />
 </label>
 <Greeting name={name} />
 </>
 );
}
```

```css
label {
 display: block;
 margin-bottom: 16px;
}
```

</Sandpack>

<Note>

Unlike `PureComponent`, [`memo`](/reference/react/memo) does not compare the new and the old state. In function components, calling the [`set` function](/reference/react/useState#setstate) with the same state [already prevents re-renders by default,](/reference/react/memo#updating-a-memoized-component-using-state) even without `memo`.

</Note>

---
title: <StrictMode>
---

<Intro>

`<StrictMode>` lets you find common bugs in your components early during development.

```js
<StrictMode>
 <App />
</StrictMode>
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `<StrictMode>` {/*strictmode*/}

Use `StrictMode` to enable additional development behaviors and warnings for the component tree inside:

```js
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

const root = createRoot(document.getElementById('root'));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

[See more examples below.](#usage)

Strict Mode enables the following development-only behaviors:

- Your components will [re-render an extra time](#fixing-bugs-found-by-double-rendering-in-development) to find bugs caused by impure rendering.
- Your components will [re-run Effects an extra time](#fixing-bugs-found-by-re-running-effects-in-development) to find bugs caused by missing Effect cleanup.
- Your components will [re-run refs callbacks an extra time](#fixing-bugs-found-by-re-running-ref-callbacks-in-development) to find bugs caused by missing ref cleanup.
- Your components will [be checked for usage of deprecated APIs.](#fixing-deprecation-warnings-enabled-by-strict-mode)

#### Props {/*props*/}

`StrictMode` accepts no props.

#### Caveats {/*caveats*/}

* There is no way to opt out of Strict Mode inside a tree wrapped in `<StrictMode>`. This gives you confidence that all components inside `<StrictMode>` are checked. If two teams working on a product disagree whether they find the checks valuable, they need to either reach consensus or move `<StrictMode>` down in the tree.

---

## Usage {/*usage*/}

### Enabling Strict Mode for entire app {/*enabling-strict-mode-for-entire-app*/}

Strict Mode enables extra development-only checks for the entire component tree inside the `<StrictMode>` component. These checks help you find common bugs in your components early in the development process.

To enable Strict Mode for your entire app, wrap your root component with `<StrictMode>` when you render it:

```js {6,8}
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

const root = createRoot(document.getElementById('root'));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

We recommend wrapping your entire app in Strict Mode, especially for newly created apps. If you use a framework that calls [`createRoot`](/reference/react-dom/client/createRoot) for you, check its documentation for how to enable Strict Mode.

Although the Strict Mode checks **only run in development,** they help you find bugs that already exist in your code but can be tricky to reliably reproduce in production. Strict Mode lets you fix bugs before your users report them.

<Note>

Strict Mode enables the following checks in development:

- Your components will [re-render an extra time](#fixing-bugs-found-by-double-rendering-in-development) to find bugs caused by impure rendering.
- Your components will [re-run Effects an extra time](#fixing-bugs-found-by-re-running-effects-in-development) to find bugs caused by missing Effect cleanup.
- Your components will [re-run ref callbacks an extra time](#fixing-bugs-found-by-re-running-ref-callbacks-in-development) to find bugs caused by missing ref cleanup.
- Your components will [be checked for usage of deprecated APIs.](#fixing-deprecation-warnings-enabled-by-strict-mode)

**All of these checks are development-only and do not impact the production build.**

</Note>

---

### Enabling Strict Mode for a part of the app {/*enabling-strict-mode-for-a-part-of-the-app*/}

You can also enable Strict Mode for any part of your application:

```js {7,12}
import { StrictMode } from 'react';

function App() {
 return (
 <>
 <Header />
 <StrictMode>
 <main>
 <Sidebar />
 <Content />
 </main>
 </StrictMode>
 <Footer />
 </>
 );
}
```

In this example, Strict Mode checks will not run against the `Header` and `Footer` components. However, they will run on `Sidebar` and `Content`, as well as all of the components inside them, no matter how deep.

<Note>

When `StrictMode` is enabled for a part of the app, React will only enable behaviors that are possible in production. For example, if `<StrictMode>` is not enabled at the root of the app, it will not [re-run Effects an extra time](#fixing-bugs-found-by-re-running-effects-in-development) on initial mount, since this would cause child effects to double fire without the parent effects, which cannot happen in production.

</Note>

---

### Fixing bugs found by double rendering in development {/*fixing-bugs-found-by-double-rendering-in-development*/}

[React assumes that every component you write is a pure function.](/learn/keeping-components-pure) This means that React components you write must always return the same JSX given the same inputs (props, state, and context).

Components breaking this rule behave unpredictably and cause bugs. To help you find accidentally impure code, Strict Mode calls some of your functions (only the ones that should be pure) **twice in development.** This includes:

- Your component function body (only top-level logic, so this doesn't include code inside event handlers)
- Functions that you pass to [`useState`](/reference/react/useState), [`set` functions](/reference/react/useState#setstate), [`useMemo`](/reference/react/useMemo), or [`useReducer`](/reference/react/useReducer)
- Some class component methods like [`constructor`](/reference/react/Component#constructor), [`render`](/reference/react/Component#render), [`shouldComponentUpdate`](/reference/react/Component#shouldcomponentupdate) ([see the whole list](https://reactjs.org/docs/strict-mode.html#detecting-unexpected-side-effects))

If a function is pure, running it twice does not change its behavior because a pure function produces the same result every time. However, if a function is impure (for example, it mutates the data it receives), running it twice tends to be noticeable (that's what makes it impure!) This helps you spot and fix the bug early.

**Here is an example to illustrate how double rendering in Strict Mode helps you find bugs early.**

This `StoryTray` component takes an array of `stories` and adds one last "Create Story" item at the end:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
root.render(<App />);
```

```js src/App.js
import { useState } from 'react';
import StoryTray from './StoryTray.js';

let initialStories = [
 {id: 0, label: "Ankit's Story" },
 {id: 1, label: "Taylor's Story" },
];

export default function App() {
 let [stories, setStories] = useState(initialStories)
 return (
 <div
 style={{
 width: '100%',
 height: '100%',
 textAlign: 'center',
 }}
 >
 <StoryTray stories={stories} />
 </div>
 );
}
```

```js src/StoryTray.js active
export default function StoryTray({ stories }) {
 const items = stories;
 items.push({ id: 'create', label: 'Create Story' });
 return (
 <ul>
 {items.map(story => (
 <li key={story.id}>
 {story.label}
 </li>
 ))}
 </ul>
 );
}
```

```css
ul {
 margin: 0;
 list-style-type: none;
 height: 100%;
 display: flex;
 flex-wrap: wrap;
 padding: 10px;
}

li {
 border: 1px solid #aaa;
 border-radius: 6px;
 float: left;
 margin: 5px;
 padding: 5px;
 width: 70px;
 height: 100px;
}
```

</Sandpack>

There is a mistake in the code above. However, it is easy to miss because the initial output appears correct.

This mistake will become more noticeable if the `StoryTray` component re-renders multiple times. For example, let's make the `StoryTray` re-render with a different background color whenever you hover over it:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById('root'));
root.render(<App />);
```

```js src/App.js
import { useState } from 'react';
import StoryTray from './StoryTray.js';

let initialStories = [
 {id: 0, label: "Ankit's Story" },
 {id: 1, label: "Taylor's Story" },
];

export default function App() {
 let [stories, setStories] = useState(initialStories)
 return (
 <div
 style={{
 width: '100%',
 height: '100%',
 textAlign: 'center',
 }}
 >
 <StoryTray stories={stories} />
 </div>
 );
}
```

```js src/StoryTray.js active
import { useState } from 'react';

export default function StoryTray({ stories }) {
 const [isHover, setIsHover] = useState(false);
 const items = stories;
 items.push({ id: 'create', label: 'Create Story' });
 return (
 <ul
 onPointerEnter={() => setIsHover(true)}
 onPointerLeave={() => setIsHover(false)}
 style={{
 backgroundColor: isHover ? '#ddd' : '#fff'
 }}
 >
 {items.map(story => (
 <li key={story.id}>
 {story.label}
 </li>
 ))}
 </ul>
 );
}
```

```css
ul {
 margin: 0;
 list-style-type: none;
 height: 100%;
 display: flex;
 flex-wrap: wrap;
 padding: 10px;
}

li {
 border: 1px solid #aaa;
 border-radius: 6px;
 float: left;
 margin: 5px;
 padding: 5px;
 width: 70px;
 height: 100px;
}
```

</Sandpack>

Notice how every time you hover over the `StoryTray` component, "Create Story" gets added to the list again. The intention of the code was to add it once at the end. But `StoryTray` directly modifies the `stories` array from the props. Every time `StoryTray` renders, it adds "Create Story" again at the end of the same array. In other words, `StoryTray` is not a pure function--running it multiple times produces different results.

To fix this problem, you can make a copy of the array, and modify that copy instead of the original one:

```js {2}
export default function StoryTray({ stories }) {
 const items = stories.slice(); // Clone the array
 // ✅ Good: Pushing into a new array
 items.push({ id: 'create', label: 'Create Story' });
```

This would [make the `StoryTray` function pure.](/learn/keeping-components-pure) Each time it is called, it would only modify a new copy of the array, and would not affect any external objects or variables. This solves the bug, but you had to make the component re-render more often before it became obvious that something is wrong with its behavior.

**In the original example, the bug wasn't obvious. Now let's wrap the original (buggy) code in `<StrictMode>`:**

<Sandpack>

```js src/index.js
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```js src/App.js
import { useState } from 'react';
import StoryTray from './StoryTray.js';

let initialStories = [
 {id: 0, label: "Ankit's Story" },
 {id: 1, label: "Taylor's Story" },
];

export default function App() {
 let [stories, setStories] = useState(initialStories)
 return (
 <div
 style={{
 width: '100%',
 height: '100%',
 textAlign: 'center',
 }}
 >
 <StoryTray stories={stories} />
 </div>
 );
}
```

```js src/StoryTray.js active
export default function StoryTray({ stories }) {
 const items = stories;
 items.push({ id: 'create', label: 'Create Story' });
 return (
 <ul>
 {items.map(story => (
 <li key={story.id}>
 {story.label}
 </li>
 ))}
 </ul>
 );
}
```

```css
ul {
 margin: 0;
 list-style-type: none;
 height: 100%;
 display: flex;
 flex-wrap: wrap;
 padding: 10px;
}

li {
 border: 1px solid #aaa;
 border-radius: 6px;
 float: left;
 margin: 5px;
 padding: 5px;
 width: 70px;
 height: 100px;
}
```

</Sandpack>

**Strict Mode *always* calls your rendering function twice, so you can see the mistake right away** ("Create Story" appears twice). This lets you notice such mistakes early in the process. When you fix your component to render in Strict Mode, you *also* fix many possible future production bugs like the hover functionality from before:

<Sandpack>

```js src/index.js
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById('root'));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```js src/App.js
import { useState } from 'react';
import StoryTray from './StoryTray.js';

let initialStories = [
 {id: 0, label: "Ankit's Story" },
 {id: 1, label: "Taylor's Story" },
];

export default function App() {
 let [stories, setStories] = useState(initialStories)
 return (
 <div
 style={{
 width: '100%',
 height: '100%',
 textAlign: 'center',
 }}
 >
 <StoryTray stories={stories} />
 </div>
 );
}
```

```js src/StoryTray.js active
import { useState } from 'react';

export default function StoryTray({ stories }) {
 const [isHover, setIsHover] = useState(false);
 const items = stories.slice(); // Clone the array
 items.push({ id: 'create', label: 'Create Story' });
 return (
 <ul
 onPointerEnter={() => setIsHover(true)}
 onPointerLeave={() => setIsHover(false)}
 style={{
 backgroundColor: isHover ? '#ddd' : '#fff'
 }}
 >
 {items.map(story => (
 <li key={story.id}>
 {story.label}
 </li>
 ))}
 </ul>
 );
}
```

```css
ul {
 margin: 0;
 list-style-type: none;
 height: 100%;
 display: flex;
 flex-wrap: wrap;
 padding: 10px;
}

li {
 border: 1px solid #aaa;
 border-radius: 6px;
 float: left;
 margin: 5px;
 padding: 5px;
 width: 70px;
 height: 100px;
}
```

</Sandpack>

Without Strict Mode, it was easy to miss the bug until you added more re-renders. Strict Mode made the same bug appear right away. Strict Mode helps you find bugs before you push them to your team and to your users.

[Read more about keeping components pure.](/learn/keeping-components-pure)

<Note>

If you have [React DevTools](/learn/react-developer-tools) installed, any `console.log` calls during the second render call will appear slightly dimmed. React DevTools also offers a setting (off by default) to suppress them completely.

</Note>

---

### Fixing bugs found by re-running Effects in development {/*fixing-bugs-found-by-re-running-effects-in-development*/}

Strict Mode can also help find bugs in [Effects.](/learn/synchronizing-with-effects)

Every Effect has some setup code and may have some cleanup code. Normally, React calls setup when the component *mounts* (is added to the screen) and calls cleanup when the component *unmounts* (is removed from the screen). React then calls cleanup and setup again if its dependencies changed since the last render.

When Strict Mode is on, React will also run **one extra setup+cleanup cycle in development for every Effect.** This may feel surprising, but it helps reveal subtle bugs that are hard to catch manually.

**Here is an example to illustrate how re-running Effects in Strict Mode helps you find bugs early.**

Consider this example that connects a component to a chat:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
root.render(<App />);
```

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';
const roomId = 'general';

export default function ChatRoom() {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 }, []);
 return <h1>Welcome to the {roomId} room!</h1>;
}
```

```js src/chat.js
let connections = 0;

export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 connections++;
 console.log('Active connections: ' + connections);
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 connections--;
 console.log('Active connections: ' + connections);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

There is an issue with this code, but it might not be immediately clear.

To make the issue more obvious, let's implement a feature. In the example below, `roomId` is not hardcoded. Instead, the user can select the `roomId` that they want to connect to from a dropdown. Click "Open chat" and then select different chat rooms one by one. Keep track of the number of active connections in the console:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
root.render(<App />);
```

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';

function ChatRoom({ roomId }) {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 }, [roomId]);

 return <h1>Welcome to the {roomId} room!</h1>;
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/chat.js
let connections = 0;

export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 connections++;
 console.log('Active connections: ' + connections);
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 connections--;
 console.log('Active connections: ' + connections);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

You'll notice that the number of open connections always keeps growing. In a real app, this would cause performance and network problems. The issue is that [your Effect is missing a cleanup function:](/learn/synchronizing-with-effects#step-3-add-cleanup-if-needed)

```js {4}
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]);
```

Now that your Effect "cleans up" after itself and destroys the outdated connections, the leak is solved. However, notice that the problem did not become visible until you've added more features (the select box).

**In the original example, the bug wasn't obvious. Now let's wrap the original (buggy) code in `<StrictMode>`:**

<Sandpack>

```js src/index.js
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';
const roomId = 'general';

export default function ChatRoom() {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 }, []);
 return <h1>Welcome to the {roomId} room!</h1>;
}
```

```js src/chat.js
let connections = 0;

export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 connections++;
 console.log('Active connections: ' + connections);
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 connections--;
 console.log('Active connections: ' + connections);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

**With Strict Mode, you immediately see that there is a problem** (the number of active connections jumps to 2). Strict Mode runs an extra setup+cleanup cycle for every Effect. This Effect has no cleanup logic, so it creates an extra connection but doesn't destroy it. This is a hint that you're missing a cleanup function.

Strict Mode lets you notice such mistakes early in the process. When you fix your Effect by adding a cleanup function in Strict Mode, you *also* fix many possible future production bugs like the select box from before:

<Sandpack>

```js src/index.js
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';

function ChatRoom({ roomId }) {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]);

 return <h1>Welcome to the {roomId} room!</h1>;
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/chat.js
let connections = 0;

export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 connections++;
 console.log('Active connections: ' + connections);
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 connections--;
 console.log('Active connections: ' + connections);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

Notice how the active connection count in the console doesn't keep growing anymore.

Without Strict Mode, it was easy to miss that your Effect needed cleanup. By running *setup → cleanup → setup* instead of *setup* for your Effect in development, Strict Mode made the missing cleanup logic more noticeable.

[Read more about implementing Effect cleanup.](/learn/synchronizing-with-effects#how-to-handle-the-effect-firing-twice-in-development)

---
### Fixing bugs found by re-running ref callbacks in development {/*fixing-bugs-found-by-re-running-ref-callbacks-in-development*/}

Strict Mode can also help find bugs in [callbacks refs.](/learn/manipulating-the-dom-with-refs)

Every callback `ref` has some setup code and may have some cleanup code. Normally, React calls setup when the element is *created* (is added to the DOM) and calls cleanup when the element is *removed* (is removed from the DOM).

When Strict Mode is on, React will also run **one extra setup+cleanup cycle in development for every callback `ref`.** This may feel surprising, but it helps reveal subtle bugs that are hard to catch manually.

Consider this example, which allows you to select an animal and then scroll to one of them. Notice when you switch from "Cats" to "Dogs", the console logs show that the number of animals in the list keeps growing, and the "Scroll to" buttons stop working:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
// ❌ Not using StrictMode.
root.render(<App />);
```

```js src/App.js active
import { useRef, useState } from "react";

export default function CatFriends() {
 const itemsRef = useRef([]);
 const [catList, setCatList] = useState(setupCatList);
 const [cat, setCat] = useState('neo');

 function scrollToCat(index) {
 const list = itemsRef.current;
 const {node} = list[index];
 node.scrollIntoView({
 behavior: "smooth",
 block: "nearest",
 inline: "center",
 });
 }

 const cats = catList.filter(c => c.type === cat)

 return (
 <>
 <nav>
 <button onClick={() => setCat('neo')}>Neo</button>
 <button onClick={() => setCat('millie')}>Millie</button>
 </nav>
 <hr />
 <nav>
 <span>Scroll to:</span>{cats.map((cat, index) => (
 <button key={cat.src} onClick={() => scrollToCat(index)}>
 {index}
 </button>
 ))}
 </nav>
 <div>
 <ul>
 {cats.map((cat) => (
 <li
 key={cat.src}
 ref={(node) => {
 const list = itemsRef.current;
 const item = {cat: cat, node};
 list.push(item);
 console.log(`✅ Adding cat to the map. Total cats: ${list.length}`);
 if (list.length > 10) {
 console.log('❌ Too many cats in the list!');
 }
 return () => {
 // 🚩 No cleanup, this is a bug!
 }
 }}
 >
 <img src={cat.src} />
 </li>
 ))}
 </ul>
 </div>
 </>
 );
}

function setupCatList() {
 const catList = [];
 for (let i = 0; i < 10; i++) {
 catList.push({type: 'neo', src: "https://placecats.com/neo/320/240?" + i});
 }
 for (let i = 0; i < 10; i++) {
 catList.push({type: 'millie', src: "https://placecats.com/millie/320/240?" + i});
 }

 return catList;
}

```

```css
div {
 width: 100%;
 overflow: hidden;
}

nav {
 text-align: center;
}

button {
 margin: .25rem;
}

ul,
li {
 list-style: none;
 white-space: nowrap;
}

li {
 display: inline;
 padding: 0.5rem;
}
```

</Sandpack>

**This is a production bug!** Since the ref callback doesn't remove animals from the list in the cleanup, the list of animals keeps growing. This is a memory leak that can cause performance problems in a real app, and breaks the behavior of the app.

The issue is the ref callback doesn't cleanup after itself:

```js {6-8}
<li
 ref={node => {
 const list = itemsRef.current;
 const item = {animal, node};
 list.push(item);
 return () => {
 // 🚩 No cleanup, this is a bug!
 }
 }}
</li>
```

Now let's wrap the original (buggy) code in `<StrictMode>`:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import {StrictMode} from 'react';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
// ✅ Using StrictMode.
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```js src/App.js active
import { useRef, useState } from "react";

export default function CatFriends() {
 const itemsRef = useRef([]);
 const [catList, setCatList] = useState(setupCatList);
 const [cat, setCat] = useState('neo');

 function scrollToCat(index) {
 const list = itemsRef.current;
 const {node} = list[index];
 node.scrollIntoView({
 behavior: "smooth",
 block: "nearest",
 inline: "center",
 });
 }

 const cats = catList.filter(c => c.type === cat)

 return (
 <>
 <nav>
 <button onClick={() => setCat('neo')}>Neo</button>
 <button onClick={() => setCat('millie')}>Millie</button>
 </nav>
 <hr />
 <nav>
 <span>Scroll to:</span>{cats.map((cat, index) => (
 <button key={cat.src} onClick={() => scrollToCat(index)}>
 {index}
 </button>
 ))}
 </nav>
 <div>
 <ul>
 {cats.map((cat) => (
 <li
 key={cat.src}
 ref={(node) => {
 const list = itemsRef.current;
 const item = {cat: cat, node};
 list.push(item);
 console.log(`✅ Adding cat to the map. Total cats: ${list.length}`);
 if (list.length > 10) {
 console.log('❌ Too many cats in the list!');
 }
 return () => {
 // 🚩 No cleanup, this is a bug!
 }
 }}
 >
 <img src={cat.src} />
 </li>
 ))}
 </ul>
 </div>
 </>
 );
}

function setupCatList() {
 const catList = [];
 for (let i = 0; i < 10; i++) {
 catList.push({type: 'neo', src: "https://placecats.com/neo/320/240?" + i});
 }
 for (let i = 0; i < 10; i++) {
 catList.push({type: 'millie', src: "https://placecats.com/millie/320/240?" + i});
 }

 return catList;
}

```

```css
div {
 width: 100%;
 overflow: hidden;
}

nav {
 text-align: center;
}

button {
 margin: .25rem;
}

ul,
li {
 list-style: none;
 white-space: nowrap;
}

li {
 display: inline;
 padding: 0.5rem;
}
```

</Sandpack>

**With Strict Mode, you immediately see that there is a problem**. Strict Mode runs an extra setup+cleanup cycle for every callback ref. This callback ref has no cleanup logic, so it adds refs but doesn't remove them. This is a hint that you're missing a cleanup function.

Strict Mode lets you eagerly find mistakes in callback refs. When you fix your callback by adding a cleanup function in Strict Mode, you *also* fix many possible future production bugs like the "Scroll to" bug from before:

<Sandpack>

```js src/index.js
import { createRoot } from 'react-dom/client';
import {StrictMode} from 'react';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById("root"));
// ✅ Using StrictMode.
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```js src/App.js active
import { useRef, useState } from "react";

export default function CatFriends() {
 const itemsRef = useRef([]);
 const [catList, setCatList] = useState(setupCatList);
 const [cat, setCat] = useState('neo');

 function scrollToCat(index) {
 const list = itemsRef.current;
 const {node} = list[index];
 node.scrollIntoView({
 behavior: "smooth",
 block: "nearest",
 inline: "center",
 });
 }

 const cats = catList.filter(c => c.type === cat)

 return (
 <>
 <nav>
 <button onClick={() => setCat('neo')}>Neo</button>
 <button onClick={() => setCat('millie')}>Millie</button>
 </nav>
 <hr />
 <nav>
 <span>Scroll to:</span>{cats.map((cat, index) => (
 <button key={cat.src} onClick={() => scrollToCat(index)}>
 {index}
 </button>
 ))}
 </nav>
 <div>
 <ul>
 {cats.map((cat) => (
 <li
 key={cat.src}
 ref={(node) => {
 const list = itemsRef.current;
 const item = {cat: cat, node};
 list.push(item);
 console.log(`✅ Adding cat to the map. Total cats: ${list.length}`);
 if (list.length > 10) {
 console.log('❌ Too many cats in the list!');
 }
 return () => {
 list.splice(list.indexOf(item), 1);
 console.log(`❌ Removing cat from the map. Total cats: ${itemsRef.current.length}`);
 }
 }}
 >
 <img src={cat.src} />
 </li>
 ))}
 </ul>
 </div>
 </>
 );
}

function setupCatList() {
 const catList = [];
 for (let i = 0; i < 10; i++) {
 catList.push({type: 'neo', src: "https://placecats.com/neo/320/240?" + i});
 }
 for (let i = 0; i < 10; i++) {
 catList.push({type: 'millie', src: "https://placecats.com/millie/320/240?" + i});
 }

 return catList;
}

```

```css
div {
 width: 100%;
 overflow: hidden;
}

nav {
 text-align: center;
}

button {
 margin: .25rem;
}

ul,
li {
 list-style: none;
 white-space: nowrap;
}

li {
 display: inline;
 padding: 0.5rem;
}
```

</Sandpack>

Now on inital mount in StrictMode, the ref callbacks are all setup, cleaned up, and setup again:

```
...
✅ Adding animal to the map. Total animals: 10
...
❌ Removing animal from the map. Total animals: 0
...
✅ Adding animal to the map. Total animals: 10
```

**This is expected.** Strict Mode confirms that the ref callbacks are cleaned up correctly, so the size never grows above the expected amount. After the fix, there are no memory leaks, and all the features work as expected.

Without Strict Mode, it was easy to miss the bug until you clicked around to app to notice broken features. Strict Mode made the bugs appear right away, before you push them to production.

---
### Fixing deprecation warnings enabled by Strict Mode {/*fixing-deprecation-warnings-enabled-by-strict-mode*/}

React warns if some component anywhere inside a `<StrictMode>` tree uses one of these deprecated APIs:

* `UNSAFE_` class lifecycle methods like [`UNSAFE_componentWillMount`](/reference/react/Component#unsafe_componentwillmount). [See alternatives.](https://reactjs.org/blog/2018/03/27/update-on-async-rendering.html#migrating-from-legacy-lifecycles)

These APIs are primarily used in older [class components](/reference/react/Component) so they rarely appear in modern apps.

---
title: <Suspense>
---

<Intro>

`<Suspense>` lets you display a fallback until its children have finished loading.

```js
<Suspense fallback={<Loading />}>
 <SomeComponent />
</Suspense>
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `<Suspense>` {/*suspense*/}

#### Props {/*props*/}
* `children`: The actual UI you intend to render. If `children` suspends while rendering, the Suspense boundary will switch to rendering `fallback`.
* `fallback`: An alternate UI to render in place of the actual UI if it has not finished loading. Any valid React node is accepted, though in practice, a fallback is a lightweight placeholder view, such as a loading spinner or skeleton. Suspense will automatically switch to `fallback` when `children` suspends, and back to `children` when the data is ready. If `fallback` suspends while rendering, it will activate the closest parent Suspense boundary.
* <ExperimentalBadge /> **optional** `defer`: A boolean. When `true`, React may show the `fallback` first and render or stream `children` later, even when nothing in them suspends. Use it for content that is expensive to render. Defaults to `false`.

#### Caveats {/*caveats*/}

- Suspense does not detect when data is fetched inside an Effect or event handler. It only activates in the [cases listed below.](#what-activates-a-suspense-boundary)
- React does not preserve any state for renders that got suspended before they were able to mount for the first time. When the component has loaded, React will retry rendering the suspended tree from scratch.
- If Suspense was displaying content for the tree, but then it suspended again, the `fallback` will be shown again unless the update causing it was caused by [`startTransition`](/reference/react/startTransition) or [`useDeferredValue`](/reference/react/useDeferredValue).
- React reveals suspended content at most once every 300ms, measured from the last reveal. Boundaries that become ready within that window are [revealed together](/blog/2025/10/01/react-19-2#batching-suspense-boundaries-for-ssr) rather than one at a time.
- If React needs to hide the already visible content because it suspended again, it will clean up [layout Effects](/reference/react/useLayoutEffect) in the content tree. When the content is ready to be shown again, React will fire the layout Effects again. This ensures that Effects measuring the DOM layout don't try to do this while the content is hidden.
- React includes under-the-hood optimizations like *Streaming Server Rendering* and *Selective Hydration* that are integrated with Suspense. Read [an architectural overview](https://github.com/reactwg/react-18/discussions/37) and watch [a technical talk](https://www.youtube.com/watch?v=pj5N-Khihgc) to learn more.

---

### What activates a Suspense boundary {/*what-activates-a-suspense-boundary*/}

A Suspense boundary waits for its content to be ready before revealing it. Any of the following keeps a boundary from revealing its content:

- Lazy-loading component code with [`lazy`](/reference/react/lazy).
- Reading a Promise with [`use`](/reference/react/use), including data streamed from [Server Components](/reference/rsc/server-components) or loaded through a [Suspense-enabled framework](#suspense-enabled-frameworks).
- Loading a stylesheet rendered with [`<link rel="stylesheet">` and a `precedence` prop.](/reference/react-dom/components/link#special-rendering-behavior) React blocks the boundary until the stylesheet loads, up to a timeout. [See an example below.](#waiting-for-a-stylesheet-to-load)
- Waiting for a large boundary's HTML to arrive during streaming server rendering. Sending HTML takes time, so a boundary with enough content activates even when nothing in it suspends. React reveals the content as the HTML arrives.
- <CanaryBadge /> Loading fonts. Suspense doesn't wait for fonts by default, but a [`<ViewTransition>`](/reference/react/ViewTransition) update waits for new fonts to load, up to a timeout, so text doesn't flash with a fallback font. [See an example below.](#waiting-for-a-font-to-load)
- <CanaryBadge /> Loading images. Suspense doesn't wait for images by default, but during a [`<ViewTransition>`](/reference/react/ViewTransition) update, React blocks the boundary until the image loads, up to a timeout. Adding an `onLoad` handler opts a specific image out. [See an example below.](#waiting-for-an-image-to-load)
- <ExperimentalBadge /> Performing CPU-bound render work inside a [`<Suspense defer>`](#props) boundary.

<Note>

#### Suspense-enabled frameworks {/*suspense-enabled-frameworks*/}

A *Suspense-enabled framework* gives you a way to read data in your component in a way that activates the closest Suspense boundary. The exact way you load your data depends on your framework, and you'll find the details in its documentation. Under the hood, a Suspense-enabled framework maintains a cache of Promises and calls [`use`](/reference/react/use) to suspend on a Promise.

Without a framework, you can read a Promise with `use` directly, as long as the Promise is [cached so the same instance is reused across renders.](/reference/react/use#caching-promises-for-client-components)

</Note>

---

## Usage {/*usage*/}

### Displaying a fallback while content is loading {/*displaying-a-fallback-while-content-is-loading*/}

You can wrap any part of your application with a Suspense boundary:

```js [[1, 1, "<Loading />"], [2, 2, "<Albums />"]]
<Suspense fallback={<Loading />}>
 <Albums />
</Suspense>
```

React will display your <CodeStep step={1}>loading fallback</CodeStep> until all the code and data needed by <CodeStep step={2}>the children</CodeStep> has been loaded.

In the example below, the `Albums` component *suspends* while fetching the list of albums. Until it's ready to render, React switches the closest Suspense boundary above to show the fallback--your `Loading` component. Then, when the data loads, React hides the `Loading` fallback and renders the `Albums` component with data.

<Sandpack>

```js src/App.js hidden
import { useState } from 'react';
import ArtistPage from './ArtistPage.js';

export default function App() {
 const [show, setShow] = useState(false);
 if (show) {
 return (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 } else {
 return (
 <button onClick={() => setShow(true)}>
 Open The Beatles artist page
 </button>
 );
 }
}
```

```js src/ArtistPage.js active
import { Suspense } from 'react';
import Albums from './Albums.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Suspense fallback={<Loading />}>
 <Albums artistId={artist.id} />
 </Suspense>
 </>
 );
}

function Loading() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else {
 throw Error('Not implemented');
 }
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

</Sandpack>

By contrast, code that fetches data outside of `use`, such as inside an Effect, does not activate the boundary:

<Sandpack>

```js src/App.js hidden
import { useState } from 'react';
import ArtistPage from './ArtistPage.js';

export default function App() {
 const [show, setShow] = useState(false);
 if (show) {
 return (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 } else {
 return (
 <button onClick={() => setShow(true)}>
 Open The Beatles artist page
 </button>
 );
 }
}
```

```js src/ArtistPage.js active
import { Suspense } from 'react';
import EffectAlbums from './EffectAlbums.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Suspense fallback={<Loading />}>
 <EffectAlbums artistId={artist.id} />
 </Suspense>
 </>
 );
}

function Loading() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/EffectAlbums.js
import { useState, useEffect } from 'react';
import { fetchData } from './data.js';

export default function EffectAlbums({ artistId }) {
 const [albums, setAlbums] = useState([]);

 useEffect(() => {
 let active = true;
 fetchData(`/${artistId}/albums`).then(result => {
 if (active) {
 setAlbums(result);
 }
 });
 return () => {
 active = false;
 };
 }, [artistId]);

 // Suspense can't see this fetch, so its fallback never
 // shows. The list stays empty until the data arrives.
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else {
 throw Error('Not implemented');
 }
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

</Sandpack>

During streaming server rendering, a boundary also activates while its HTML is still streaming in. With any streaming server rendering API, React sends [the shell](/reference/react-dom/server/renderToPipeableStream#specifying-what-goes-into-the-shell) with the `fallback` first, then streams in each boundary's HTML and swaps out its `fallback` as that content arrives. Press "Render the page" to watch the page stream in:

<Sandpack>

```js src/App.js hidden
```

```html public/index.html
<!DOCTYPE html>
<html lang="en">
<head>
 <meta charset="UTF-8" />
 <title>Streaming SSR</title>
</head>
<body>
 <button id="render">Render the page</button>
 <br /><br />
 <iframe id="container" style="width: 100%; height: 180px; border: 1px solid #aaa;"></iframe>
</body>
</html>
```

```js src/index.js
import { flushReadableStreamToFrame } from './demo-helpers.js';
import { Suspense, use } from 'react';
import { renderToReadableStream } from 'react-dom/server';

let posts = null;

function Posts() {
 const text = use(posts.promise);
 return <p>{text}</p>;
}

function ProfilePage() {
 return (
 <html>
 <body>
 <h1>Alice</h1>
 <p>Photographer and traveler.</p>
 <Suspense fallback={<p>⌛ Loading posts...</p>}>
 <Posts />
 </Suspense>
 </body>
 </html>
 );
}

async function main(frame) {
 posts = Promise.withResolvers();
 const stream = await renderToReadableStream(<ProfilePage />);

 // The posts resolve after the shell has streamed, so React
 // streams their HTML in and swaps out the fallback.
 setTimeout(() => {
 posts.resolve(
 'Just got back from two weeks along the coast. The drive ' +
 'was longer than expected, but every stop was worth it. ' +
 'A full write-up and more photos are coming soon.'
 );
 }, 1500);

 await flushReadableStreamToFrame(stream, frame);
}

document.getElementById('render').addEventListener('click', () => {
 main(document.getElementById('container'));
});
```

```js src/demo-helpers.js hidden
export async function flushReadableStreamToFrame(readable, frame) {
 const doc = frame.contentWindow.document;
 const decoder = new TextDecoder();
 for await (const chunk of readable) {
 doc.write(decoder.decode(chunk, { stream: true }));
 }
 doc.close();
}
```

</Sandpack>

---

### Revealing content together at once {/*revealing-content-together-at-once*/}

By default, the whole tree inside Suspense is treated as a single unit. For example, even if *only one* of these components suspends waiting for some data, *all* of them together will be replaced by the loading indicator:

```js {2-5}
<Suspense fallback={<Loading />}>
 <Biography />
 <Panel>
 <Albums />
 </Panel>
</Suspense>
```

Then, after all of them are ready to be displayed, they will all appear together at once.

In the example below, both `Biography` and `Albums` fetch some data. However, because they are grouped under a single Suspense boundary, these components always "pop in" together at the same time.

<Sandpack>

```js src/App.js hidden
import { useState } from 'react';
import ArtistPage from './ArtistPage.js';

export default function App() {
 const [show, setShow] = useState(false);
 if (show) {
 return (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 } else {
 return (
 <button onClick={() => setShow(true)}>
 Open The Beatles artist page
 </button>
 );
 }
}
```

```js src/ArtistPage.js active
import { Suspense } from 'react';
import Albums from './Albums.js';
import Biography from './Biography.js';
import Panel from './Panel.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Suspense fallback={<Loading />}>
 <Biography artistId={artist.id} />
 <Panel>
 <Albums artistId={artist.id} />
 </Panel>
 </Suspense>
 </>
 );
}

function Loading() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/Panel.js
export default function Panel({ children }) {
 return (
 <section className="panel">
 {children}
 </section>
 );
}
```

```js src/Biography.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Biography({ artistId }) {
 const bio = use(fetchData(`/${artistId}/bio`));
 return (
 <section>
 <p className="bio">{bio}</p>
 </section>
 );
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else if (url === '/the-beatles/bio') {
 return await getBio();
 } else {
 throw Error('Not implemented');
 }
}

async function getBio() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1500);
 });

 return `The Beatles were an English rock band,
 formed in Liverpool in 1960, that comprised
 John Lennon, Paul McCartney, George Harrison
 and Ringo Starr.`;
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

```css
.bio { font-style: italic; }

.panel {
 border: 1px solid #aaa;
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}
```

</Sandpack>

Components that load data don't have to be direct children of the Suspense boundary. For example, you can move `Biography` and `Albums` into a new `Details` component. This doesn't change the behavior. `Biography` and `Albums` share the same closest parent Suspense boundary, so their reveal is coordinated together.

```js {2,8-11}
<Suspense fallback={<Loading />}>
 <Details artistId={artist.id} />
</Suspense>

function Details({ artistId }) {
 return (
 <>
 <Biography artistId={artistId} />
 <Panel>
 <Albums artistId={artistId} />
 </Panel>
 </>
 );
}
```

---

### Revealing nested content as it loads {/*revealing-nested-content-as-it-loads*/}

When a component suspends, the closest parent Suspense component shows the fallback. This lets you nest multiple Suspense components to create a loading sequence. Each Suspense boundary's fallback will be filled in as the next level of content becomes available. For example, you can give the album list its own fallback:

```js {3,7}
<Suspense fallback={<BigSpinner />}>
 <Biography />
 <Suspense fallback={<AlbumsGlimmer />}>
 <Panel>
 <Albums />
 </Panel>
 </Suspense>
</Suspense>
```

With this change, displaying the `Biography` doesn't need to "wait" for the `Albums` to load.

The sequence will be:

1. If `Biography` hasn't loaded yet, `BigSpinner` is shown in place of the entire content area.
2. Once `Biography` finishes loading, `BigSpinner` is replaced by the content.
3. If `Albums` hasn't loaded yet, `AlbumsGlimmer` is shown in place of `Albums` and its parent `Panel`.
4. Finally, once `Albums` finishes loading, it replaces `AlbumsGlimmer`.

<Sandpack>

```js src/App.js hidden
import { useState } from 'react';
import ArtistPage from './ArtistPage.js';

export default function App() {
 const [show, setShow] = useState(false);
 if (show) {
 return (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 } else {
 return (
 <button onClick={() => setShow(true)}>
 Open The Beatles artist page
 </button>
 );
 }
}
```

```js src/ArtistPage.js active
import { Suspense } from 'react';
import Albums from './Albums.js';
import Biography from './Biography.js';
import Panel from './Panel.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Suspense fallback={<BigSpinner />}>
 <Biography artistId={artist.id} />
 <Suspense fallback={<AlbumsGlimmer />}>
 <Panel>
 <Albums artistId={artist.id} />
 </Panel>
 </Suspense>
 </Suspense>
 </>
 );
}

function BigSpinner() {
 return <h2>🌀 Loading...</h2>;
}

function AlbumsGlimmer() {
 return (
 <div className="glimmer-panel">
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 </div>
 );
}
```

```js src/Panel.js
export default function Panel({ children }) {
 return (
 <section className="panel">
 {children}
 </section>
 );
}
```

```js src/Biography.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Biography({ artistId }) {
 const bio = use(fetchData(`/${artistId}/bio`));
 return (
 <section>
 <p className="bio">{bio}</p>
 </section>
 );
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else if (url === '/the-beatles/bio') {
 return await getBio();
 } else {
 throw Error('Not implemented');
 }
}

async function getBio() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 return `The Beatles were an English rock band,
 formed in Liverpool in 1960, that comprised
 John Lennon, Paul McCartney, George Harrison
 and Ringo Starr.`;
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

```css
.bio { font-style: italic; }

.panel {
 border: 1px solid #aaa;
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-panel {
 border: 1px dashed #aaa;
 background: linear-gradient(90deg, rgba(221,221,221,1) 0%, rgba(255,255,255,1) 100%);
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-line {
 display: block;
 width: 60%;
 height: 20px;
 margin: 10px;
 border-radius: 4px;
 background: #f0f0f0;
}
```

</Sandpack>

Suspense boundaries let you coordinate which parts of your UI should always "pop in" together at the same time, and which parts should progressively reveal more content in a sequence of loading states. You can add, move, or delete Suspense boundaries in any place in the tree without affecting the rest of your app's behavior.

Don't put a Suspense boundary around every component. Suspense boundaries should not be more granular than the loading sequence that you want the user to experience. If you work with a designer, ask them where the loading states should be placed--it's likely that they've already included them in their design wireframes.

---

### Showing stale content while fresh content is loading {/*showing-stale-content-while-fresh-content-is-loading*/}

In this example, the `SearchResults` component suspends while fetching the search results. Type `"a"`, wait for the results, and then edit it to `"ab"`. The results for `"a"` will get replaced by the loading fallback.

<Sandpack>

```js src/App.js
import { Suspense, useState } from 'react';
import SearchResults from './SearchResults.js';

export default function App() {
 const [query, setQuery] = useState('');
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <SearchResults query={query} />
 </Suspense>
 </>
 );
}
```

```js src/SearchResults.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function SearchResults({ query }) {
 if (query === '') {
 return null;
 }
 const albums = use(fetchData(`/search?q=${query}`));
 if (albums.length === 0) {
 return <p>No matches for <i>"{query}"</i></p>;
 }
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/search?q=')) {
 return await getSearchResults(url.slice('/search?q='.length));
 } else {
 throw Error('Not implemented');
 }
}

async function getSearchResults(query) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 const allAlbums = [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];

 const lowerQuery = query.trim().toLowerCase();
 return allAlbums.filter(album => {
 const lowerTitle = album.title.toLowerCase();
 return (
 lowerTitle.startsWith(lowerQuery) ||
 lowerTitle.indexOf(' ' + lowerQuery) !== -1
 )
 });
}
```

```css
input { margin: 10px; }
```

</Sandpack>

A common alternative UI pattern is to *defer* updating the list and to keep showing the previous results until the new results are ready. The [`useDeferredValue`](/reference/react/useDeferredValue) Hook lets you pass a deferred version of the query down:

```js {3,11}
export default function App() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <SearchResults query={deferredQuery} />
 </Suspense>
 </>
 );
}
```

The `query` will update immediately, so the input will display the new value. However, the `deferredQuery` will keep its previous value until the data has loaded, so `SearchResults` will show the stale results for a bit.

To make it more obvious to the user, you can add a visual indication when the stale result list is displayed:

```js {2}
<div style={{
 opacity: query !== deferredQuery ? 0.5 : 1
}}>
 <SearchResults query={deferredQuery} />
</div>
```

Enter `"a"` in the example below, wait for the results to load, and then edit the input to `"ab"`. Notice how instead of the Suspense fallback, you now see the dimmed stale result list until the new results have loaded:

<Sandpack>

```js src/App.js
import { Suspense, useState, useDeferredValue } from 'react';
import SearchResults from './SearchResults.js';

export default function App() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 const isStale = query !== deferredQuery;
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <div style={{ opacity: isStale ? 0.5 : 1 }}>
 <SearchResults query={deferredQuery} />
 </div>
 </Suspense>
 </>
 );
}
```

```js src/SearchResults.js hidden
import {use} from 'react';
import { fetchData } from './data.js';

export default function SearchResults({ query }) {
 if (query === '') {
 return null;
 }
 const albums = use(fetchData(`/search?q=${query}`));
 if (albums.length === 0) {
 return <p>No matches for <i>"{query}"</i></p>;
 }
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/search?q=')) {
 return await getSearchResults(url.slice('/search?q='.length));
 } else {
 throw Error('Not implemented');
 }
}

async function getSearchResults(query) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 const allAlbums = [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];

 const lowerQuery = query.trim().toLowerCase();
 return allAlbums.filter(album => {
 const lowerTitle = album.title.toLowerCase();
 return (
 lowerTitle.startsWith(lowerQuery) ||
 lowerTitle.indexOf(' ' + lowerQuery) !== -1
 )
 });
}
```

```css
input { margin: 10px; }
```

</Sandpack>

<Note>

Both deferred values and [Transitions](#preventing-already-revealed-content-from-hiding) let you avoid showing Suspense fallback in favor of inline indicators. Transitions mark the whole update as non-urgent so they are typically used by frameworks and router libraries for navigation. Deferred values, on the other hand, are mostly useful in application code where you want to mark a part of UI as non-urgent and let it "lag behind" the rest of the UI.

</Note>

---

### Preventing already revealed content from hiding {/*preventing-already-revealed-content-from-hiding*/}

When a component suspends, the closest parent Suspense boundary switches to showing the fallback. This can lead to a jarring user experience if it was already displaying some content. Try pressing this button:

<Sandpack>

```js src/App.js
import { Suspense, useState } from 'react';
import IndexPage from './IndexPage.js';
import ArtistPage from './ArtistPage.js';
import Layout from './Layout.js';

export default function App() {
 return (
 <Suspense fallback={<BigSpinner />}>
 <Router />
 </Suspense>
 );
}

function Router() {
 const [page, setPage] = useState('/');

 function navigate(url) {
 setPage(url);
 }

 let content;
 if (page === '/') {
 content = (
 <IndexPage navigate={navigate} />
 );
 } else if (page === '/the-beatles') {
 content = (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 }
 return (
 <Layout>
 {content}
 </Layout>
 );
}

function BigSpinner() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/Layout.js
export default function Layout({ children }) {
 return (
 <div className="layout">
 <section className="header">
 Music Browser
 </section>
 <main>
 {children}
 </main>
 </div>
 );
}
```

```js src/IndexPage.js
export default function IndexPage({ navigate }) {
 return (
 <button onClick={() => navigate('/the-beatles')}>
 Open The Beatles artist page
 </button>
 );
}
```

```js src/ArtistPage.js
import { Suspense } from 'react';
import Albums from './Albums.js';
import Biography from './Biography.js';
import Panel from './Panel.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Biography artistId={artist.id} />
 <Suspense fallback={<AlbumsGlimmer />}>
 <Panel>
 <Albums artistId={artist.id} />
 </Panel>
 </Suspense>
 </>
 );
}

function AlbumsGlimmer() {
 return (
 <div className="glimmer-panel">
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 </div>
 );
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/Biography.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Biography({ artistId }) {
 const bio = use(fetchData(`/${artistId}/bio`));
 return (
 <section>
 <p className="bio">{bio}</p>
 </section>
 );
}
```

```js src/Panel.js
export default function Panel({ children }) {
 return (
 <section className="panel">
 {children}
 </section>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else if (url === '/the-beatles/bio') {
 return await getBio();
 } else {
 throw Error('Not implemented');
 }
}

async function getBio() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 return `The Beatles were an English rock band,
 formed in Liverpool in 1960, that comprised
 John Lennon, Paul McCartney, George Harrison
 and Ringo Starr.`;
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

```css
main {
 min-height: 200px;
 padding: 10px;
}

.layout {
 border: 1px solid black;
}

.header {
 background: #222;
 padding: 10px;
 text-align: center;
 color: white;
}

.bio { font-style: italic; }

.panel {
 border: 1px solid #aaa;
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-panel {
 border: 1px dashed #aaa;
 background: linear-gradient(90deg, rgba(221,221,221,1) 0%, rgba(255,255,255,1) 100%);
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-line {
 display: block;
 width: 60%;
 height: 20px;
 margin: 10px;
 border-radius: 4px;
 background: #f0f0f0;
}
```

</Sandpack>

When you pressed the button, the `Router` component rendered `ArtistPage` instead of `IndexPage`. A component inside `ArtistPage` suspended, so the closest Suspense boundary started showing the fallback. The closest Suspense boundary was near the root, so the whole site layout got replaced by `BigSpinner`.

To prevent this, you can mark the navigation state update as a *Transition* with [`startTransition`:](/reference/react/startTransition)

```js {5,7}
function Router() {
 const [page, setPage] = useState('/');

 function navigate(url) {
 startTransition(() => {
 setPage(url);
 });
 }
 // ...
```

This tells React that the state transition is not urgent, and it's better to keep showing the previous page instead of hiding any already revealed content. Now clicking the button "waits" for the `Biography` to load:

<Sandpack>

```js src/App.js
import { Suspense, startTransition, useState } from 'react';
import IndexPage from './IndexPage.js';
import ArtistPage from './ArtistPage.js';
import Layout from './Layout.js';

export default function App() {
 return (
 <Suspense fallback={<BigSpinner />}>
 <Router />
 </Suspense>
 );
}

function Router() {
 const [page, setPage] = useState('/');

 function navigate(url) {
 startTransition(() => {
 setPage(url);
 });
 }

 let content;
 if (page === '/') {
 content = (
 <IndexPage navigate={navigate} />
 );
 } else if (page === '/the-beatles') {
 content = (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 }
 return (
 <Layout>
 {content}
 </Layout>
 );
}

function BigSpinner() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/Layout.js
export default function Layout({ children }) {
 return (
 <div className="layout">
 <section className="header">
 Music Browser
 </section>
 <main>
 {children}
 </main>
 </div>
 );
}
```

```js src/IndexPage.js
export default function IndexPage({ navigate }) {
 return (
 <button onClick={() => navigate('/the-beatles')}>
 Open The Beatles artist page
 </button>
 );
}
```

```js src/ArtistPage.js
import { Suspense } from 'react';
import Albums from './Albums.js';
import Biography from './Biography.js';
import Panel from './Panel.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Biography artistId={artist.id} />
 <Suspense fallback={<AlbumsGlimmer />}>
 <Panel>
 <Albums artistId={artist.id} />
 </Panel>
 </Suspense>
 </>
 );
}

function AlbumsGlimmer() {
 return (
 <div className="glimmer-panel">
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 </div>
 );
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/Biography.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Biography({ artistId }) {
 const bio = use(fetchData(`/${artistId}/bio`));
 return (
 <section>
 <p className="bio">{bio}</p>
 </section>
 );
}
```

```js src/Panel.js
export default function Panel({ children }) {
 return (
 <section className="panel">
 {children}
 </section>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else if (url === '/the-beatles/bio') {
 return await getBio();
 } else {
 throw Error('Not implemented');
 }
}

async function getBio() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 return `The Beatles were an English rock band,
 formed in Liverpool in 1960, that comprised
 John Lennon, Paul McCartney, George Harrison
 and Ringo Starr.`;
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

```css
main {
 min-height: 200px;
 padding: 10px;
}

.layout {
 border: 1px solid black;
}

.header {
 background: #222;
 padding: 10px;
 text-align: center;
 color: white;
}

.bio { font-style: italic; }

.panel {
 border: 1px solid #aaa;
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-panel {
 border: 1px dashed #aaa;
 background: linear-gradient(90deg, rgba(221,221,221,1) 0%, rgba(255,255,255,1) 100%);
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-line {
 display: block;
 width: 60%;
 height: 20px;
 margin: 10px;
 border-radius: 4px;
 background: #f0f0f0;
}
```

</Sandpack>

A Transition doesn't wait for *all* content to load. It only waits long enough to avoid hiding already revealed content. For example, the website `Layout` was already revealed, so it would be bad to hide it behind a loading spinner. However, the nested `Suspense` boundary around `Albums` is new, so the Transition doesn't wait for it.

<Note>

Suspense-enabled routers are expected to wrap the navigation updates into Transitions by default.

</Note>

---

### Indicating that a Transition is happening {/*indicating-that-a-transition-is-happening*/}

In the above example, once you click the button, there is no visual indication that a navigation is in progress. To add an indicator, you can replace [`startTransition`](/reference/react/startTransition) with [`useTransition`](/reference/react/useTransition) which gives you a boolean `isPending` value. In the example below, it's used to change the website header styling while a Transition is happening:

<Sandpack>

```js src/App.js
import { Suspense, useState, useTransition } from 'react';
import IndexPage from './IndexPage.js';
import ArtistPage from './ArtistPage.js';
import Layout from './Layout.js';

export default function App() {
 return (
 <Suspense fallback={<BigSpinner />}>
 <Router />
 </Suspense>
 );
}

function Router() {
 const [page, setPage] = useState('/');
 const [isPending, startTransition] = useTransition();

 function navigate(url) {
 startTransition(() => {
 setPage(url);
 });
 }

 let content;
 if (page === '/') {
 content = (
 <IndexPage navigate={navigate} />
 );
 } else if (page === '/the-beatles') {
 content = (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 }
 return (
 <Layout isPending={isPending}>
 {content}
 </Layout>
 );
}

function BigSpinner() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/Layout.js
export default function Layout({ children, isPending }) {
 return (
 <div className="layout">
 <section className="header" style={{
 opacity: isPending ? 0.7 : 1
 }}>
 Music Browser
 </section>
 <main>
 {children}
 </main>
 </div>
 );
}
```

```js src/IndexPage.js
export default function IndexPage({ navigate }) {
 return (
 <button onClick={() => navigate('/the-beatles')}>
 Open The Beatles artist page
 </button>
 );
}
```

```js src/ArtistPage.js
import { Suspense } from 'react';
import Albums from './Albums.js';
import Biography from './Biography.js';
import Panel from './Panel.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Biography artistId={artist.id} />
 <Suspense fallback={<AlbumsGlimmer />}>
 <Panel>
 <Albums artistId={artist.id} />
 </Panel>
 </Suspense>
 </>
 );
}

function AlbumsGlimmer() {
 return (
 <div className="glimmer-panel">
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 </div>
 );
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/Biography.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Biography({ artistId }) {
 const bio = use(fetchData(`/${artistId}/bio`));
 return (
 <section>
 <p className="bio">{bio}</p>
 </section>
 );
}
```

```js src/Panel.js
export default function Panel({ children }) {
 return (
 <section className="panel">
 {children}
 </section>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else if (url === '/the-beatles/bio') {
 return await getBio();
 } else {
 throw Error('Not implemented');
 }
}

async function getBio() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 return `The Beatles were an English rock band,
 formed in Liverpool in 1960, that comprised
 John Lennon, Paul McCartney, George Harrison
 and Ringo Starr.`;
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

```css
main {
 min-height: 200px;
 padding: 10px;
}

.layout {
 border: 1px solid black;
}

.header {
 background: #222;
 padding: 10px;
 text-align: center;
 color: white;
}

.bio { font-style: italic; }

.panel {
 border: 1px solid #aaa;
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-panel {
 border: 1px dashed #aaa;
 background: linear-gradient(90deg, rgba(221,221,221,1) 0%, rgba(255,255,255,1) 100%);
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-line {
 display: block;
 width: 60%;
 height: 20px;
 margin: 10px;
 border-radius: 4px;
 background: #f0f0f0;
}
```

</Sandpack>

---

### Resetting Suspense boundaries on navigation {/*resetting-suspense-boundaries-on-navigation*/}

During a Transition, React avoids hiding already revealed content. However, when you navigate to *different* content, such as another user's profile, you'll want the boundary to show the fallback instead of the previous content. You can express this with a `key`:

```js
<ProfilePage key={queryParams.id} />
```

With a different `key`, React treats the profiles as different content and resets the Suspense boundary during navigation. The `key` can go on the boundary itself or on a component above it. Suspense-integrated routers should do this automatically.

In the example below, opening the profile page loads the first profile. Pressing "Bob" navigates to a different profile, and the `key` resets the boundary, so the fallback shows instead of the previous user's bio. Try removing the `key`: the previous bio stays visible while the next one loads:

<Sandpack>

```js src/App.js hidden
import { useState } from 'react';
import ProfilePage from './ProfilePage.js';

export default function App() {
 const [show, setShow] = useState(false);
 if (show) {
 return <ProfilePage />;
 }
 return (
 <button onClick={() => setShow(true)}>
 Open profile page
 </button>
 );
}
```

```js src/ProfilePage.js active
import { Suspense, useState, startTransition } from 'react';
import Bio from './Bio.js';
import { fetchBio } from './data.js';

export default function ProfilePage() {
 const [user, setUser] = useState(() => ({
 id: 'alice',
 bioPromise: fetchBio('alice'),
 }));
 function navigate(id) {
 startTransition(() => {
 setUser({ id, bioPromise: fetchBio(id) });
 });
 }
 return (
 <>
 <button onClick={() => navigate('alice')}>
 Alice
 </button>
 <button onClick={() => navigate('bob')}>
 Bob
 </button>
 <Suspense key={user.id} fallback={<p>⌛ Loading profile...</p>}>
 <Bio bioPromise={user.bioPromise} />
 </Suspense>
 </>
 );
}
```

```js src/Bio.js
import { use } from 'react';

export default function Bio({ bioPromise }) {
 const bio = use(bioPromise);
 return <p>{bio}</p>;
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.

export async function fetchBio(userId) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1500);
 });

 return userId === 'alice'
 ? 'Alice is a photographer and traveler.'
 : 'Bob collects vintage synthesizers.';
}
```

```css
button {
 margin-right: 8px;
}
```

</Sandpack>

---

### Providing a fallback for server errors and client-only content {/*providing-a-fallback-for-server-errors-and-client-only-content*/}

If you use one of the [streaming server rendering APIs](/reference/react-dom/server) (or a framework that relies on them), React will also use your `<Suspense>` boundaries to handle errors on the server. If a component throws an error on the server, React will not abort the server render. Instead, it will find the closest `<Suspense>` component above it and include its fallback (such as a spinner) into the generated server HTML. The user will see a spinner at first.

On the client, React will attempt to render the same component again. If it errors on the client too, React will throw the error and display the closest [Error Boundary.](/reference/react/Component#static-getderivedstatefromerror) However, if it does not error on the client, React will not display the error to the user since the content was eventually displayed successfully.

You can use this to opt out some components from rendering on the server. To do this, throw an error in the server environment and then wrap them in a `<Suspense>` boundary to replace their HTML with fallbacks:

```js
<Suspense fallback={<Loading />}>
 <Chat />
</Suspense>

function Chat() {
 if (typeof window === 'undefined') {
 throw Error('Chat should only render on the client.');
 }
 // ...
}
```

The server HTML will include the loading indicator. It will be replaced by the `Chat` component on the client.

---

### Waiting for a stylesheet to load {/*waiting-for-a-stylesheet-to-load*/}

A stylesheet rendered with [`<link rel="stylesheet">` and a `precedence` prop](/reference/react-dom/components/link#special-rendering-behavior) blocks the Suspense boundary until the stylesheet loads, up to a timeout, so the content doesn't appear unstyled.

In the example below, the `Card` component renders a stylesheet with `precedence`. Press "Show card": React shows the fallback until the stylesheet has loaded, and then reveals the card with its styles applied.

For comparison, the second button performs the same update without React, in a separate document. Nothing waits for the stylesheet, so the card's text appears in a fallback font first and then switches:

<Sandpack>

```js
import { Suspense, useState, startTransition } from 'react';
import { freshStylesheetUrl } from './styles.js';
import VanillaCard from './VanillaCard.js';

function Card({ href }) {
 return (
 <>
 <link rel="stylesheet" href={href} precedence="default" />
 <div className="fancy-card">This card uses a font from the stylesheet.</div>
 </>
 );
}

export default function App() {
 const [href, setHref] = useState(null);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setHref(freshStylesheetUrl());
 });
 }}>
 Show card
 </button>
 {href && (
 <Suspense fallback={<p>⌛ Loading styles...</p>}>
 <Card href={href} />
 </Suspense>
 )}
 <hr />
 <VanillaCard />
 </>
 );
}
```

```js src/VanillaCard.js
import { useRef } from 'react';
import { freshStylesheetUrl } from './styles.js';

export default function VanillaCard() {
 const ref = useRef(null);
 function show() {
 const doc = ref.current.contentWindow.document;
 doc.open();
 doc.write(`
 <style>
 body { margin: 0; }
 .fancy-card {
 padding: 20px;
 border-radius: 8px;
 color: white;
 font-family: 'Caveat', sans-serif;
 font-size: 24px;
 background: linear-gradient(135deg, #087ea4, #2b3491);
 }
 </style>
 <div class="fancy-card">This card uses a font from the stylesheet.</div>
 <link rel="stylesheet" href="${freshStylesheetUrl()}">
 `);
 doc.close();
 }
 return (
 <>
 <button onClick={show}>Show card (without React)</button>
 <iframe ref={ref} title="Vanilla card" className="vanilla-frame" />
 </>
 );
}
```

```js src/styles.js hidden
// Add a unique parameter so the stylesheet isn't cached,
// and every run shows the loading state.
export function freshStylesheetUrl() {
 return (
 'https://fonts.googleapis.com/css2?family=Caveat&display=swap' +
 '&t=' +
 Date.now()
 );
}
```

```css
#root {
 min-height: 300px;
}
button {
 margin-right: 8px;
}
hr {
 margin: 16px 0;
}
.fancy-card {
 margin-top: 1em;
 padding: 20px;
 border-radius: 8px;
 color: white;
 font-family: 'Caveat', sans-serif;
 font-size: 24px;
 background: linear-gradient(135deg, #087ea4, #2b3491);
}
.vanilla-frame {
 display: block;
 margin-top: 1em;
 border: none;
 width: 100%;
 height: 90px;
}
```

</Sandpack>

---

### <CanaryBadge /> Animating from Suspense content {/*animating-from-suspense-content*/}

Suspense composes with [`<ViewTransition>`](/reference/react/ViewTransition) to animate the swap from the fallback to the content. Wrap the boundary in a `<ViewTransition>`, and React treats the swap as an update, cross-fading between the fallback and the content by default:

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}

export function VideoPlaceholder() {
 const video = {image: 'loading'};
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title loading" />
 <div className="video-description loading" />
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition, Suspense} from 'react';
import {Video, VideoPlaceholder} from './Video';
import {useLazyVideoData} from './data';

function LazyVideo() {
 const video = useLazyVideoData();
 return <Video video={video} />;
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>
 {showItem ? (
 <ViewTransition>
 <Suspense fallback={<VideoPlaceholder />}>
 <LazyVideo />
 </Suspense>
 </ViewTransition>
 ) : null}
 </>
 );
}
```

```js src/data.js hidden
import {use} from 'react';

let cache = null;

function fetchVideo() {
 if (!cache) {
 cache = new Promise((resolve) => {
 setTimeout(() => {
 resolve({
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 });
 }, 1000);
 });
 }
 return cache;
}

export function useLazyVideoData() {
 return use(fetchVideo());
}
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.loading {
 background-image: linear-gradient(
 90deg,
 rgba(173, 216, 230, 0.3) 25%,
 rgba(135, 206, 250, 0.5) 50%,
 rgba(173, 216, 230, 0.3) 75%
 );
 background-size: 200% 100%;
 animation: shimmer 1.5s infinite;
}
@keyframes shimmer {
 0% {
 background-position: -200% 0;
 }
 100% {
 background-position: 200% 0;
 }
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-title.loading {
 height: 20px;
 width: 80px;
 border-radius: 0.5rem;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
 border-radius: 0.5rem;
}
.video-description.loading {
 height: 15px;
 width: 100px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

<Note>

Where you place the `<ViewTransition>` relative to the boundary determines whether the fallback and content cross-fade as one update or animate as separate exit and enter animations. You can also [customize the animation](/reference/react/ViewTransition#customizing-animations) with View Transition classes.

[Learn more about animating from Suspense content.](/reference/react/ViewTransition#animating-from-suspense-content)

</Note>

---

### <CanaryBadge /> Waiting for a font to load {/*waiting-for-a-font-to-load*/}

When a [`<ViewTransition>`](/reference/react/ViewTransition) animates a Suspense boundary's reveal, React waits for new fonts the content introduces, up to a timeout, so the text doesn't flash with a fallback font. This only happens during a `<ViewTransition>` update.

In the example below, the Suspense boundary is wrapped in a `<ViewTransition>`, and the `Quote` component suspends while its data loads. Rendering the quote starts its font download. React keeps the fallback visible until the font has loaded, so the quote appears already in its font.

For comparison, the second button performs the same update without React. Nothing waits for the font, so the text appears in a fallback font first and then switches:

<Sandpack>

```js
import { ViewTransition, Suspense, use, useState, startTransition } from 'react';
import { fetchQuote } from './data.js';
import { freshFontUrl } from './font.js';
import VanillaQuote from './VanillaQuote.js';

function Quote({ fontSrc }) {
 const quote = use(fetchQuote());
 return (
 <>
 <style href={fontSrc} precedence="default">
 {`@font-face {
 font-family: 'Fancy';
 src: url(${fontSrc}) format('truetype');
 font-display: swap;
 }`}
 </style>
 <p className="quote fancy">{quote}</p>
 </>
 );
}

export default function App() {
 const [fontSrc, setFontSrc] = useState(null);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setFontSrc(freshFontUrl());
 });
 }}>
 Show quote
 </button>
 {fontSrc && (
 <ViewTransition>
 <Suspense fallback={<p className="quote">⌛ Loading quote...</p>}>
 <Quote fontSrc={fontSrc} />
 </Suspense>
 </ViewTransition>
 )}
 <hr />
 <VanillaQuote />
 </>
 );
}
```

```js src/VanillaQuote.js
import { useRef } from 'react';
import { freshFontUrl } from './font.js';

export default function VanillaQuote() {
 const ref = useRef(null);
 function show() {
 const style = document.createElement('style');
 style.textContent = `@font-face {
 font-family: 'VanillaFancy';
 src: url(${freshFontUrl()}) format('truetype');
 font-display: swap;
 }`;
 document.head.appendChild(style);
 ref.current.innerHTML = `<p class="quote vanilla-fancy">The best way to predict the future is to invent it.</p>`;
 }
 return (
 <>
 <button onClick={show}>Show quote (without React)</button>
 <div ref={ref} />
 </>
 );
}
```

```js src/font.js hidden
// Add a unique parameter so the font isn't cached,
// and every run shows the loading state.
export function freshFontUrl() {
 return (
 'https://raw.githubusercontent.com/google/fonts/main/ofl/caveat/Caveat%5Bwght%5D.ttf' +
 '?t=' +
 Date.now()
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = null;

export function fetchQuote() {
 if (!cache) {
 cache = new Promise((resolve) => {
 // Add a fake delay to make waiting noticeable.
 setTimeout(() => {
 resolve(
 'The best way to predict the future is to invent it.'
 );
 }, 500);
 });
 }
 return cache;
}
```

```css
#root {
 min-height: 260px;
}
.quote {
 font-size: 20px;
 margin-top: 1em;
}
.fancy {
 font-family: 'Fancy', sans-serif;
}
.vanilla-fancy {
 font-family: 'VanillaFancy', sans-serif;
}
hr {
 margin: 16px 0;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### <CanaryBadge /> Waiting for an image to load {/*waiting-for-an-image-to-load*/}

When a [`<ViewTransition>`](/reference/react/ViewTransition) animates a Suspense boundary's reveal, React waits for visible images to load, up to a timeout, so the animation doesn't start with a half-loaded image. This only happens during a `<ViewTransition>` update. Adding an `onLoad` handler opts a specific image out, even inside a `<ViewTransition>`.

In the example below, the Suspense boundary is wrapped in a `<ViewTransition>` and shows a profile skeleton until the portrait has loaded.

For comparison, the second button performs the same update without React. Nothing waits for the image, so the card appears immediately and the image pops in when it loads:

<Sandpack>

```js
import { ViewTransition, Suspense, useState, startTransition } from 'react';
import { freshImageUrl } from './image.js';
import VanillaProfile from './VanillaProfile.js';

function Profile({ src }) {
 return (
 <div className="card">
 <img src={src} alt="Jack Pope" width={80} height={80} />
 <p>Jack Pope</p>
 </div>
 );
}

function ProfilePlaceholder() {
 return (
 <div className="card">
 <div className="avatar-placeholder" />
 <p className="name-placeholder">&nbsp;</p>
 </div>
 );
}

export default function App() {
 const [src, setSrc] = useState(null);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setSrc(freshImageUrl());
 });
 }}>
 Show profile
 </button>
 {src && (
 <ViewTransition>
 <Suspense fallback={<ProfilePlaceholder />}>
 <Profile src={src} />
 </Suspense>
 </ViewTransition>
 )}
 <hr />
 <VanillaProfile />
 </>
 );
}
```

```js src/VanillaProfile.js
import { useRef } from 'react';
import { freshImageUrl } from './image.js';

export default function VanillaProfile() {
 const ref = useRef(null);
 function show() {
 ref.current.innerHTML = `<div class="card">
 <img src="${freshImageUrl()}" alt="Jack Pope" width="80" height="80" />
 <p>Jack Pope</p>
 </div>`;
 }
 return (
 <>
 <button onClick={show}>Show profile (without React)</button>
 <div ref={ref} />
 </>
 );
}
```

```js src/image.js hidden
// Add a unique parameter so the image isn't cached,
// and every run shows the loading state.
export function freshImageUrl() {
 return 'https://react.dev/images/team/jack-pope.jpg?t=' + Date.now();
}
```

```css
#root {
 min-height: 390px;
}
.card {
 margin-top: 1em;
}
.card img {
 display: block;
 border-radius: 50%;
 background: #dfe3e9;
}
.card p {
 font-weight: bold;
}
.avatar-placeholder {
 width: 80px;
 height: 80px;
 border-radius: 50%;
 background: #dfe3e9;
}
.name-placeholder {
 width: 90px;
 border-radius: 4px;
 background: #dfe3e9;
}
hr {
 margin: 16px 0;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### <CanaryBadge /> Coordinating fonts, images, and stylesheets {/*coordinating-fonts-images-and-stylesheets*/}

A Suspense boundary can wait for data, stylesheets, fonts, and images at once. Waiting for fonts and images only happens during a [`<ViewTransition>`](/reference/react/ViewTransition) update. In the example below, the `ProfileCard` component suspends while its data loads, and renders a stylesheet with `precedence`, text in a new font, and a portrait. React keeps the skeleton visible while the data and the stylesheet load. The `<ViewTransition>` reveal then waits for the font and the image, so the card appears complete.

For comparison, the version without React loads the same data and shows every resource arriving on its own schedule:

<Sandpack>

```js
import { ViewTransition, Suspense, use, useState, startTransition } from 'react';
import { fetchQuote } from './data.js';
import { freshStylesheetUrl, freshImageUrl } from './resources.js';
import VanillaProfileCard from './VanillaProfileCard.js';

function ProfileCard({ resources }) {
 const quote = use(resources.quotePromise);
 return (
 <>
 <link rel="stylesheet" href={resources.stylesheet} precedence="default" />
 <div className="profile-card">
 <img src={resources.image} alt="Jack Pope" width={80} height={80} />
 <div>
 <p className="name">Jack Pope</p>
 <p className="bio">{quote}</p>
 </div>
 </div>
 </>
 );
}

function ProfileCardPlaceholder() {
 return (
 <div className="profile-card">
 <div className="avatar-placeholder" />
 <div>
 <p className="name name-placeholder">&nbsp;</p>
 <p className="bio bio-placeholder">&nbsp;</p>
 </div>
 </div>
 );
}

export default function App() {
 const [resources, setResources] = useState(null);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setResources({
 quotePromise: fetchQuote(),
 stylesheet: freshStylesheetUrl(),
 image: freshImageUrl(),
 });
 });
 }}>
 Show profile
 </button>
 {resources && (
 <ViewTransition>
 <Suspense fallback={<ProfileCardPlaceholder />}>
 <ProfileCard resources={resources} />
 </Suspense>
 </ViewTransition>
 )}
 <hr />
 <VanillaProfileCard />
 </>
 );
}
```

```js src/VanillaProfileCard.js
import { useRef } from 'react';
import { fetchQuote } from './data.js';
import { freshStylesheetUrl, freshImageUrl } from './resources.js';

export default function VanillaProfileCard() {
 const ref = useRef(null);
 async function show() {
 const quote = await fetchQuote();
 const doc = ref.current.contentWindow.document;
 doc.open();
 doc.write(`
 <style>
 body { margin: 0; font-family: sans-serif; }
 .profile-card { display: flex; gap: 12px; align-items: center; }
 .profile-card img { border-radius: 50%; background: #dfe3e9; }
 .name { margin: 0 0 4px; font-family: 'Caveat', sans-serif; font-size: 22px; line-height: 28px; font-weight: bold; }
 .bio { margin: 0; font-family: 'Caveat', sans-serif; font-size: 20px; line-height: 26px; }
 </style>
 <div class="profile-card">
 <img src="${freshImageUrl()}" alt="Jack Pope" width="80" height="80" />
 <div>
 <p class="name">Jack Pope</p>
 <p class="bio">${quote}</p>
 </div>
 </div>
 <link rel="stylesheet" href="${freshStylesheetUrl()}">
 `);
 doc.close();
 }
 return (
 <>
 <button onClick={show}>Show profile (without React)</button>
 <iframe ref={ref} title="Vanilla profile card" className="vanilla-frame" />
 </>
 );
}
```

```js src/resources.js hidden
// Add a unique parameter so the resources aren't cached,
// and every run shows the loading state.
export function freshStylesheetUrl() {
 return (
 'https://fonts.googleapis.com/css2?family=Caveat&display=swap' +
 '&t=' +
 Date.now()
 );
}

export function freshImageUrl() {
 return 'https://react.dev/images/team/jack-pope.jpg?t=' + Date.now();
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.

export async function fetchQuote() {
 // Add a fake delay to make waiting noticeable.
 await new Promise((resolve) => {
 setTimeout(resolve, 1000);
 });
 return 'The best way to predict the future is to invent it.';
}
```

```css
#root {
 min-height: 320px;
}
button {
 margin-right: 8px;
}
hr {
 margin: 16px 0;
}
.profile-card {
 display: flex;
 gap: 12px;
 align-items: center;
 margin-top: 1em;
}
.profile-card img {
 border-radius: 50%;
 background: #dfe3e9;
}
.name {
 margin: 0 0 4px;
 font-family: 'Caveat', sans-serif;
 font-size: 22px;
 line-height: 28px;
 font-weight: bold;
}
.bio {
 margin: 0;
 font-family: 'Caveat', sans-serif;
 font-size: 20px;
 line-height: 26px;
}
.profile-card img {
 display: block;
}
.avatar-placeholder {
 width: 80px;
 height: 80px;
 border-radius: 50%;
 background: #dfe3e9;
}
.name-placeholder,
.bio-placeholder {
 border-radius: 4px;
 background: #dfe3e9;
 color: transparent;
}
.name-placeholder {
 width: 90px;
}
.bio-placeholder {
 width: 220px;
}
.vanilla-frame {
 display: block;
 margin-top: 1em;
 border: none;
 width: 100%;
 height: 110px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

## Troubleshooting {/*troubleshooting*/}

### How do I prevent the UI from being replaced by a fallback during an update? {/*preventing-unwanted-fallbacks*/}

Replacing visible UI with a fallback creates a jarring user experience. This can happen when an update causes a component to suspend, and the nearest Suspense boundary is already showing content to the user.

To prevent this from happening, [mark the update as non-urgent using `startTransition`](#preventing-already-revealed-content-from-hiding). During a Transition, React will wait until enough data has loaded to prevent an unwanted fallback from appearing:

```js {2-3,5}
function handleNextPageClick() {
 // If this update suspends, don't hide the already displayed content
 startTransition(() => {
 setCurrentPage(currentPage + 1);
 });
}
```

This will avoid hiding existing content. However, any newly rendered `Suspense` boundaries will still immediately display fallbacks to avoid blocking the UI and let the user see the content as it becomes available.

**React will only prevent unwanted fallbacks during non-urgent updates**. It will not delay a render if it's the result of an urgent update. You must opt in with an API like [`startTransition`](/reference/react/startTransition) or [`useDeferredValue`](/reference/react/useDeferredValue).

If your router is integrated with Suspense, it should wrap its updates into [`startTransition`](/reference/react/startTransition) automatically.

---
title: <ViewTransition>
version: canary
---

<Intro>

<Canary>

**The `<ViewTransition />` API is currently only available in React’s Canary and Experimental channels.**

[Learn more about React’s release channels here.](/community/versioning-policy#all-release-channels)

</Canary>

`<ViewTransition>` lets you animate a component tree with Transitions and Suspense.

```js
import {ViewTransition} from 'react';

<ViewTransition>
 <div>...</div>
</ViewTransition>
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `<ViewTransition>` {/*viewtransition*/}

Wrap a component tree in `<ViewTransition>` to animate it:

```js
<ViewTransition>
 <Page />
</ViewTransition>
```

[See more examples below.](#usage)

<DeepDive>

#### How does `<ViewTransition>` work? {/*how-does-viewtransition-work*/}

Under the hood, React applies `view-transition-name` to inline styles of the nearest DOM node nested inside the `<ViewTransition>` component. If there are multiple sibling DOM nodes like `<ViewTransition><div /><div /></ViewTransition>` then React adds a suffix to the name to make each unique but conceptually they're part of the same one. React doesn't apply these eagerly but only at the time that boundary should participate in an animation.

React automatically calls `startViewTransition` itself behind the scenes so you should never do that yourself. In fact, if you have something else on the page running a ViewTransition React will interrupt it. So it's recommended that you use React itself to coordinate these. If you had other ways to trigger ViewTransitions in the past, we recommend that you migrate to the built-in way.

If there are other React ViewTransitions already running then React will wait for them to finish before starting the next one. However, importantly if there are multiple updates happening while the first one is running, those will all be batched into one. If you start A->B. Then in the meantime you get an update to go to C and then D. When the first A->B animation finishes the next one will animate from B->D.

The `getSnapshotBeforeUpdate` lifecycle will be called before `startViewTransition` and some `view-transition-name` will update at the same time.

Then React calls `startViewTransition`. Inside the `updateCallback`, React will:

- Apply its mutations to the DOM and invoke `useInsertionEffect`.
- Wait for fonts to load.
- Call `componentDidMount`, `componentDidUpdate`, `useLayoutEffect` and refs.
- Wait for any pending Navigation to finish.
- Then React will measure any changes to the layout to see which boundaries will need to animate.

After the ready Promise of the `startViewTransition` is resolved, React will then revert the `view-transition-name`. Then React will invoke the `onEnter`, `onExit`, `onUpdate` and `onShare` callbacks to allow for manual programmatic control over the animations. This will be after the built-in default ones have already been computed.

If a `flushSync` happens to get in the middle of this sequence, then React will skip the Transition since it relies on being able to complete synchronously.

After the finished Promise of the `startViewTransition` is resolved, React will then invoke `useEffect`. This prevents those from interfering with the performance of the animation. However, this is not a guarantee because if another `setState` happens while the animation is running it'll still have to invoke the `useEffect` earlier to preserve the sequential guarantees.

</DeepDive>

#### Props {/*props*/}

- **optional** `name`: A string or object. The name of the View Transition used for shared element transitions. If not provided, React will use a unique name for each View Transition to prevent unexpected animations.
- [View Transition Class](#view-transition-class) props.
- [View Transition Event](#view-transition-event) props.

#### Caveats {/*caveats*/}

- Only use `name` for [shared element transitions](#animating-a-shared-element). For all other animations, React automatically generates a unique name to prevent unexpected animations.
- By default, `setState` updates immediately and does not activate `<ViewTransition>`, only updates wrapped in a [Transition](/reference/react/useTransition), [`<Suspense>`](/reference/react/Suspense), or `useDeferredValue` activate ViewTransition.
- `<ViewTransition>` creates an image that can be moved around, scaled and cross-faded. Unlike Layout Animations you may have seen in React Native or Motion, this means that not every individual Element inside of it animates its position. This can lead to better performance and a more continuous feeling, smooth animation compared to animating every individual piece. However, it can also lose continuity in things that should be moving by themselves. So you might have to add more `<ViewTransition>` boundaries manually as a result.
- Currently, `<ViewTransition>` only works in the DOM. We're working on adding support for React Native and other platforms.

#### Animation triggers {/*animation-triggers*/}

React automatically decides the type of View Transition animation to trigger:

- `enter`: If a `ViewTransition` is the first component inserted in this Transition, then this will activate.
- `exit`: If a `ViewTransition` is the first component deleted in this Transition, then this will activate.
- `update`: If a `ViewTransition` has any DOM mutations inside it that React is doing (such as a prop changing) or if the `ViewTransition` boundary itself changes size or position due to an immediate sibling. If there are nested `ViewTransition` then the mutation applies to them and not the parent.
- `share`: If a named `ViewTransition` is inside a deleted subtree and another named `ViewTransition` with the same name is part of an inserted subtree in the same Transition, they form a Shared Element Transition, and it animates from the deleted one to the inserted one.

By default, `<ViewTransition>` animates with a smooth cross-fade (the browser default view transition).

You can customize the animation by providing a [View Transition Class](#view-transition-class) to the `<ViewTransition>` component for each kind of trigger (see [Styling View Transitions](#styling-view-transitions)), or by using [ViewTransition Events](#view-transition-events) to control the animation with JavaScript using the [Web Animations API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Animations_API).

<Note>

#### Always check `prefers-reduced-motion` {/*always-check-prefers-reduced-motion*/}

Many users may prefer not having animations on the page. React doesn't automatically disable animations for this case.

We recommend always using the `@media (prefers-reduced-motion)` media query to disable animations or tone them down based on user preference.

In the future, CSS libraries may have this built-in to their presets.

</Note>

### View Transition Class {/*view-transition-class*/}

`<ViewTransition>` provides props to define what animations trigger:

```js
<ViewTransition
 default="none"
 enter="slide-up"
 exit="slide-down"
/>
```

#### Props {/*view-transition-class-props*/}

- **optional** `enter`: `"auto"`, `"none"`, a string, or an object.
- **optional** `exit`: `"auto"`, `"none"`, a string, or an object.
- **optional** `update`: `"auto"`, `"none"`, a string, or an object.
- **optional** `share`: `"auto"`, `"none"`, a string, or an object.
- **optional** `default`: `"auto"`, `"none"`, a string, or an object.

#### Caveats {/*view-transition-class-caveats*/}

- If `default` is `"none"` then all other triggers are turned off unless explicitly listed.

#### Values {/*view-transition-values*/}

View Transition class values can be:
- `auto`: the default. Uses the browser default animation.
- `none`: disable animations for this type.
- `<classname>`: a custom CSS class name to use for [customizing View Transitions](#styling-view-transitions).

Object values can be an object with string keys and a value of `auto`, `none` or a custom className:
- `{[type]: value}`: applies `value` if the animation matches the [Transition Type](/reference/react/addTransitionType).
- `{default: value}`: the default value to apply if no [Transition Type](/reference/react/addTransitionType) is matched.

For example, you can define a ViewTransition as:

```js
<ViewTransition
 /* turn off any animation not defined below */
 default="none"
 enter={{
 /* apply slide-in for Transition Type `forward` */
 "forward": 'slide-in',
 /* otherwise use the browser default animation */
 "default": 'auto'
 }}
 /* use the browser default for exit animations*/
 exit="auto"
 /* apply a custom `cross-fade` class for updates */
 update="cross-fade"
>
```

See [Styling View Transitions](#styling-view-transitions) for how to define CSS classes for custom animations.

---

### View Transition Event {/*view-transition-event*/}

View Transition Events allow you to control the animation with JavaScript using the [Web Animations API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Animations_API):

```js
<ViewTransition
 onEnter={instance => {/* ... */}}
 onExit={instance => {/* ... */}}
/>
```

#### Props {/*view-transition-event-props*/}

- **optional** `onEnter`: Called when an "enter" animation is triggered.
- **optional** `onExit`: Called when an "exit" animation is triggered.
- **optional** `onShare`: Called when a "share" animation is triggered.
- **optional** `onUpdate`: Called when an "update" animation is triggered.

#### Caveats {/*view-transition-event-caveats*/}
- Only one event fires per `<ViewTransition>` per Transition. `onShare` takes precedence over `onEnter` and `onExit`.
- Each event should return a **cleanup function**. The cleanup function is called when the View Transition finishes, allowing you to cancel or cleanup any animations.

#### Arguments {/*view-transition-event-arguments*/}

Each event receives two arguments:

- `instance`: A View Transition instance that provides access to the view transition [pseudo-elements](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API/Using#the_view_transition_process)
 - `old`: The `::view-transition-old` pseudo-element.
 - `new`: The `::view-transition-new` pseudo-element.
 - `name`: The `view-transition-name` string for this boundary.
 - `group`: The `::view-transition-group` pseudo-element.
 - `imagePair`: The `::view-transition-image-pair` pseudo-element.
- `types`: An `Array<string>` of [Transition Types](/reference/react/addTransitionType) included in the animation. Empty array if no types were specified.

For example, you can define a `onEnter` event that drives the animation using JavaScript:

```js
<ViewTransition
 onEnter={(instance, types) => {
 const anim = instance.new.animate([{opacity: 0}, {opacity: 1}], {
 duration: 500,
 });
 return () => anim.cancel();
 }}>
 <div>...</div>
</ViewTransition>
```

See [Animating with JavaScript](#animating-with-javascript) for more examples.

---

## Styling View Transitions {/*styling-view-transitions*/}

<Note>

In many early examples of View Transitions around the web, you'll have seen using a [`view-transition-name`](https://developer.mozilla.org/en-US/docs/Web/CSS/view-transition-name) and then style it using `::view-transition-...(my-name)` selectors. We don't recommend that for styling. Instead, we normally recommend using a View Transition Class instead.

</Note>

To customize the animation for a `<ViewTransition>` you can provide a View Transition Class to one of the activation props. The View Transition Class is a CSS class name that React applies to the child elements when the ViewTransition activates.

For example, to customize an "enter" animation, provide a class name to the `enter` prop:

```js
<ViewTransition enter="slide-in">
```

When the `<ViewTransition>` activates an "enter" animation, React will add the class name `slide-in`. Then you can refer to this class using [view transition pseudo selectors](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API#pseudo-elements) to build reusable animations:

```css
::view-transition-group(.slide-in) {
}
::view-transition-old(.slide-in) {
}
::view-transition-new(.slide-in) {
}
```

In the future, CSS libraries may add built-in animations using View Transition Classes to make this easier to use.

---

## Usage {/*usage*/}

### Animating an element on enter/exit {/*animating-an-element-on-enter*/}

Enter/Exit Transitions trigger when a `<ViewTransition>` is added or removed by a component in a transition:

```js {3}
function Child() {
 return (
 <ViewTransition enter="auto" exit="auto" default="none">
 <div>Hi</div>
 </ViewTransition>
 );
}

function Parent() {
 const [show, setShow] = useState();
 if (show) {
 return <Child />;
 }
 return null;
}
```

When `setShow` is called, `show` switches to `true` and the `Child` component is rendered. When `setShow` is called inside `startTransition`, and `Child` renders a `ViewTransition` before any other DOM nodes, an `enter` animation is triggered.

When `show` switches back to `false`, an `exit` animation is triggered.

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video} from './Video';
import videos from './data';

function Item() {
 return (
 <ViewTransition enter="auto" exit="auto" default="none">
 <Video video={videos[0]} />
 </ViewTransition>
 );
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>

 {showItem ? <Item /> : null}
 </>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

<Pitfall>

#### Only top-level ViewTransitions animate on exit/enter {/*only-top-level-viewtransition-animates-on-exit-enter*/}

`<ViewTransition>` only activates exit/enter if it is placed _before_ any DOM nodes.

If there's a `<div>` above `<ViewTransition>`, no exit/enter animations trigger:

```js [3, 5]
function Item() {
 return (
 <div> {/* 🚩<div> above <ViewTransition> breaks exit/enter */}
 <ViewTransition enter="auto" exit="auto" default="none">
 <Video video={videos[0]} />
 </ViewTransition>
 </div>
 );
}
```

This constraint prevents subtle bugs where too much or too little animates.

</Pitfall>

---

### Animating enter/exit with Activity {/*animating-enter-exit-with-activity*/}

If you want to animate a component in and out while preserving its state, or pre-rendering content for an animation, you can use [`<Activity>`](/reference/react/Activity). When a `<ViewTransition>` inside an `<Activity>` becomes visible, the `enter` animation activates. When it becomes hidden, the `exit` animation activates:

```js
<Activity mode={isVisible ? 'visible' : 'hidden'}>
 <ViewTransition enter="auto" exit="auto">
 <Counter />
 </ViewTransition>
</Activity>

```

In this example, `Counter` has a counter with internal state. Try incrementing the counter, hiding it, then showing it again. The counter's value is preserved while the sidebar animates in and out:

<Sandpack>

```js
import { Activity, ViewTransition, useState, startTransition } from 'react';

export default function App() {
 const [show, setShow] = useState(true);
 return (
 <div className="layout">
 <Toggle show={show} setShow={setShow} />
 <Activity mode={show ? 'visible' : 'hidden'}>
 <ViewTransition enter="auto" exit="auto" default="none">
 <Counter />
 </ViewTransition>
 </Activity>
 </div>
 );
}
function Toggle({show, setShow}) {
 return (
 <button
 className="toggle"
 onClick={() => {
 startTransition(() => {
 setShow(s => !s);
 });
 }}>
 {show ? 'Hide' : 'Show'}
 </button>
 )
}
function Counter() {
 const [count, setCount] = useState(0);
 return (
 <div className="counter">
 <h2>Counter</h2>
 <p>Count: {count}</p>
 <button onClick={() => setCount(count + 1)}>
 Increment
 </button>
 </div>
 );
}

```

```css
.layout {
 display: flex;
 flex-direction: column;
 align-items: flex-start;
 gap: 10px;
 min-height: 200px;
}
.counter {
 padding: 15px;
 background: #f0f4f8;
 border-radius: 8px;
 width: 200px;
}
.counter h2 {
 margin: 0 0 10px 0;
 font-size: 16px;
}
.counter p {
 margin: 0 0 10px 0;
}
.toggle {
 padding: 8px 16px;
 border: 1px solid #ccc;
 border-radius: 6px;
 background: #f0f8ff;
 cursor: pointer;
 font-size: 14px;
}
.toggle:hover {
 background: #e0e8ff;
}
.counter button {
 padding: 4px 12px;
 border: 1px solid #ccc;
 border-radius: 4px;
 background: white;
 cursor: pointer;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

Without `<Activity>`, the counter would reset to `0` every time the sidebar reappears.

---

### Animating a shared element {/*animating-a-shared-element*/}

Normally, we don't recommend assigning a name to a `<ViewTransition>` and instead let React assign it an automatic name. The reason you might want to assign a name is to animate between completely different components when one tree unmounts and another tree mounts at the same time, to preserve continuity.

```js
<ViewTransition name={UNIQUE_NAME}>
 <Child />
</ViewTransition>
```

When one tree unmounts and another mounts, if there's a pair where the same name exists in the unmounting tree and the mounting tree, they trigger the "share" animation on both. It animates from the unmounting side to the mounting side.

Unlike an exit/enter animation this can be deeply inside the deleted/mounted tree. If a `<ViewTransition>` would also be eligible for exit/enter, then the "share" animation takes precedence.

If Transition first unmounts one side and then leads to a `<Suspense>` fallback being shown before eventually the new name being mounted, then no shared element transition happens.

<Sandpack>

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video, Thumbnail, FullscreenVideo} from './Video';
import videos from './data';

export default function Component() {
 const [fullscreen, setFullscreen] = useState(false);
 if (fullscreen) {
 return (
 <FullscreenVideo
 video={videos[0]}
 onExit={() => startTransition(() => setFullscreen(false))}
 />
 );
 }
 return (
 <Video
 video={videos[0]}
 onClick={() => startTransition(() => setFullscreen(true))}
 />
 );
}
```

```js src/Video.js
import {ViewTransition} from 'react';

const THUMBNAIL_NAME = 'video-thumbnail';

export function Thumbnail({video, children}) {
 return (
 <ViewTransition name={THUMBNAIL_NAME}>
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 </ViewTransition>
 );
}

export function Video({video, onClick}) {
 return (
 <div className="video">
 <div className="link" onClick={onClick}>
 <Thumbnail video={video} />
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}

export function FullscreenVideo({video, onExit}) {
 return (
 <div className="fullscreenLayout">
 <ViewTransition name={THUMBNAIL_NAME}>
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image} fullscreen`}
 />
 <button className="close-button" onClick={onExit}>
 ✖
 </button>
 </ViewTransition>
 </div>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 height: 300px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.thumbnail.red {
 background-image: conic-gradient(at top right, #c76a15, #a6423a, #2b3491);
}
.thumbnail.fullscreen {
 width: 100%;
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
.fullscreenLayout {
 position: relative;
 height: 100%;
 width: 100%;
}
.close-button {
 position: absolute;
 top: 10px;
 right: 10px;
 color: black;
}
@keyframes progress-animation {
 from {
 width: 0;
 }
 to {
 width: 100%;
 }
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

<Note>

If either the mounted or unmounted side of a pair is outside the viewport, then no pair is formed. This ensures that it doesn't fly in or out of the viewport when something is scrolled. Instead it's treated as a regular enter/exit by itself.

This does not happen if the same Component instance changes position, which triggers an "update". Those animate regardless of whether one position is outside the viewport.

There is a known case where if a deeply nested unmounted `<ViewTransition>` is inside the viewport but the mounted side is not within the viewport, then the unmounted side animates as its own "exit" animation even if it's deeply nested instead of as part of the parent animation.

</Note>

<Pitfall>

It's important that there's only one thing with the same name mounted at a time in the entire app. Therefore it's important to use unique namespaces for the name to avoid conflicts. To ensure you can do this you might want to add a constant in a separate module that you import.

```js
export const MY_NAME = "my-globally-unique-name";
import {MY_NAME} from './shared-name';
...
<ViewTransition name={MY_NAME}>
```

</Pitfall>

---

### Animating reorder of items in a list {/*animating-reorder-of-items-in-a-list*/}

```js
items.map((item) => <Component key={item.id} item={item} />);
```

When reordering a list, without updating the content, the "update" animation triggers on each `<ViewTransition>` in the list if they're outside a DOM node. Similar to enter/exit animations.

This means that this will trigger the animation on this `<ViewTransition>`:

```js
function Component() {
 return (
 <ViewTransition>
 <div>...</div>
 </ViewTransition>
 );
}
```

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video} from './Video';
import videos from './data';

export default function Component() {
 const [orderedVideos, setOrderedVideos] = useState(videos);
 const reorder = () => {
 startTransition(() => {
 setOrderedVideos((prev) => {
 return [...prev.sort(() => Math.random() - 0.5)];
 });
 });
 };
 return (
 <>
 <button onClick={reorder}>🎲</button>
 <div className="listContainer">
 {orderedVideos.map((video, i) => {
 return (
 <ViewTransition key={video.title}>
 <Video video={video} />
 </ViewTransition>
 );
 })}
 </div>
 </>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
 {
 id: '2',
 title: 'Second video',
 description: 'Video description',
 image: 'red',
 },
 {
 id: '3',
 title: 'Third video',
 description: 'Video description',
 image: 'green',
 },
 {
 id: '4',
 title: 'Fourth video',
 description: 'Video description',
 image: 'purple',
 },
];
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 150px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.thumbnail.red {
 background-image: conic-gradient(at top right, #c76a15, #a6423a, #2b3491);
}
.thumbnail.green {
 background-image: conic-gradient(at top right, #c76a15, #388f7f, #2b3491);
}
.thumbnail.purple {
 background-image: conic-gradient(at top right, #c76a15, #575fb7, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

However, this wouldn't animate each individual item:

```js
function Component() {
 return (
 <div>
 <ViewTransition>...</ViewTransition>
 </div>
 );
}
```

Instead, any parent `<ViewTransition>` would cross-fade. If there is no parent `<ViewTransition>` then there's no animation in that case.

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video} from './Video';
import videos from './data';

export default function Component() {
 const [orderedVideos, setOrderedVideos] = useState(videos);
 const reorder = () => {
 startTransition(() => {
 setOrderedVideos((prev) => {
 return [...prev.sort(() => Math.random() - 0.5)];
 });
 });
 };
 return (
 <>
 <button onClick={reorder}>🎲</button>
 <ViewTransition>
 <div className="listContainer">
 {orderedVideos.map((video, i) => {
 return <Video video={video} key={video.title} />;
 })}
 </div>
 </ViewTransition>
 </>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
 {
 id: '2',
 title: 'Second video',
 description: 'Video description',
 image: 'red',
 },
 {
 id: '3',
 title: 'Third video',
 description: 'Video description',
 image: 'green',
 },
 {
 id: '4',
 title: 'Fourth video',
 description: 'Video description',
 image: 'purple',
 },
];
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 150px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.thumbnail.red {
 background-image: conic-gradient(at top right, #c76a15, #a6423a, #2b3491);
}
.thumbnail.green {
 background-image: conic-gradient(at top right, #c76a15, #388f7f, #2b3491);
}
.thumbnail.purple {
 background-image: conic-gradient(at top right, #c76a15, #575fb7, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

This means you might want to avoid wrapper elements in lists where you want to allow the Component to control its own reorder animation:

```
items.map(item => <div><Component key={item.id} item={item} /></div>)
```

The above rule also applies if one of the items updates to resize, which then causes the siblings to resize, it'll also animate its sibling `<ViewTransition>` but only if they're immediate siblings.

This means that during an update, which causes a lot of re-layout, it doesn't individually animate every `<ViewTransition>` on the page. That would lead to a lot of noisy animations which distracts from the actual change. Therefore React is more conservative about when an individual animation triggers.

<Pitfall>

It's important to properly use keys to preserve identity when reordering lists. It might seem like you could use "name", shared element transitions, to animate reorders but that would not trigger if one side was outside the viewport. To animate a reorder you often want to show that it went to a position outside the viewport.

</Pitfall>

---

### Animating from Suspense content {/*animating-from-suspense-content*/}

Like any Transition, React waits for data and new CSS (`<link rel="stylesheet" precedence="...">`) before running the animation. In addition to this, ViewTransitions also wait up to 500ms for new fonts to load before starting the animation to avoid them flickering in later. For the same reason, an image wrapped in ViewTransition will wait for the image to load. See examples of [waiting for a font](/reference/react/Suspense#waiting-for-a-font-to-load) and [waiting for an image](/reference/react/Suspense#waiting-for-an-image-to-load) on the Suspense page.

If it's inside a new Suspense boundary instance, then the fallback is shown first. After the Suspense boundary fully loads, it triggers the `<ViewTransition>` to animate the reveal to the content.

There are two ways to animate Suspense boundaries depending on where you place the `<ViewTransition>`:

**Update:**

```
<ViewTransition>
 <Suspense fallback={<A />}>
 <B />
 </Suspense>
</ViewTransition>
```

In this scenario when the content goes from A to B, it'll be treated as an "update" and apply that class if appropriate. Both A and B will get the same view-transition-name and therefore they're acting as a cross-fade by default.

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}

export function VideoPlaceholder() {
 const video = {image: 'loading'};
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title loading" />
 <div className="video-description loading" />
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition, Suspense} from 'react';
import {Video, VideoPlaceholder} from './Video';
import {useLazyVideoData} from './data';

function LazyVideo() {
 const video = useLazyVideoData();
 return <Video video={video} />;
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>
 {showItem ? (
 <ViewTransition>
 <Suspense fallback={<VideoPlaceholder />}>
 <LazyVideo />
 </Suspense>
 </ViewTransition>
 ) : null}
 </>
 );
}
```

```js src/data.js hidden
import {use} from 'react';

let cache = null;

function fetchVideo() {
 if (!cache) {
 cache = new Promise((resolve) => {
 setTimeout(() => {
 resolve({
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 });
 }, 1000);
 });
 }
 return cache;
}

export function useLazyVideoData() {
 return use(fetchVideo());
}
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.loading {
 background-image: linear-gradient(
 90deg,
 rgba(173, 216, 230, 0.3) 25%,
 rgba(135, 206, 250, 0.5) 50%,
 rgba(173, 216, 230, 0.3) 75%
 );
 background-size: 200% 100%;
 animation: shimmer 1.5s infinite;
}
@keyframes shimmer {
 0% {
 background-position: -200% 0;
 }
 100% {
 background-position: 200% 0;
 }
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-title.loading {
 height: 20px;
 width: 80px;
 border-radius: 0.5rem;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
 border-radius: 0.5rem;
}
.video-description.loading {
 height: 15px;
 width: 100px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

**Enter/Exit:**

```
<Suspense fallback={<ViewTransition><A /></ViewTransition>}>
 <ViewTransition><B /></ViewTransition>
</Suspense>
```

In this scenario, these are two separate ViewTransition instances each with their own `view-transition-name`. This will be treated as an "exit" of the `<A>` and an "enter" of the `<B>`.

You can achieve different effects depending on where you choose to place the `<ViewTransition>` boundary.

---

### Opting-out of an animation {/*opting-out-of-an-animation*/}

Sometimes you're wrapping a large existing component, like a whole page, and you want to animate some updates, such as changing the theme. However, you don't want it to opt-in all updates inside the whole page to cross-fade when they're updating. Especially if you're incrementally adding more animations.

You can use the class "none" to opt-out of an animation. By wrapping your children in a "none" you can disable animations for updates to them while the parent still triggers.

```js
<ViewTransition>
 <div className={theme}>
 <ViewTransition update="none">{children}</ViewTransition>
 </div>
</ViewTransition>
```

This will only animate if the theme changes and not if only the children update. The children can still opt-in again with their own `<ViewTransition>` but at least it's manual again.

---

### Customizing animations {/*customizing-animations*/}

By default, `<ViewTransition>` includes the default cross-fade from the browser.

To customize animations, you can provide props to the `<ViewTransition>` component to specify which animations to use, based on how the `<ViewTransition>` activates.

For example, we can slow down the default cross fade animation:

```js
<ViewTransition default="slow-fade">
 <Video />
</ViewTransition>
```

And define slow-fade in CSS using view transition classes:

```css
::view-transition-old(.slow-fade) {
 animation-duration: 500ms;
}

::view-transition-new(.slow-fade) {
 animation-duration: 500ms;
}
```

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>

 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video} from './Video';
import videos from './data';

function Item() {
 return (
 <ViewTransition default="slow-fade">
 <Video video={videos[0]} />
 </ViewTransition>
 );
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>

 {showItem ? <Item /> : null}
 </>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
::view-transition-old(.slow-fade) {
 animation-duration: 500ms;
}

::view-transition-new(.slow-fade) {
 animation-duration: 500ms;
}

#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

In addition to setting the `default`, you can also provide configurations for `enter`, `exit`, `update`, and `share` animations.

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>

 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video} from './Video';
import videos from './data';

function Item() {
 return (
 <ViewTransition enter="slide-in" exit="slide-out">
 <Video video={videos[0]} />
 </ViewTransition>
 );
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>

 {showItem ? <Item /> : null}
 </>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
::view-transition-old(.slide-in) {
 animation-name: slideOutRight;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-new(.slide-in) {
 animation-name: slideInRight;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-old(.slide-out) {
 animation-name: slideOutLeft;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-new(.slide-out) {
 animation-name: slideInLeft;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

@keyframes slideOutLeft {
 from {
 transform: translateX(0);
 opacity: 1;
 }
 to {
 transform: translateX(-100%);
 opacity: 0;
 }
}

@keyframes slideInLeft {
 from {
 transform: translateX(-100%);
 opacity: 0;
 }
 to {
 transform: translateX(0);
 opacity: 1;
 }
}

@keyframes slideOutRight {
 from {
 transform: translateX(0);
 opacity: 1;
 }
 to {
 transform: translateX(100%);
 opacity: 0;
 }
}

@keyframes slideInRight {
 from {
 transform: translateX(100%);
 opacity: 0;
 }
 to {
 transform: translateX(0);
 opacity: 1;
 }
}

@keyframes slideInRight {
 from {
 transform: translateX(100%);
 opacity: 0;
 }
 to {
 transform: translateX(0);
 opacity: 1;
 }
}

#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### Customizing animations with types {/*customizing-animations-with-types*/}

You can use the [`addTransitionType`](/reference/react/addTransitionType) API to add a class name to the child elements when a specific transition type is activated for a specific activation trigger. This allows you to customize the animation for each type of transition.

For example, to customize the animation for all forward and backward navigations:

```js
<ViewTransition
 default={{
 'navigation-back': 'slide-right',
 'navigation-forward': 'slide-left',
 }}>
 <div>...</div>
</ViewTransition>;

// in your router:
startTransition(() => {
 addTransitionType('navigation-' + navigationType);
});
```

When the ViewTransition activates a "navigation-back" animation, React will add the class name "slide-right". When the ViewTransition activates a "navigation-forward" animation, React will add the class name "slide-left".

In the future, routers and other libraries may add support for standard view-transition types and styles.

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>
 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {
 ViewTransition,
 addTransitionType,
 useState,
 startTransition,
} from 'react';
import {Video} from './Video';
import videos from './data';

function Item() {
 return (
 <ViewTransition
 enter={{
 'add-video-back': 'slide-in-back',
 'add-video-forward': 'slide-in-forward',
 }}
 exit={{
 'remove-video-back': 'slide-in-forward',
 'remove-video-forward': 'slide-in-back',
 }}>
 <Video video={videos[0]} />
 </ViewTransition>
 );
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <div className="button-container">
 <button
 onClick={() => {
 startTransition(() => {
 if (showItem) {
 addTransitionType('remove-video-back');
 } else {
 addTransitionType('add-video-back');
 }
 setShowItem((prev) => !prev);
 });
 }}>
 ⬅️
 </button>
 <button
 onClick={() => {
 startTransition(() => {
 if (showItem) {
 addTransitionType('remove-video-forward');
 } else {
 addTransitionType('add-video-forward');
 }
 setShowItem((prev) => !prev);
 });
 }}>
 ➡️
 </button>
 </div>
 {showItem ? <Item /> : null}
 </>
 );
}
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
::view-transition-old(.slide-in-back) {
 animation-name: slideOutRight;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-new(.slide-in-back) {
 animation-name: slideInRight;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-old(.slide-out-back) {
 animation-name: slideOutLeft;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-new(.slide-out-back) {
 animation-name: slideInLeft;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-old(.slide-in-forward) {
 animation-name: slideOutLeft;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-new(.slide-in-forward) {
 animation-name: slideInLeft;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-old(.slide-out-forward) {
 animation-name: slideOutRight;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

::view-transition-new(.slide-out-forward) {
 animation-name: slideInRight;
 animation-duration: 500ms;
 animation-timing-function: ease-in-out;
}

@keyframes slideOutLeft {
 from {
 transform: translateX(0);
 opacity: 1;
 }
 to {
 transform: translateX(-100%);
 opacity: 0;
 }
}

@keyframes slideInLeft {
 from {
 transform: translateX(-100%);
 opacity: 0;
 }
 to {
 transform: translateX(0);
 opacity: 1;
 }
}

@keyframes slideOutRight {
 from {
 transform: translateX(0);
 opacity: 1;
 }
 to {
 transform: translateX(100%);
 opacity: 0;
 }
}

@keyframes slideInRight {
 from {
 transform: translateX(100%);
 opacity: 0;
 }
 to {
 transform: translateX(0);
 opacity: 1;
 }
}

@keyframes slideInRight {
 from {
 transform: translateX(100%);
 opacity: 0;
 }
 to {
 transform: translateX(0);
 opacity: 1;
 }
}

#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.button-container {
 display: flex;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### Animating with JavaScript {/*animating-with-javascript*/}

While [View Transition Classes](#view-transition-class) let you define animations with CSS, sometimes you need imperative control over the animation. The `onEnter`, `onExit`, `onUpdate`, and `onShare` callbacks give you direct access to the view transition pseudo-elements so you can animate them using the [Web Animations API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Animations_API).

Each callback receives an `instance` with `.old` and `.new` properties representing the view transition pseudo-elements. You can call `.animate()` on them just like you would on a DOM element:

```js
<ViewTransition
 onEnter={(instance) => {
 const anim = instance.new.animate(
 [
 {transform: 'scale(0.8)'},
 {transform: 'scale(1)'},
 ],
 {duration: 300, easing: 'ease-out'}
 );
 return () => anim.cancel();
 }}>
 <div>...</div>
</ViewTransition>
```

This allows you to combine CSS-driven animations and JavaScript-driven animations.

In the following example, the default cross-fade is handled by CSS, and the slide animations are driven by JavaScript in the `onEnter` and `onExit` animations:

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>

 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition} from 'react';
import {Video} from './Video';
import videos from './data';
import {SLIDE_IN, SLIDE_OUT} from './animations';

function Item() {
 return (
 <ViewTransition
 default="none"
 /* CSS driven cross fade defaults */
 enter="auto"
 exit="auto"
 /* JS driven slide animations */
 onEnter={(instance) => {
 const anim = instance.new.animate(
 SLIDE_IN,
 {duration: 500, easing: 'ease-out'}
 );
 return () => anim.cancel();
 }}
 onExit={(instance) => {
 const anim = instance.old.animate(
 SLIDE_OUT,
 {duration: 300, easing: 'ease-in'}
 );
 return () => anim.cancel();
 }}>
 <Video video={videos[0]} />
 </ViewTransition>
 );
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 return (
 <>
 <button
 onClick={() => {
 startTransition(() => {
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>

 {showItem ? <Item /> : null}
 </>
 );
}
```

```js src/animations.js
export const SLIDE_IN = [
 {transform: 'translateY(20px)'},
 {transform: 'translateY(0)'},
];

export const SLIDE_OUT = [
 {transform: 'translateY(0)'},
 {transform: 'translateY(-20px)'},
];
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}

```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

<Note>

#### Always clean up View Transition Events {/*always-clean-up-view-transition-events*/}

View Transition Events should always return a cleanup function:

```js {7}
<ViewTransition
 onEnter={(instance) => {
 const anim = instance.new.animate(
 SLIDE_IN,
 {duration: 500, easing: 'ease-out'}
 );
 return () => anim.cancel();
 }}
>
```

This allows the browser to cancel the animation when the View Transition is interrupted.

</Note>

---

### Animating transition types with JavaScript {/*animating-transition-types-with-javascript*/}

You can use `types` passed to `ViewTransition` events to conditionally apply different animations based on how the Transition was triggered.

```js {3}
 <ViewTransition
 onEnter={(instance, types) => {
 const duration = types.includes('fast') ? 150 : 2000;
 const anim = instance.new.animate(
 SLIDE_IN,
 {duration: duration, easing: 'ease-out'}
 );
 return () => anim.cancel();
 }}
>
```

This example calls [`addTransitionType`](/reference/react/addTransitionType) to mark a Transition as "fast" and then adjust the animation duration:

<Sandpack>

```js src/Video.js hidden
function Thumbnail({video, children}) {
 return (
 <div
 aria-hidden="true"
 tabIndex={-1}
 className={`thumbnail ${video.image}`}
 />
 );
}

export function Video({video}) {
 return (
 <div className="video">
 <div className="link">
 <Thumbnail video={video}></Thumbnail>

 <div className="info">
 <div className="video-title">{video.title}</div>
 <div className="video-description">{video.description}</div>
 </div>
 </div>
 </div>
 );
}
```

```js
import {ViewTransition, useState, startTransition, addTransitionType} from 'react';
import {Video} from './Video';
import videos from './data';
import {SLIDE_IN, SLIDE_OUT} from './animations';

function Item() {
 return (
 <ViewTransition
 onEnter={(instance, types) => {
 const duration = types.includes('fast') ? 150 : 2000;
 const anim = instance.new.animate(
 SLIDE_IN,
 {duration: duration, easing: 'ease-out'}
 );
 return () => anim.cancel();
 }}
 onExit={(instance, types) => {
 const duration = types.includes('fast') ? 150 : 500;
 const anim = instance.old.animate(
 SLIDE_OUT,
 {duration: duration, easing: 'ease-in'}
 );
 return () => anim.cancel();
 }}>
 <Video video={videos[0]} />
 </ViewTransition>
 );
}

export default function Component() {
 const [showItem, setShowItem] = useState(false);
 const [isFast, setIsFast] = useState(false);
 return (
 <>
 <div>
 Fast: <input type="checkbox" onChange={() => {setIsFast(f => !f)}} value={isFast}></input>
 </div><br />
 <button
 onClick={() => {
 startTransition(() => {
 if (isFast) {
 addTransitionType('fast');
 }
 setShowItem((prev) => !prev);
 });
 }}>
 {showItem ? '➖' : '➕'}
 </button>

 {showItem ? <Item /> : null}
 </>
 );
}
```

```js src/animations.js
export const SLIDE_IN = [
 {opacity: 0, transform: 'translateY(20px)'},
 {opacity: 1, transform: 'translateY(0)'},
];

export const SLIDE_OUT = [
 {opacity: 1, transform: 'translateY(0)'},
 {opacity: 0, transform: 'translateY(-20px)'},
];
```

```js src/data.js hidden
export default [
 {
 id: '1',
 title: 'First video',
 description: 'Video description',
 image: 'blue',
 },
];
```

```css
#root {
 display: flex;
 flex-direction: column;
 align-items: center;
 min-height: 200px;
}
button {
 border: none;
 border-radius: 50%;
 width: 50px;
 height: 50px;
 display: flex;
 justify-content: center;
 align-items: center;
 background-color: #f0f8ff;
 color: white;
 font-size: 20px;
 cursor: pointer;
 transition: background-color 0.3s, border 0.3s;
}
button:hover {
 border: 2px solid #ccc;
 background-color: #e0e8ff;
}
.thumbnail {
 position: relative;
 aspect-ratio: 16 / 9;
 display: flex;
 overflow: hidden;
 flex-direction: column;
 justify-content: center;
 align-items: center;
 border-radius: 0.5rem;
 outline-offset: 2px;
 width: 8rem;
 vertical-align: middle;
 background-color: #ffffff;
 background-size: cover;
 user-select: none;
}
.thumbnail.blue {
 background-image: conic-gradient(at top right, #c76a15, #087ea4, #2b3491);
}
.video {
 display: flex;
 flex-direction: row;
 gap: 0.75rem;
 align-items: center;
 margin-top: 1em;
}
.video .link {
 display: flex;
 flex-direction: row;
 flex: 1 1 0;
 gap: 0.125rem;
 outline-offset: 4px;
 cursor: pointer;
}
.video .info {
 display: flex;
 flex-direction: column;
 justify-content: center;
 margin-left: 8px;
 gap: 0.125rem;
}
.video .info:hover {
 text-decoration: underline;
}
.video-title {
 font-size: 15px;
 line-height: 1.25;
 font-weight: 700;
 color: #23272f;
}
.video-description {
 color: #5e687e;
 font-size: 13px;
}

```

```json package.json hidden
{
 "dependencies": {
 "react": "canary",
 "react-dom": "canary",
 "react-scripts": "latest"
 }
}
```

</Sandpack>

---

### Building View Transition enabled routers {/*building-view-transition-enabled-routers*/}

React waits for any pending Navigation to finish to ensure that scroll restoration happens within the animation. If the Navigation is blocked on React, your router must unblock in `useLayoutEffect` since `useEffect` would lead to a deadlock.

If a `startTransition` is started from the legacy popstate event, such as during a "back"-navigation then it must finish synchronously to ensure scroll and form restoration works correctly. This is in conflict with running a View Transition animation. Therefore, React will skip animations from popstate and animations won't run for the back button. You can fix this by upgrading your router to use the Navigation API.

---

## Troubleshooting {/*troubleshooting*/}

### My `<ViewTransition>` is not activating {/*my-viewtransition-is-not-activating*/}

`<ViewTransition>` only activates if it is placed before any DOM node:

```js [3, 5]
function Component() {
 return (
 <div>
 <ViewTransition>Hi</ViewTransition>
 </div>
 );
}
```

To fix, ensure that the `<ViewTransition>` comes before any other DOM nodes:

```js [3, 5]
function Component() {
 return (
 <ViewTransition>
 <div>Hi</div>
 </ViewTransition>
 );
}
```

### I'm getting an error "There are two `<ViewTransition name=%s>` components with the same name mounted at the same time." {/*two-viewtransition-with-same-name*/}

This error occurs when two `<ViewTransition>` components with the same `name` are mounted at the same time:

```js [3]
function Item() {
 // 🚩 All items will get the same "name".
 return <ViewTransition name="item">...</ViewTransition>;
}

function ItemList({items}) {
 return (
 <>
 {items.map((item) => (
 <Item key={item.id} />
 ))}
 </>
 );
}
```

This will cause the View Transition to error. In development, React detects this issue to surface it and logs two errors:

<ConsoleBlockMulti>
<ConsoleLogLine level="error">

There are two `<ViewTransition name=%s>` components with the same name mounted at the same time. This is not supported and will cause View Transitions to error. Try to use a more unique name e.g. by using a namespace prefix and adding the id of an item to the name.
{' '}at Item
{' '}at ItemList

</ConsoleLogLine>

<ConsoleLogLine level="error">

The existing `<ViewTransition name=%s>` duplicate has this stack trace.
{' '}at Item
{' '}at ItemList

</ConsoleLogLine>
</ConsoleBlockMulti>

To fix, ensure that there's only one `<ViewTransition>` with the same name mounted at a time in the entire app by ensuring the `name` is unique, or adding an `id` to the name:

```js [3]
function Item({id}) {
 // ✅ All items will get a unique name.
 return <ViewTransition name={`item-${id}`}>...</ViewTransition>;
}

function ItemList({items}) {
 return (
 <>
 {items.map((item) => (
 <Item key={item.id} item={item} />
 ))}
 </>
 );
}
```

---
title: act
---

<Intro>

`act` is a test helper to apply pending React updates before making assertions.

```js
await act(async actFn)
```

</Intro>

To prepare a component for assertions, wrap the code rendering it and performing updates inside an `await act()` call. This makes your test run closer to how React works in the browser.

<Note>
You might find using `act()` directly a bit too verbose. To avoid some of the boilerplate, you could use a library like [React Testing Library](https://testing-library.com/docs/react-testing-library/intro), whose helpers are wrapped with `act()`.
</Note>

<InlineToc />

---

## Reference {/*reference*/}

### `await act(async actFn)` {/*await-act-async-actfn*/}

When writing UI tests, tasks like rendering, user events, or data fetching can be considered as “units” of interaction with a user interface. React provides a helper called `act()` that makes sure all updates related to these “units” have been processed and applied to the DOM before you make any assertions.

The name `act` comes from the [Arrange-Act-Assert](https://wiki.c2.com/?ArrangeActAssert) pattern.

```js {2,4}
it ('renders with button disabled', async () => {
 await act(async () => {
 root.render(<TestComponent />)
 });
 expect(container.querySelector('button')).toBeDisabled();
});
```

<Note>

We recommend using `act` with `await` and an `async` function. Although the sync version works in many cases, it doesn't work in all cases and due to the way React schedules updates internally, it's difficult to predict when you can use the sync version.

We will deprecate and remove the sync version in the future.

</Note>

#### Parameters {/*parameters*/}

* `async actFn`: An async function wrapping renders or interactions for components being tested. Any updates triggered within the `actFn`, are added to an internal act queue, which are then flushed together to process and apply any changes to the DOM. Since it is async, React will also run any code that crosses an async boundary, and flush any updates scheduled.

#### Returns {/*returns*/}

`act` does not return anything.

## Usage {/*usage*/}

When testing a component, you can use `act` to make assertions about its output.

For example, let’s say we have this `Counter` component, the usage examples below show how to test it:

```js
function Counter() {
 const [count, setCount] = useState(0);
 const handleClick = () => {
 setCount(prev => prev + 1);
 }

 useEffect(() => {
 document.title = `You clicked ${count} times`;
 }, [count]);

 return (
 <div>
 <p>You clicked {count} times</p>
 <button onClick={handleClick}>
 Click me
 </button>
 </div>
 )
}
```

### Rendering components in tests {/*rendering-components-in-tests*/}

To test the render output of a component, wrap the render inside `act()`:

```js {10,12}
import {act} from 'react';
import ReactDOMClient from 'react-dom/client';
import Counter from './Counter';

it('can render and update a counter', async () => {
 container = document.createElement('div');
 document.body.appendChild(container);

 // ✅ Render the component inside act().
 await act(() => {
 ReactDOMClient.createRoot(container).render(<Counter />);
 });

 const button = container.querySelector('button');
 const label = container.querySelector('p');
 expect(label.textContent).toBe('You clicked 0 times');
 expect(document.title).toBe('You clicked 0 times');
});
```

Here, we create a container, append it to the document, and render the `Counter` component inside `act()`. This ensures that the component is rendered and its effects are applied before making assertions.

Using `act` ensures that all updates have been applied before we make assertions.

### Dispatching events in tests {/*dispatching-events-in-tests*/}

To test events, wrap the event dispatch inside `act()`:

```js {14,16}
import {act} from 'react';
import ReactDOMClient from 'react-dom/client';
import Counter from './Counter';

it.only('can render and update a counter', async () => {
 const container = document.createElement('div');
 document.body.appendChild(container);

 await act( async () => {
 ReactDOMClient.createRoot(container).render(<Counter />);
 });

 // ✅ Dispatch the event inside act().
 await act(async () => {
 button.dispatchEvent(new MouseEvent('click', { bubbles: true }));
 });

 const button = container.querySelector('button');
 const label = container.querySelector('p');
 expect(label.textContent).toBe('You clicked 1 times');
 expect(document.title).toBe('You clicked 1 times');
});
```

Here, we render the component with `act`, and then dispatch the event inside another `act()`. This ensures that all updates from the event are applied before making assertions.

<Pitfall>

Don’t forget that dispatching DOM events only works when the DOM container is added to the document. You can use a library like [React Testing Library](https://testing-library.com/docs/react-testing-library/intro) to reduce the boilerplate code.

</Pitfall>

## Troubleshooting {/*troubleshooting*/}

### I'm getting an error: "The current testing environment is not configured to support act(...)" {/*error-the-current-testing-environment-is-not-configured-to-support-act*/}

Using `act` requires setting `global.IS_REACT_ACT_ENVIRONMENT=true` in your test environment. This is to ensure that `act` is only used in the correct environment.

If you don't set the global, you will see an error like this:

<ConsoleBlock level="error">

Warning: The current testing environment is not configured to support act(...)

</ConsoleBlock>

To fix, add this to your global setup file for React tests:

```js
global.IS_REACT_ACT_ENVIRONMENT=true
```

<Note>

In testing frameworks like [React Testing Library](https://testing-library.com/docs/react-testing-library/intro), `IS_REACT_ACT_ENVIRONMENT` is already set for you.

</Note>

---
title: addTransitionType
version: canary
---

<Canary>

**The `addTransitionType` API is currently only available in React’s Canary and Experimental channels.**

[Learn more about React’s release channels here.](/community/versioning-policy#all-release-channels)

</Canary>

<Intro>

`addTransitionType` lets you specify the cause of a transition.

```js
startTransition(() => {
 addTransitionType('my-transition-type');
 setState(newState);
});
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `addTransitionType` {/*addtransitiontype*/}

#### Parameters {/*parameters*/}

- `type`: The type of transition to add. This can be any string.

#### Returns {/*returns*/}

`addTransitionType` does not return anything.

#### Caveats {/*caveats*/}

- If multiple transitions are combined, all Transition Types are collected. You can also add more than one type to a Transition.
- Transition Types are reset after each commit. This means a `<Suspense>` fallback will associate the types after a `startTransition`, but revealing the content does not.

---

## Usage {/*usage*/}

### Adding the cause of a transition {/*adding-the-cause-of-a-transition*/}

Call `addTransitionType` inside of `startTransition` to indicate the cause of a transition:

``` [[1, 6, "addTransitionType"], [2, 5, "startTransition", [3, 6, "'submit-click'"]]
import { startTransition, addTransitionType } from 'react';

function Submit({action) {
 function handleClick() {
 startTransition(() => {
 addTransitionType('submit-click');
 action();
 });
 }

 return <button onClick={handleClick}>Click me</button>;
}

```

When you call <CodeStep step={1}>addTransitionType</CodeStep> inside the scope of <CodeStep step={2}>startTransition</CodeStep>, React will associate <CodeStep step={3}>submit-click</CodeStep> as one of the causes for the Transition.

Currently, Transition Types can be used to customize different animations based on what caused the Transition. You have three different ways to choose from for how to use them:

- [Customize animations using browser view transition types](#customize-animations-using-browser-view-transition-types)
- [Customize animations using `View Transition` Class](#customize-animations-using-view-transition-class)
- [Customize animations using `ViewTransition` events](#customize-animations-using-viewtransition-events)

In the future, we plan to support more use cases for using the cause of a transition.

---
### Customize animations using browser view transition types {/*customize-animations-using-browser-view-transition-types*/}

When a [`ViewTransition`](/reference/react/ViewTransition) activates from a transition, React adds all the Transition Types as browser [view transition types](https://www.w3.org/TR/css-view-transitions-2/#active-view-transition-pseudo-examples) to the element.

This allows you to customize different animations based on CSS scopes:

```js [11]
function Component() {
 return (
 <ViewTransition>
 <div>Hello</div>
 </ViewTransition>
 );
}

startTransition(() => {
 addTransitionType('my-transition-type');
 setShow(true);
});
```

```css
:root:active-view-transition-type(my-transition-type) {
 &::view-transition-...(...) {
 ...
 }
}
```

---

### Customize animations using `View Transition` Class {/*customize-animations-using-view-transition-class*/}

You can customize animations for an activated `ViewTransition` based on type by passing an object to the View Transition Class:

```js
function Component() {
 return (
 <ViewTransition enter={{
 'my-transition-type': 'my-transition-class',
 }}>
 <div>Hello</div>
 </ViewTransition>
 );
}

// ...
startTransition(() => {
 addTransitionType('my-transition-type');
 setState(newState);
});
```

If multiple types match, then they're joined together. If no types match then the special "default" entry is used instead. If any type has the value "none" then that wins and the ViewTransition is disabled (not assigned a name).

These can be combined with enter/exit/update/layout/share props to match based on kind of trigger and Transition Type.

```js
<ViewTransition enter={{
 'navigation-back': 'enter-right',
 'navigation-forward': 'enter-left',
}}
exit={{
 'navigation-back': 'exit-right',
 'navigation-forward': 'exit-left',
}}>
```

---

### Customize animations using `ViewTransition` events {/*customize-animations-using-viewtransition-events*/}

You can imperatively customize animations for an activated `ViewTransition` based on type using View Transition events:

```
<ViewTransition onUpdate={(inst, types) => {
 if (types.includes('navigation-back')) {
 ...
 } else if (types.includes('navigation-forward')) {
 ...
 } else {
 ...
 }
}}>
```

This allows you to pick different imperative Animations based on the cause.

---
title: "Built-in React APIs"
---

<Intro>

In addition to [Hooks](/reference/react/hooks) and [Components](/reference/react/components), the `react` package exports a few other APIs that are useful for defining components. This page lists all the remaining modern React APIs.

</Intro>

---

* [`createContext`](/reference/react/createContext) lets you define and provide context to the child components. Used with [`useContext`.](/reference/react/useContext)
* [`lazy`](/reference/react/lazy) lets you defer loading a component's code until it's rendered for the first time.
* [`memo`](/reference/react/memo) lets your component skip re-renders with same props. Used with [`useMemo`](/reference/react/useMemo) and [`useCallback`.](/reference/react/useCallback)
* [`startTransition`](/reference/react/startTransition) lets you mark a state update as non-urgent. Similar to [`useTransition`.](/reference/react/useTransition)
* [`act`](/reference/react/act) lets you wrap renders and interactions in tests to ensure updates have processed before making assertions.

---

## Resource APIs {/*resource-apis*/}

*Resources* can be accessed by a component without having them as part of their state. For example, a component can read a message from a Promise or read styling information from a context.

To read a value from a resource, use this API:

* [`use`](/reference/react/use) lets you read the value of a resource like a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) or [context](/learn/passing-data-deeply-with-context).
```js
function MessageComponent({ messagePromise }) {
 const message = use(messagePromise);
 const theme = use(ThemeContext);
 // ...
}
```

---
title: cache
---

<RSC>

`cache` is only for use with [React Server Components](/reference/rsc/server-components).

</RSC>

<Intro>

`cache` lets you cache the result of a data fetch or computation.

```js
const cachedFn = cache(fn);
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `cache(fn)` {/*cache*/}

Call `cache` outside of any components to create a version of the function with caching.

```js {4,7}
import {cache} from 'react';
import calculateMetrics from 'lib/metrics';

const getMetrics = cache(calculateMetrics);

function Chart({data}) {
 const report = getMetrics(data);
 // ...
}
```

When `getMetrics` is first called with `data`, `getMetrics` will call `calculateMetrics(data)` and store the result in cache. If `getMetrics` is called again with the same `data`, it will return the cached result instead of calling `calculateMetrics(data)` again.

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

- `fn`: The function you want to cache results for. `fn` can take any arguments and return any value.

#### Returns {/*returns*/}

`cache` returns a cached version of `fn` with the same type signature. It does not call `fn` in the process.

When calling `cachedFn` with given arguments, it first checks if a cached result exists in the cache. If a cached result exists, it returns the result. If not, it calls `fn` with the arguments, stores the result in the cache, and returns the result. The only time `fn` is called is when there is a cache miss.

<Note>

The optimization of caching return values based on inputs is known as [_memoization_](https://en.wikipedia.org/wiki/Memoization). We refer to the function returned from `cache` as a memoized function.

</Note>

#### Caveats {/*caveats*/}

- React will invalidate the cache for all memoized functions for each server request.
- Each call to `cache` creates a new function. This means that calling `cache` with the same function multiple times will return different memoized functions that do not share the same cache.
- `cachedFn` will also cache errors. If `fn` throws an error for certain arguments, it will be cached, and the same error is re-thrown when `cachedFn` is called with those same arguments.
- `cache` is for use in [Server Components](/reference/rsc/server-components) only.

---

## Usage {/*usage*/}

### Cache an expensive computation {/*cache-expensive-computation*/}

Use `cache` to skip duplicate work.

```js [[1, 7, "getUserMetrics(user)"],[2, 13, "getUserMetrics(user)"]]
import {cache} from 'react';
import calculateUserMetrics from 'lib/user';

const getUserMetrics = cache(calculateUserMetrics);

function Profile({user}) {
 const metrics = getUserMetrics(user);
 // ...
}

function TeamReport({users}) {
 for (let user in users) {
 const metrics = getUserMetrics(user);
 // ...
 }
 // ...
}
```

If the same `user` object is rendered in both `Profile` and `TeamReport`, the two components can share work and only call `calculateUserMetrics` once for that `user`.

Assume `Profile` is rendered first. It will call <CodeStep step={1}>`getUserMetrics`</CodeStep>, and check if there is a cached result. Since it is the first time `getUserMetrics` is called with that `user`, there will be a cache miss. `getUserMetrics` will then call `calculateUserMetrics` with that `user` and write the result to cache.

When `TeamReport` renders its list of `users` and reaches the same `user` object, it will call <CodeStep step={2}>`getUserMetrics`</CodeStep> and read the result from cache.

If `calculateUserMetrics` can be aborted by passing an [`AbortSignal`](https://developer.mozilla.org/en-US/docs/Web/API/AbortSignal), you can use [`cacheSignal()`](/reference/react/cacheSignal) to cancel the expensive computation if React has finished rendering. `calculateUserMetrics` may already handle cancellation internally by using `cacheSignal` directly.

<Pitfall>

##### Calling different memoized functions will read from different caches. {/*pitfall-different-memoized-functions*/}

To access the same cache, components must call the same memoized function.

```js [[1, 7, "getWeekReport"], [1, 7, "cache(calculateWeekReport)"], [1, 8, "getWeekReport"]]
// Temperature.js
import {cache} from 'react';
import {calculateWeekReport} from './report';

export function Temperature({cityData}) {
 // 🚩 Wrong: Calling `cache` in component creates new `getWeekReport` for each render
 const getWeekReport = cache(calculateWeekReport);
 const report = getWeekReport(cityData);
 // ...
}
```

```js [[2, 6, "getWeekReport"], [2, 6, "cache(calculateWeekReport)"], [2, 9, "getWeekReport"]]
// Precipitation.js
import {cache} from 'react';
import {calculateWeekReport} from './report';

// 🚩 Wrong: `getWeekReport` is only accessible for `Precipitation` component.
const getWeekReport = cache(calculateWeekReport);

export function Precipitation({cityData}) {
 const report = getWeekReport(cityData);
 // ...
}
```

In the above example, <CodeStep step={2}>`Precipitation`</CodeStep> and <CodeStep step={1}>`Temperature`</CodeStep> each call `cache` to create a new memoized function with their own cache look-up. If both components render for the same `cityData`, they will do duplicate work to call `calculateWeekReport`.

In addition, `Temperature` creates a <CodeStep step={1}>new memoized function</CodeStep> each time the component is rendered which doesn't allow for any cache sharing.

To maximize cache hits and reduce work, the two components should call the same memoized function to access the same cache. Instead, define the memoized function in a dedicated module that can be [`import`-ed](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import) across components.

```js [[3, 5, "export default cache(calculateWeekReport)"]]
// getWeekReport.js
import {cache} from 'react';
import {calculateWeekReport} from './report';

export default cache(calculateWeekReport);
```

```js [[3, 2, "getWeekReport", 0], [3, 5, "getWeekReport"]]
// Temperature.js
import getWeekReport from './getWeekReport';

export default function Temperature({cityData}) {
	const report = getWeekReport(cityData);
 // ...
}
```

```js [[3, 2, "getWeekReport", 0], [3, 5, "getWeekReport"]]
// Precipitation.js
import getWeekReport from './getWeekReport';

export default function Precipitation({cityData}) {
 const report = getWeekReport(cityData);
 // ...
}
```
Here, both components call the <CodeStep step={3}>same memoized function</CodeStep> exported from `./getWeekReport.js` to read and write to the same cache.
</Pitfall>

### Share a snapshot of data {/*take-and-share-snapshot-of-data*/}

To share a snapshot of data between components, call `cache` with a data-fetching function like `fetch`. When multiple components make the same data fetch, only one request is made and the data returned is cached and shared across components. All components refer to the same snapshot of data across the server render.

```js [[1, 4, "city"], [1, 5, "fetchTemperature(city)"], [2, 4, "getTemperature"], [2, 9, "getTemperature"], [1, 9, "city"], [2, 14, "getTemperature"], [1, 14, "city"]]
import {cache} from 'react';
import {fetchTemperature} from './api.js';

const getTemperature = cache(async (city) => {
	return await fetchTemperature(city);
});

async function AnimatedWeatherCard({city}) {
	const temperature = await getTemperature(city);
	// ...
}

async function MinimalWeatherCard({city}) {
	const temperature = await getTemperature(city);
	// ...
}
```

If `AnimatedWeatherCard` and `MinimalWeatherCard` both render for the same <CodeStep step={1}>city</CodeStep>, they will receive the same snapshot of data from the <CodeStep step={2}>memoized function</CodeStep>.

If `AnimatedWeatherCard` and `MinimalWeatherCard` supply different <CodeStep step={1}>city</CodeStep> arguments to <CodeStep step={2}>`getTemperature`</CodeStep>, then `fetchTemperature` will be called twice and each call site will receive different data.

The <CodeStep step={1}>city</CodeStep> acts as a cache key.

<Note>

<CodeStep step={3}>Asynchronous rendering</CodeStep> is only supported for Server Components.

```js [[3, 1, "async"], [3, 2, "await"]]
async function AnimatedWeatherCard({city}) {
	const temperature = await getTemperature(city);
	// ...
}
```

To render components that use asynchronous data in Client Components, see [`use()` documentation](/reference/react/use).

</Note>

### Preload data {/*preload-data*/}

By caching a long-running data fetch, you can kick off asynchronous work prior to rendering the component.

```jsx [[2, 6, "await getUser(id)"], [1, 17, "getUser(id)"]]
const getUser = cache(async (id) => {
 return await db.user.query(id);
});

async function Profile({id}) {
 const user = await getUser(id);
 return (
 <section>
 <img src={user.profilePic} />
 <h2>{user.name}</h2>
 </section>
 );
}

function Page({id}) {
 // ✅ Good: start fetching the user data
 getUser(id);
 // ... some computational work
 return (
 <>
 <Profile id={id} />
 </>
 );
}
```

When rendering `Page`, the component calls <CodeStep step={1}>`getUser`</CodeStep> but note that it doesn't use the returned data. This early <CodeStep step={1}>`getUser`</CodeStep> call kicks off the asynchronous database query that occurs while `Page` is doing other computational work and rendering children.

When rendering `Profile`, we call <CodeStep step={2}>`getUser`</CodeStep> again. If the initial <CodeStep step={1}>`getUser`</CodeStep> call has already returned and cached the user data, when `Profile` <CodeStep step={2}>asks and waits for this data</CodeStep>, it can simply read from the cache without requiring another remote procedure call. If the <CodeStep step={1}> initial data request</CodeStep> hasn't been completed, preloading data in this pattern reduces delay in data-fetching.

<DeepDive>

#### Caching asynchronous work {/*caching-asynchronous-work*/}

When evaluating an [asynchronous function](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function), you will receive a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) for that work. The promise holds the state of that work (_pending_, _fulfilled_, _failed_) and its eventual settled result.

In this example, the asynchronous function <CodeStep step={1}>`fetchData`</CodeStep> returns a promise that is awaiting the `fetch`.

```js [[1, 1, "fetchData()"], [2, 8, "getData()"], [3, 10, "getData()"]]
async function fetchData() {
 return await fetch(`https://...`);
}

const getData = cache(fetchData);

async function MyComponent() {
 getData();
 // ... some computational work
 await getData();
 // ...
}
```

In calling <CodeStep step={2}>`getData`</CodeStep> the first time, the promise returned from <CodeStep step={1}>`fetchData`</CodeStep> is cached. Subsequent look-ups will then return the same promise.

Notice that the first <CodeStep step={2}>`getData`</CodeStep> call does not `await` whereas the <CodeStep step={3}>second</CodeStep> does. [`await`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/await) is a JavaScript operator that will wait and return the settled result of the promise. The first <CodeStep step={2}>`getData`</CodeStep> call simply initiates the `fetch` to cache the promise for the second <CodeStep step={3}>`getData`</CodeStep> to look-up.

If by the <CodeStep step={3}>second call</CodeStep> the promise is still _pending_, then `await` will pause for the result. The optimization is that while we wait on the `fetch`, React can continue with computational work, thus reducing the wait time for the <CodeStep step={3}>second call</CodeStep>.

If the promise is already settled, either to an error or the _fulfilled_ result, `await` will return that value immediately. In both outcomes, there is a performance benefit.
</DeepDive>

<Pitfall>

##### Calling a memoized function outside of a component will not use the cache. {/*pitfall-memoized-call-outside-component*/}

```jsx [[1, 3, "getUser"]]
import {cache} from 'react';

const getUser = cache(async (userId) => {
 return await db.user.query(userId);
});

// 🚩 Wrong: Calling memoized function outside of component will not memoize.
getUser('demo-id');

async function DemoProfile() {
 // ✅ Good: `getUser` will memoize.
 const user = await getUser('demo-id');
 return <Profile user={user} />;
}
```

React only provides cache access to the memoized function in a component. When calling <CodeStep step={1}>`getUser`</CodeStep> outside of a component, it will still evaluate the function but not read or update the cache.

This is because cache access is provided through a [context](/learn/passing-data-deeply-with-context) which is only accessible from a component.

</Pitfall>

<DeepDive>

#### When should I use `cache`, [`memo`](/reference/react/memo), or [`useMemo`](/reference/react/useMemo)? {/*cache-memo-usememo*/}

All mentioned APIs offer memoization but the difference is what they're intended to memoize, who can access the cache, and when their cache is invalidated.

#### `useMemo` {/*deep-dive-use-memo*/}

In general, you should use [`useMemo`](/reference/react/useMemo) for caching an expensive computation in a Client Component across renders. As an example, to memoize a transformation of data within a component.

```jsx {expectedErrors: {'react-compiler': [4]}} {4}
'use client';

function WeatherReport({record}) {
 const avgTemp = useMemo(() => calculateAvg(record), record);
 // ...
}

function App() {
 const record = getRecord();
 return (
 <>
 <WeatherReport record={record} />
 <WeatherReport record={record} />
 </>
 );
}
```
In this example, `App` renders two `WeatherReport`s with the same record. Even though both components do the same work, they cannot share work. `useMemo`'s cache is only local to the component.

However, `useMemo` does ensure that if `App` re-renders and the `record` object doesn't change, each component instance would skip work and use the memoized value of `avgTemp`. `useMemo` will only cache the last computation of `avgTemp` with the given dependencies.

#### `cache` {/*deep-dive-cache*/}

In general, you should use `cache` in Server Components to memoize work that can be shared across components.

```js [[1, 12, "<WeatherReport city={city} />"], [3, 13, "<WeatherReport city={city} />"], [2, 1, "cache(fetchReport)"]]
const cachedFetchReport = cache(fetchReport);

function WeatherReport({city}) {
 const report = cachedFetchReport(city);
 // ...
}

function App() {
 const city = "Los Angeles";
 return (
 <>
 <WeatherReport city={city} />
 <WeatherReport city={city} />
 </>
 );
}
```
Re-writing the previous example to use `cache`, in this case the <CodeStep step={3}>second instance of `WeatherReport`</CodeStep> will be able to skip duplicate work and read from the same cache as the <CodeStep step={1}>first `WeatherReport`</CodeStep>. Another difference from the previous example is that `cache` is also recommended for <CodeStep step={2}>memoizing data fetches</CodeStep>, unlike `useMemo` which should only be used for computations.

At this time, `cache` should only be used in Server Components and the cache will be invalidated across server requests.

#### `memo` {/*deep-dive-memo*/}

You should use [`memo`](reference/react/memo) to prevent a component re-rendering if its props are unchanged.

```js
'use client';

function WeatherReport({record}) {
 const avgTemp = calculateAvg(record);
 // ...
}

const MemoWeatherReport = memo(WeatherReport);

function App() {
 const record = getRecord();
 return (
 <>
 <MemoWeatherReport record={record} />
 <MemoWeatherReport record={record} />
 </>
 );
}
```

In this example, both `MemoWeatherReport` components will call `calculateAvg` when first rendered. However, if `App` re-renders, with no changes to `record`, none of the props have changed and `MemoWeatherReport` will not re-render.

Compared to `useMemo`, `memo` memoizes the component render based on props vs. specific computations. Similar to `useMemo`, the memoized component only caches the last render with the last prop values. Once the props change, the cache invalidates and the component re-renders.

</DeepDive>

---

## Troubleshooting {/*troubleshooting*/}

### My memoized function still runs even though I've called it with the same arguments {/*memoized-function-still-runs*/}

See prior mentioned pitfalls
* [Calling different memoized functions will read from different caches.](#pitfall-different-memoized-functions)
* [Calling a memoized function outside of a component will not use the cache.](#pitfall-memoized-call-outside-component)

If none of the above apply, it may be a problem with how React checks if something exists in cache.

If your arguments are not [primitives](https://developer.mozilla.org/en-US/docs/Glossary/Primitive) (ex. objects, functions, arrays), ensure you're passing the same object reference.

When calling a memoized function, React will look up the input arguments to see if a result is already cached. React will use shallow equality of the arguments to determine if there is a cache hit.

```js
import {cache} from 'react';

const calculateNorm = cache((vector) => {
 // ...
});

function MapMarker(props) {
 // 🚩 Wrong: props is an object that changes every render.
 const length = calculateNorm(props);
 // ...
}

function App() {
 return (
 <>
 <MapMarker x={10} y={10} z={10} />
 <MapMarker x={10} y={10} z={10} />
 </>
 );
}
```

In this case the two `MapMarker`s look like they're doing the same work and calling `calculateNorm` with the same value of `{x: 10, y: 10, z:10}`. Even though the objects contain the same values, they are not the same object reference as each component creates its own `props` object.

React will call [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) on the input to verify if there is a cache hit.

```js {3,9}
import {cache} from 'react';

const calculateNorm = cache((x, y, z) => {
 // ...
});

function MapMarker(props) {
 // ✅ Good: Pass primitives to memoized function
 const length = calculateNorm(props.x, props.y, props.z);
 // ...
}

function App() {
 return (
 <>
 <MapMarker x={10} y={10} z={10} />
 <MapMarker x={10} y={10} z={10} />
 </>
 );
}
```

One way to address this could be to pass the vector dimensions to `calculateNorm`. This works because the dimensions themselves are primitives.

Another solution may be to pass the vector object itself as a prop to the component. We'll need to pass the same object to both component instances.

```js {3,9,14}
import {cache} from 'react';

const calculateNorm = cache((vector) => {
 // ...
});

function MapMarker(props) {
 // ✅ Good: Pass the same `vector` object
 const length = calculateNorm(props.vector);
 // ...
}

function App() {
 const vector = [10, 10, 10];
 return (
 <>
 <MapMarker vector={vector} />
 <MapMarker vector={vector} />
 </>
 );
}
```

---
title: cacheSignal
---

<RSC>

`cacheSignal` is currently only used with [React Server Components](/blog/2023/03/22/react-labs-what-we-have-been-working-on-march-2023#react-server-components).

</RSC>

<Intro>

`cacheSignal` allows you to know when the `cache()` lifetime is over.

```js
const signal = cacheSignal();
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `cacheSignal` {/*cachesignal*/}

Call `cacheSignal` to get an `AbortSignal`.

```js {3,7}
import {cacheSignal} from 'react';
async function Component() {
 await fetch(url, { signal: cacheSignal() });
}
```

When React has finished rendering, the `AbortSignal` will be aborted. This allows you to cancel any in-flight work that is no longer needed.
Rendering is considered finished when:
- React has successfully completed rendering
- the render was aborted
- the render has failed

#### Parameters {/*parameters*/}

This function does not accept any parameters.

#### Returns {/*returns*/}

`cacheSignal` returns an `AbortSignal` if called during rendering. Otherwise `cacheSignal()` returns `null`.

#### Caveats {/*caveats*/}

- `cacheSignal` is currently for use in [React Server Components](/reference/rsc/server-components) only. In Client Components, it will always return `null`. In the future it will also be used for Client Component when a client cache refreshes or invalidates. You should not assume it'll always be null on the client.
- If called outside of rendering, `cacheSignal` will return `null` to make it clear that the current scope isn't cached forever.

---

## Usage {/*usage*/}

### Cancel in-flight requests {/*cancel-in-flight-requests*/}

Call <CodeStep step={1}>`cacheSignal`</CodeStep> to abort in-flight requests.

```js [[1, 4, "cacheSignal()"]]
import {cache, cacheSignal} from 'react';
const dedupedFetch = cache(fetch);
async function Component() {
 await dedupedFetch(url, { signal: cacheSignal() });
}
```

<Pitfall>
You can't use `cacheSignal` to abort async work that was started outside of rendering e.g.

```js
import {cacheSignal} from 'react';
// 🚩 Pitfall: The request will not actually be aborted if the rendering of `Component` is finished.
const response = fetch(url, { signal: cacheSignal() });
async function Component() {
 await response;
}
```
</Pitfall>

### Ignore errors after React has finished rendering {/*ignore-errors-after-react-has-finished-rendering*/}

If a function throws, it may be due to cancellation (e.g. <CodeStep step={1}>the Database connection</CodeStep> has been closed). You can use the <CodeStep step={2}>`aborted` property</CodeStep> to check if the error was due to cancellation or a real error. You may want to <CodeStep step={3}>ignore errors</CodeStep> that were due to cancellation.

```js [[1, 2, "./database"], [2, 8, "cacheSignal()?.aborted"], [3, 12, "return null"]]
import {cacheSignal} from "react";
import {queryDatabase, logError} from "./database";

async function getData(id) {
 try {
 return await queryDatabase(id);
 } catch (x) {
 if (!cacheSignal()?.aborted) {
 // only log if it's a real error and not due to cancellation
 logError(x);
 }
 return null;
 }
}

async function Component({id}) {
 const data = await getData(id);
 if (data === null) {
 return <div>No data available</div>;
 }
 return <div>{data.name}</div>;
}
```

---
title: captureOwnerStack
---

<Intro>

`captureOwnerStack` reads the current Owner Stack in development and returns it as a string if available.

```js
const stack = captureOwnerStack();
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `captureOwnerStack()` {/*captureownerstack*/}

Call `captureOwnerStack` to get the current Owner Stack.

```js {5,5}
import * as React from 'react';

function Component() {
 if (process.env.NODE_ENV !== 'production') {
 const ownerStack = React.captureOwnerStack();
 console.log(ownerStack);
 }
}
```

#### Parameters {/*parameters*/}

`captureOwnerStack` does not take any parameters.

#### Returns {/*returns*/}

`captureOwnerStack` returns `string | null`.

Owner Stacks are available in
- Component render
- Effects (e.g. `useEffect`)
- React's event handlers (e.g. `<button onClick={...} />`)
- React error handlers ([React Root options](/reference/react-dom/client/createRoot#parameters) `onCaughtError`, `onRecoverableError`, and `onUncaughtError`)

If no Owner Stack is available, `null` is returned (see [Troubleshooting: The Owner Stack is `null`](#the-owner-stack-is-null)).

#### Caveats {/*caveats*/}

- Owner Stacks are only available in development. `captureOwnerStack` will always return `null` outside of development.

<DeepDive>

#### Owner Stack vs Component Stack {/*owner-stack-vs-component-stack*/}

The Owner Stack is different from the Component Stack available in React error handlers like [`errorInfo.componentStack` in `onUncaughtError`](/reference/react-dom/client/hydrateRoot#error-logging-in-production).

For example, consider the following code:

<Sandpack>

```js src/App.js
import {Suspense} from 'react';

function SubComponent({disabled}) {
 if (disabled) {
 throw new Error('disabled');
 }
}

export function Component({label}) {
 return (
 <fieldset>
 <legend>{label}</legend>
 <SubComponent key={label} disabled={label === 'disabled'} />
 </fieldset>
 );
}

function Navigation() {
 return null;
}

export default function App({children}) {
 return (
 <Suspense fallback="loading...">
 <main>
 <Navigation />
 {children}
 </main>
 </Suspense>
 );
}
```

```js src/index.js
import {captureOwnerStack} from 'react';
import {createRoot} from 'react-dom/client';
import App, {Component} from './App.js';
import './styles.css';

createRoot(document.createElement('div'), {
 onUncaughtError: (error, errorInfo) => {
 // The stacks are logged instead of showing them in the UI directly to
 // highlight that browsers will apply sourcemaps to the logged stacks.
 // Note that sourcemapping is only applied in the real browser console not
 // in the fake one displayed on this page.
 // Press "fork" to be able to view the sourcemapped stack in a real console.
 console.log(errorInfo.componentStack);
 console.log(captureOwnerStack());
 },
}).render(
 <App>
 <Component label="disabled" />
 </App>
);
```

```html public/index.html hidden
<!DOCTYPE html>
<html lang="en">
 <head>
 <meta charset="UTF-8" />
 <meta name="viewport" content="width=device-width, initial-scale=1.0" />
 <title>Document</title>
 </head>
 <body>
 <p>Check the console output.</p>
 </body>
</html>
```

</Sandpack>

`SubComponent` would throw an error.
The Component Stack of that error would be

```
at SubComponent
at fieldset
at Component
at main
at React.Suspense
at App
```

However, the Owner Stack would only read

```
at Component
```

Neither `App` nor the DOM components (e.g. `fieldset`) are considered Owners in this Stack since they didn't contribute to "creating" the node containing `SubComponent`. `App` and DOM components only forwarded the node. `App` just rendered the `children` node as opposed to `Component` which created a node containing `SubComponent` via `<SubComponent />`.

Neither `Navigation` nor `legend` are in the stack at all since it's only a sibling to a node containing `<SubComponent />`.

`SubComponent` is omitted because it's already part of the callstack.

</DeepDive>

## Usage {/*usage*/}

### Enhance a custom error overlay {/*enhance-a-custom-error-overlay*/}

```js [[1, 5, "console.error"], [4, 7, "captureOwnerStack"]]
import { captureOwnerStack } from "react";
import { instrumentedConsoleError } from "./errorOverlay";

const originalConsoleError = console.error;
console.error = function patchedConsoleError(...args) {
 originalConsoleError.apply(console, args);
 const ownerStack = captureOwnerStack();
 onConsoleError({
 // Keep in mind that in a real application, console.error can be
 // called with multiple arguments which you should account for.
 consoleMessage: args[0],
 ownerStack,
 });
};
```

If you intercept <CodeStep step={1}>`console.error`</CodeStep> calls to highlight them in an error overlay, you can call <CodeStep step={2}>`captureOwnerStack`</CodeStep> to include the Owner Stack.

<Sandpack>

```css src/styles.css
* {
 box-sizing: border-box;
}

body {
 font-family: sans-serif;
 margin: 20px;
 padding: 0;
}

h1 {
 margin-top: 0;
 font-size: 22px;
}

h2 {
 margin-top: 0;
 font-size: 20px;
}

code {
 font-size: 1.2em;
}

ul {
 padding-inline-start: 20px;
}

label, button { display: block; margin-bottom: 20px; }
html, body { min-height: 300px; }

#error-dialog {
 position: absolute;
 top: 0;
 right: 0;
 bottom: 0;
 left: 0;
 background-color: white;
 padding: 15px;
 opacity: 0.9;
 text-wrap: wrap;
 overflow: scroll;
}

.text-red {
 color: red;
}

.-mb-20 {
 margin-bottom: -20px;
}

.mb-0 {
 margin-bottom: 0;
}

.mb-10 {
 margin-bottom: 10px;
}

pre {
 text-wrap: wrap;
}

pre.nowrap {
 text-wrap: nowrap;
}

.hidden {
 display: none;
}
```

```html public/index.html hidden
<!DOCTYPE html>
<html>
<head>
 <title>My app</title>
</head>
<body>
<!--
 Error dialog in raw HTML
 since an error in the React app may crash.
-->
<div id="error-dialog" class="hidden">
 <h1 id="error-title" class="text-red">Error</h1>
 <p>
 <pre id="error-body"></pre>
 </p>
 <h2 class="-mb-20">Owner Stack:</h4>
 <pre id="error-owner-stack" class="nowrap"></pre>
 <button
 id="error-close"
 class="mb-10"
 onclick="document.getElementById('error-dialog').classList.add('hidden')"
 >
 Close
 </button>
</div>
<!-- This is the DOM node -->
<div id="root"></div>
</body>
</html>

```

```js src/errorOverlay.js

export function onConsoleError({ consoleMessage, ownerStack }) {
 const errorDialog = document.getElementById("error-dialog");
 const errorBody = document.getElementById("error-body");
 const errorOwnerStack = document.getElementById("error-owner-stack");

 // Display console.error() message
 errorBody.innerText = consoleMessage;

 // Display owner stack
 errorOwnerStack.innerText = ownerStack;

 // Show the dialog
 errorDialog.classList.remove("hidden");
}
```

```js src/index.js active
import { captureOwnerStack } from "react";
import { createRoot } from "react-dom/client";
import App from './App';
import { onConsoleError } from "./errorOverlay";
import './styles.css';

const originalConsoleError = console.error;
console.error = function patchedConsoleError(...args) {
 originalConsoleError.apply(console, args);
 const ownerStack = captureOwnerStack();
 onConsoleError({
 // Keep in mind that in a real application, console.error can be
 // called with multiple arguments which you should account for.
 consoleMessage: args[0],
 ownerStack,
 });
};

const container = document.getElementById("root");
createRoot(container).render(<App />);
```

```js src/App.js
function Component() {
 return <button onClick={() => console.error('Some console error')}>Trigger console.error()</button>;
}

export default function App() {
 return <Component />;
}
```

</Sandpack>

## Troubleshooting {/*troubleshooting*/}

### The Owner Stack is `null` {/*the-owner-stack-is-null*/}

The call of `captureOwnerStack` happened outside of a React controlled function e.g. in a `setTimeout` callback, after a `fetch` call or in a custom DOM event handler. During render, Effects, React event handlers, and React error handlers (e.g. `hydrateRoot#options.onCaughtError`) Owner Stacks should be available.

In the example below, clicking the button will log an empty Owner Stack because `captureOwnerStack` was called during a custom DOM event handler. The Owner Stack must be captured earlier e.g. by moving the call of `captureOwnerStack` into the Effect body.
<Sandpack>

```js
import {captureOwnerStack, useEffect} from 'react';

export default function App() {
 useEffect(() => {
 // Should call `captureOwnerStack` here.
 function handleEvent() {
 // Calling it in a custom DOM event handler is too late.
 // The Owner Stack will be `null` at this point.
 console.log('Owner Stack: ', captureOwnerStack());
 }

 document.addEventListener('click', handleEvent);

 return () => {
 document.removeEventListener('click', handleEvent);
 }
 })

 return <button>Click me to see that Owner Stacks are not available in custom DOM event handlers</button>;
}
```

</Sandpack>

### `captureOwnerStack` is not available {/*captureownerstack-is-not-available*/}

`captureOwnerStack` is only exported in development builds. It will be `undefined` in production builds. If `captureOwnerStack` is used in files that are bundled for production and development, you should conditionally access it from a namespace import.

```js
// Don't use named imports of `captureOwnerStack` in files that are bundled for development and production.
import {captureOwnerStack} from 'react';
// Use a namespace import instead and access `captureOwnerStack` conditionally.
import * as React from 'react';

if (process.env.NODE_ENV !== 'production') {
 const ownerStack = React.captureOwnerStack();
 console.log('Owner Stack', ownerStack);
}
```

---
title: cloneElement
---

<Pitfall>

Using `cloneElement` is uncommon and can lead to fragile code. [See common alternatives.](#alternatives)

</Pitfall>

<Intro>

`cloneElement` lets you create a new React element using another element as a starting point.

```js
const clonedElement = cloneElement(element, props, ...children)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `cloneElement(element, props, ...children)` {/*cloneelement*/}

Call `cloneElement` to create a React element based on the `element`, but with different `props` and `children`:

```js
import { cloneElement } from 'react';

// ...
const clonedElement = cloneElement(
 <Row title="Cabbage">
 Hello
 </Row>,
 { isHighlighted: true },
 'Goodbye'
);

console.log(clonedElement); // <Row title="Cabbage" isHighlighted={true}>Goodbye</Row>
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `element`: The `element` argument must be a valid React element. For example, it could be a JSX node like `<Something />`, the result of calling [`createElement`](/reference/react/createElement), or the result of another `cloneElement` call.

* `props`: The `props` argument must either be an object or `null`. If you pass `null`, the cloned element will retain all of the original `element.props`. Otherwise, for every prop in the `props` object, the returned element will "prefer" the value from `props` over the value from `element.props`. The rest of the props will be filled from the original `element.props`. If you pass `props.key` or `props.ref`, they will replace the original ones.

* **optional** `...children`: Zero or more child nodes. They can be any React nodes, including React elements, strings, numbers, [portals](/reference/react-dom/createPortal), empty nodes (`null`, `undefined`, `true`, and `false`), and arrays of React nodes. If you don't pass any `...children` arguments, the original `element.props.children` will be preserved.

#### Returns {/*returns*/}

`cloneElement` returns a React element object with a few properties:

* `type`: Same as `element.type`.
* `props`: The result of shallowly merging `element.props` with the overriding `props` you have passed.
* `ref`: The original `element.ref`, unless it was overridden by `props.ref`.
* `key`: The original `element.key`, unless it was overridden by `props.key`.

Usually, you'll return the element from your component or make it a child of another element. Although you may read the element's properties, it's best to treat every element as opaque after it's created, and only render it.

#### Caveats {/*caveats*/}

* Cloning an element **does not modify the original element.**

* You should only **pass children as multiple arguments to `cloneElement` if they are all statically known,** like `cloneElement(element, null, child1, child2, child3)`. If your children are dynamic, pass the entire array as the third argument: `cloneElement(element, null, listItems)`. This ensures that React will [warn you about missing `key`s](/learn/rendering-lists#keeping-list-items-in-order-with-key) for any dynamic lists. For static lists this is not necessary because they never reorder.

* `cloneElement` makes it harder to trace the data flow, so **try the [alternatives](#alternatives) instead.**

---

## Usage {/*usage*/}

### Overriding props of an element {/*overriding-props-of-an-element*/}

To override the props of some <CodeStep step={1}>React element</CodeStep>, pass it to `cloneElement` with the <CodeStep step={2}>props you want to override</CodeStep>:

```js [[1, 5, "<Row title=\\"Cabbage\\" />"], [2, 6, "{ isHighlighted: true }"], [3, 4, "clonedElement"]]
import { cloneElement } from 'react';

// ...
const clonedElement = cloneElement(
 <Row title="Cabbage" />,
 { isHighlighted: true }
);
```

Here, the resulting <CodeStep step={3}>cloned element</CodeStep> will be `<Row title="Cabbage" isHighlighted={true} />`.

**Let's walk through an example to see when it's useful.**

Imagine a `List` component that renders its [`children`](/learn/passing-props-to-a-component#passing-jsx-as-children) as a list of selectable rows with a "Next" button that changes which row is selected. The `List` component needs to render the selected `Row` differently, so it clones every `<Row>` child that it has received, and adds an extra `isHighlighted: true` or `isHighlighted: false` prop:

```js {6-8}
export default function List({ children }) {
 const [selectedIndex, setSelectedIndex] = useState(0);
 return (
 <div className="List">
 {Children.map(children, (child, index) =>
 cloneElement(child, {
 isHighlighted: index === selectedIndex
 })
 )}
```

Let's say the original JSX received by `List` looks like this:

```js {2-4}
<List>
 <Row title="Cabbage" />
 <Row title="Garlic" />
 <Row title="Apple" />
</List>
```

By cloning its children, the `List` can pass extra information to every `Row` inside. The result looks like this:

```js {4,8,12}
<List>
 <Row
 title="Cabbage"
 isHighlighted={true}
 />
 <Row
 title="Garlic"
 isHighlighted={false}
 />
 <Row
 title="Apple"
 isHighlighted={false}
 />
</List>
```

Notice how pressing "Next" updates the state of the `List`, and highlights a different row:

<Sandpack>

```js
import List from './List.js';
import Row from './Row.js';
import { products } from './data.js';

export default function App() {
 return (
 <List>
 {products.map(product =>
 <Row
 key={product.id}
 title={product.title}
 />
 )}
 </List>
 );
}
```

```js src/List.js active
import { Children, cloneElement, useState } from 'react';

export default function List({ children }) {
 const [selectedIndex, setSelectedIndex] = useState(0);
 return (
 <div className="List">
 {Children.map(children, (child, index) =>
 cloneElement(child, {
 isHighlighted: index === selectedIndex
 })
 )}
 <hr />
 <button onClick={() => {
 setSelectedIndex(i =>
 (i + 1) % Children.count(children)
 );
 }}>
 Next
 </button>
 </div>
 );
}
```

```js src/Row.js
export default function Row({ title, isHighlighted }) {
 return (
 <div className={[
 'Row',
 isHighlighted ? 'RowHighlighted' : ''
 ].join(' ')}>
 {title}
 </div>
 );
}
```

```js src/data.js
export const products = [
 { title: 'Cabbage', id: 1 },
 { title: 'Garlic', id: 2 },
 { title: 'Apple', id: 3 },
];
```

```css
.List {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}

.RowHighlighted {
 background: #ffa;
}

button {
 height: 40px;
 font-size: 20px;
}
```

</Sandpack>

To summarize, the `List` cloned the `<Row />` elements it received and added an extra prop to them.

<Pitfall>

Cloning children makes it hard to tell how the data flows through your app. Try one of the [alternatives.](#alternatives)

</Pitfall>

---

## Alternatives {/*alternatives*/}

### Passing data with a render prop {/*passing-data-with-a-render-prop*/}

Instead of using `cloneElement`, consider accepting a *render prop* like `renderItem`. Here, `List` receives `renderItem` as a prop. `List` calls `renderItem` for every item and passes `isHighlighted` as an argument:

```js {1,7}
export default function List({ items, renderItem }) {
 const [selectedIndex, setSelectedIndex] = useState(0);
 return (
 <div className="List">
 {items.map((item, index) => {
 const isHighlighted = index === selectedIndex;
 return renderItem(item, isHighlighted);
 })}
```

The `renderItem` prop is called a "render prop" because it's a prop that specifies how to render something. For example, you can pass a `renderItem` implementation that renders a `<Row>` with the given `isHighlighted` value:

```js {3,7}
<List
 items={products}
 renderItem={(product, isHighlighted) =>
 <Row
 key={product.id}
 title={product.title}
 isHighlighted={isHighlighted}
 />
 }
/>
```

The end result is the same as with `cloneElement`:

```js {4,8,12}
<List>
 <Row
 title="Cabbage"
 isHighlighted={true}
 />
 <Row
 title="Garlic"
 isHighlighted={false}
 />
 <Row
 title="Apple"
 isHighlighted={false}
 />
</List>
```

However, you can clearly trace where the `isHighlighted` value is coming from.

<Sandpack>

```js
import List from './List.js';
import Row from './Row.js';
import { products } from './data.js';

export default function App() {
 return (
 <List
 items={products}
 renderItem={(product, isHighlighted) =>
 <Row
 key={product.id}
 title={product.title}
 isHighlighted={isHighlighted}
 />
 }
 />
 );
}
```

```js src/List.js active
import { useState } from 'react';

export default function List({ items, renderItem }) {
 const [selectedIndex, setSelectedIndex] = useState(0);
 return (
 <div className="List">
 {items.map((item, index) => {
 const isHighlighted = index === selectedIndex;
 return renderItem(item, isHighlighted);
 })}
 <hr />
 <button onClick={() => {
 setSelectedIndex(i =>
 (i + 1) % items.length
 );
 }}>
 Next
 </button>
 </div>
 );
}
```

```js src/Row.js
export default function Row({ title, isHighlighted }) {
 return (
 <div className={[
 'Row',
 isHighlighted ? 'RowHighlighted' : ''
 ].join(' ')}>
 {title}
 </div>
 );
}
```

```js src/data.js
export const products = [
 { title: 'Cabbage', id: 1 },
 { title: 'Garlic', id: 2 },
 { title: 'Apple', id: 3 },
];
```

```css
.List {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}

.RowHighlighted {
 background: #ffa;
}

button {
 height: 40px;
 font-size: 20px;
}
```

</Sandpack>

This pattern is preferred to `cloneElement` because it is more explicit.

---

### Passing data through context {/*passing-data-through-context*/}

Another alternative to `cloneElement` is to [pass data through context.](/learn/passing-data-deeply-with-context)

For example, you can call [`createContext`](/reference/react/createContext) to define a `HighlightContext`:

```js
export const HighlightContext = createContext(false);
```

Your `List` component can wrap every item it renders into a `HighlightContext` provider:

```js {8,10}
export default function List({ items, renderItem }) {
 const [selectedIndex, setSelectedIndex] = useState(0);
 return (
 <div className="List">
 {items.map((item, index) => {
 const isHighlighted = index === selectedIndex;
 return (
 <HighlightContext key={item.id} value={isHighlighted}>
 {renderItem(item)}
 </HighlightContext>
 );
 })}
```

With this approach, `Row` does not need to receive an `isHighlighted` prop at all. Instead, it reads the context:

```js src/Row.js {2}
export default function Row({ title }) {
 const isHighlighted = useContext(HighlightContext);
 // ...
```

This allows the calling component to not know or worry about passing `isHighlighted` to `<Row>`:

```js {4}
<List
 items={products}
 renderItem={product =>
 <Row title={product.title} />
 }
/>
```

Instead, `List` and `Row` coordinate the highlighting logic through context.

<Sandpack>

```js
import List from './List.js';
import Row from './Row.js';
import { products } from './data.js';

export default function App() {
 return (
 <List
 items={products}
 renderItem={(product) =>
 <Row title={product.title} />
 }
 />
 );
}
```

```js src/List.js active
import { useState } from 'react';
import { HighlightContext } from './HighlightContext.js';

export default function List({ items, renderItem }) {
 const [selectedIndex, setSelectedIndex] = useState(0);
 return (
 <div className="List">
 {items.map((item, index) => {
 const isHighlighted = index === selectedIndex;
 return (
 <HighlightContext
 key={item.id}
 value={isHighlighted}
 >
 {renderItem(item)}
 </HighlightContext>
 );
 })}
 <hr />
 <button onClick={() => {
 setSelectedIndex(i =>
 (i + 1) % items.length
 );
 }}>
 Next
 </button>
 </div>
 );
}
```

```js src/Row.js
import { useContext } from 'react';
import { HighlightContext } from './HighlightContext.js';

export default function Row({ title }) {
 const isHighlighted = useContext(HighlightContext);
 return (
 <div className={[
 'Row',
 isHighlighted ? 'RowHighlighted' : ''
 ].join(' ')}>
 {title}
 </div>
 );
}
```

```js src/HighlightContext.js
import { createContext } from 'react';

export const HighlightContext = createContext(false);
```

```js src/data.js
export const products = [
 { title: 'Cabbage', id: 1 },
 { title: 'Garlic', id: 2 },
 { title: 'Apple', id: 3 },
];
```

```css
.List {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}

.RowHighlighted {
 background: #ffa;
}

button {
 height: 40px;
 font-size: 20px;
}
```

</Sandpack>

[Learn more about passing data through context.](/reference/react/useContext#passing-data-deeply-into-the-tree)

---

### Extracting logic into a custom Hook {/*extracting-logic-into-a-custom-hook*/}

Another approach you can try is to extract the "non-visual" logic into your own Hook, and use the information returned by your Hook to decide what to render. For example, you could write a `useList` custom Hook like this:

```js
import { useState } from 'react';

export default function useList(items) {
 const [selectedIndex, setSelectedIndex] = useState(0);

 function onNext() {
 setSelectedIndex(i =>
 (i + 1) % items.length
 );
 }

 const selected = items[selectedIndex];
 return [selected, onNext];
}
```

Then you could use it like this:

```js {2,9,13}
export default function App() {
 const [selected, onNext] = useList(products);
 return (
 <div className="List">
 {products.map(product =>
 <Row
 key={product.id}
 title={product.title}
 isHighlighted={selected === product}
 />
 )}
 <hr />
 <button onClick={onNext}>
 Next
 </button>
 </div>
 );
}
```

The data flow is explicit, but the state is inside the `useList` custom Hook that you can use from any component:

<Sandpack>

```js
import Row from './Row.js';
import useList from './useList.js';
import { products } from './data.js';

export default function App() {
 const [selected, onNext] = useList(products);
 return (
 <div className="List">
 {products.map(product =>
 <Row
 key={product.id}
 title={product.title}
 isHighlighted={selected === product}
 />
 )}
 <hr />
 <button onClick={onNext}>
 Next
 </button>
 </div>
 );
}
```

```js src/useList.js
import { useState } from 'react';

export default function useList(items) {
 const [selectedIndex, setSelectedIndex] = useState(0);

 function onNext() {
 setSelectedIndex(i =>
 (i + 1) % items.length
 );
 }

 const selected = items[selectedIndex];
 return [selected, onNext];
}
```

```js src/Row.js
export default function Row({ title, isHighlighted }) {
 return (
 <div className={[
 'Row',
 isHighlighted ? 'RowHighlighted' : ''
 ].join(' ')}>
 {title}
 </div>
 );
}
```

```js src/data.js
export const products = [
 { title: 'Cabbage', id: 1 },
 { title: 'Garlic', id: 2 },
 { title: 'Apple', id: 3 },
];
```

```css
.List {
 display: flex;
 flex-direction: column;
 border: 2px solid grey;
 padding: 5px;
}

.Row {
 border: 2px dashed black;
 padding: 5px;
 margin: 5px;
}

.RowHighlighted {
 background: #ffa;
}

button {
 height: 40px;
 font-size: 20px;
}
```

</Sandpack>

This approach is particularly useful if you want to reuse this logic between different components.

---
title: "Built-in React Components"
---

<Intro>

React exposes a few built-in components that you can use in your JSX.

</Intro>

---

## Built-in components {/*built-in-components*/}

* [`<Fragment>`](/reference/react/Fragment), alternatively written as `<>...</>`, lets you group multiple JSX nodes together.
* [`<Profiler>`](/reference/react/Profiler) lets you measure rendering performance of a React tree programmatically.
* [`<Suspense>`](/reference/react/Suspense) lets you display a fallback while the child components are loading.
* [`<StrictMode>`](/reference/react/StrictMode) enables extra development-only checks that help you find bugs early.
* [`<Activity>`](/reference/react/Activity) lets you hide and restore the UI and internal state of its children.

---

## Your own components {/*your-own-components*/}

You can also [define your own components](/learn/your-first-component) as JavaScript functions.

---
title: createContext
---

<Intro>

`createContext` lets you create a [context](/learn/passing-data-deeply-with-context) that components can provide or read.

```js
const SomeContext = createContext(defaultValue)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `createContext(defaultValue)` {/*createcontext*/}

Call `createContext` outside of any components to create a context.

```js
import { createContext } from 'react';

const ThemeContext = createContext('light');
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `defaultValue`: The value that you want the context to have when there is no matching context provider in the tree above the component that reads context. If you don't have any meaningful default value, specify `null`. The default value is meant as a "last resort" fallback. It is static and never changes over time.

#### Returns {/*returns*/}

`createContext` returns a context object.

**The context object itself does not hold any information.** It represents _which_ context other components read or provide. Typically, you will use [`SomeContext`](#provider) in components above to specify the context value, and call [`useContext(SomeContext)`](/reference/react/useContext) in components below to read it. The context object has a few properties:

* `SomeContext` lets you provide the context value to components.
* `SomeContext.Consumer` is an alternative and rarely used way to read the context value.
* `SomeContext.Provider` is a legacy way to provide the context value before React 19.

---

### `SomeContext` Provider {/*provider*/}

Wrap your components into a context provider to specify the value of this context for all components inside:

```js
function App() {
 const [theme, setTheme] = useState('light');
 // ...
 return (
 <ThemeContext value={theme}>
 <Page />
 </ThemeContext>
 );
}
```

<Note>

Starting in React 19, you can render `<SomeContext>` as a provider.

In older versions of React, use `<SomeContext.Provider>`.

</Note>

#### Props {/*provider-props*/}

* `value`: The value that you want to pass to all the components reading this context inside this provider, no matter how deep. The context value can be of any type. A component calling [`useContext(SomeContext)`](/reference/react/useContext) inside of the provider receives the `value` of the innermost corresponding context provider above it.

---

### `SomeContext.Consumer` {/*consumer*/}

Before `useContext` existed, there was an older way to read context:

```js
function Button() {
 // 🟡 Legacy way (not recommended)
 return (
 <ThemeContext.Consumer>
 {theme => (
 <button className={theme} />
 )}
 </ThemeContext.Consumer>
 );
}
```

Although this older way still works, **newly written code should read context with [`useContext()`](/reference/react/useContext) instead:**

```js
function Button() {
 // ✅ Recommended way
 const theme = useContext(ThemeContext);
 return <button className={theme} />;
}
```

#### Props {/*consumer-props*/}

* `children`: A function. React will call the function you pass with the current context value determined by the same algorithm as [`useContext()`](/reference/react/useContext) does, and render the result you return from this function. React will also re-run this function and update the UI whenever the context from the parent components changes.

---

## Usage {/*usage*/}

### Creating context {/*creating-context*/}

Context lets components [pass information deep down](/learn/passing-data-deeply-with-context) without explicitly passing props.

Call `createContext` outside any components to create one or more contexts.

```js [[1, 3, "ThemeContext"], [1, 4, "AuthContext"], [3, 3, "'light'"], [3, 4, "null"]]
import { createContext } from 'react';

const ThemeContext = createContext('light');
const AuthContext = createContext(null);
```

`createContext` returns a <CodeStep step={1}>context object</CodeStep>. Components can read context by passing it to [`useContext()`](/reference/react/useContext):

```js [[1, 2, "ThemeContext"], [1, 7, "AuthContext"]]
function Button() {
 const theme = useContext(ThemeContext);
 // ...
}

function Profile() {
 const currentUser = useContext(AuthContext);
 // ...
}
```

By default, the values they receive will be the <CodeStep step={3}>default values</CodeStep> you have specified when creating the contexts. However, by itself this isn't useful because the default values never change.

Context is useful because you can **provide other, dynamic values from your components:**

```js {8-9,11-12}
function App() {
 const [theme, setTheme] = useState('dark');
 const [currentUser, setCurrentUser] = useState({ name: 'Taylor' });

 // ...

 return (
 <ThemeContext value={theme}>
 <AuthContext value={currentUser}>
 <Page />
 </AuthContext>
 </ThemeContext>
 );
}
```

Now the `Page` component and any components inside it, no matter how deep, will "see" the passed context values. If the passed context values change, React will re-render the components reading the context as well.

[Read more about reading and providing context and see examples.](/reference/react/useContext)

---

### Importing and exporting context from a file {/*importing-and-exporting-context-from-a-file*/}

Often, components in different files will need access to the same context. This is why it's common to declare contexts in a separate file. Then you can use the [`export` statement](https://developer.mozilla.org/en-US/docs/web/javascript/reference/statements/export) to make context available for other files:

```js {4-5}
// Contexts.js
import { createContext } from 'react';

export const ThemeContext = createContext('light');
export const AuthContext = createContext(null);
```

Components declared in other files can then use the [`import`](https://developer.mozilla.org/en-US/docs/web/javascript/reference/statements/import) statement to read or provide this context:

```js {2}
// Button.js
import { ThemeContext } from './Contexts.js';

function Button() {
 const theme = useContext(ThemeContext);
 // ...
}
```

```js {2}
// App.js
import { ThemeContext, AuthContext } from './Contexts.js';

function App() {
 // ...
 return (
 <ThemeContext value={theme}>
 <AuthContext value={currentUser}>
 <Page />
 </AuthContext>
 </ThemeContext>
 );
}
```

This works similar to [importing and exporting components.](/learn/importing-and-exporting-components)

---

## Troubleshooting {/*troubleshooting*/}

### I can't find a way to change the context value {/*i-cant-find-a-way-to-change-the-context-value*/}

Code like this specifies the *default* context value:

```js
const ThemeContext = createContext('light');
```

This value never changes. React only uses this value as a fallback if it can't find a matching provider above.

To make context change over time, [add state and wrap components in a context provider.](/reference/react/useContext#updating-data-passed-via-context)

---
title: createElement
---

<Intro>

`createElement` lets you create a React element. It serves as an alternative to writing [JSX.](/learn/writing-markup-with-jsx)

```js
const element = createElement(type, props, ...children)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `createElement(type, props, ...children)` {/*createelement*/}

Call `createElement` to create a React element with the given `type`, `props`, and `children`.

```js
import { createElement } from 'react';

function Greeting({ name }) {
 return createElement(
 'h1',
 { className: 'greeting' },
 'Hello'
 );
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `type`: The `type` argument must be a valid React component type. For example, it could be a tag name string (such as `'div'` or `'span'`), or a React component (a function, a class, or a special component like [`Fragment`](/reference/react/Fragment)).

* `props`: The `props` argument must either be an object or `null`. If you pass `null`, it will be treated the same as an empty object. React will create an element with props matching the `props` you have passed. Note that `ref` and `key` from your `props` object are special and will *not* be available as `element.props.ref` and `element.props.key` on the returned `element`. They will be available as `element.ref` and `element.key`.

* **optional** `...children`: Zero or more child nodes. They can be any React nodes, including React elements, strings, numbers, [portals](/reference/react-dom/createPortal), empty nodes (`null`, `undefined`, `true`, and `false`), and arrays of React nodes.

#### Returns {/*returns*/}

`createElement` returns a React element object with a few properties:

* `type`: The `type` you have passed.
* `props`: The `props` you have passed except for `ref` and `key`.
* `ref`: The `ref` you have passed. If missing, `null`.
* `key`: The `key` you have passed, coerced to a string. If missing, `null`.

Usually, you'll return the element from your component or make it a child of another element. Although you may read the element's properties, it's best to treat every element as opaque after it's created, and only render it.

#### Caveats {/*caveats*/}

* You must **treat React elements and their props as [immutable](https://en.wikipedia.org/wiki/Immutable_object)** and never change their contents after creation. In development, React will [freeze](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/freeze) the returned element and its `props` property shallowly to enforce this.

* When you use JSX, **you must start a tag with a capital letter to render your own custom component.** In other words, `<Something />` is equivalent to `createElement(Something)`, but `<something />` (lowercase) is equivalent to `createElement('something')` (note it's a string, so it will be treated as a built-in HTML tag).

* You should only **pass children as multiple arguments to `createElement` if they are all statically known,** like `createElement('h1', {}, child1, child2, child3)`. If your children are dynamic, pass the entire array as the third argument: `createElement('ul', {}, listItems)`. This ensures that React will [warn you about missing `key`s](/learn/rendering-lists#keeping-list-items-in-order-with-key) for any dynamic lists. For static lists this is not necessary because they never reorder.

---

## Usage {/*usage*/}

### Creating an element without JSX {/*creating-an-element-without-jsx*/}

If you don't like [JSX](/learn/writing-markup-with-jsx) or can't use it in your project, you can use `createElement` as an alternative.

To create an element without JSX, call `createElement` with some <CodeStep step={1}>type</CodeStep>, <CodeStep step={2}>props</CodeStep>, and <CodeStep step={3}>children</CodeStep>:

```js [[1, 5, "'h1'"], [2, 6, "{ className: 'greeting' }"], [3, 7, "'Hello ',"], [3, 8, "createElement('i', null, name),"], [3, 9, "'. Welcome!'"]]
import { createElement } from 'react';

function Greeting({ name }) {
 return createElement(
 'h1',
 { className: 'greeting' },
 'Hello ',
 createElement('i', null, name),
 '. Welcome!'
 );
}
```

The <CodeStep step={3}>children</CodeStep> are optional, and you can pass as many as you need (the example above has three children). This code will display a `<h1>` header with a greeting. For comparison, here is the same example rewritten with JSX:

```js [[1, 3, "h1"], [2, 3, "className=\\"greeting\\""], [3, 4, "Hello <i>{name}</i>. Welcome!"], [1, 5, "h1"]]
function Greeting({ name }) {
 return (
 <h1 className="greeting">
 Hello <i>{name}</i>. Welcome!
 </h1>
 );
}
```

To render your own React component, pass a function like `Greeting` as the <CodeStep step={1}>type</CodeStep> instead of a string like `'h1'`:

```js [[1, 2, "Greeting"], [2, 2, "{ name: 'Taylor' }"]]
export default function App() {
 return createElement(Greeting, { name: 'Taylor' });
}
```

With JSX, it would look like this:

```js [[1, 2, "Greeting"], [2, 2, "name=\\"Taylor\\""]]
export default function App() {
 return <Greeting name="Taylor" />;
}
```

Here is a complete example written with `createElement`:

<Sandpack>

```js
import { createElement } from 'react';

function Greeting({ name }) {
 return createElement(
 'h1',
 { className: 'greeting' },
 'Hello ',
 createElement('i', null, name),
 '. Welcome!'
 );
}

export default function App() {
 return createElement(
 Greeting,
 { name: 'Taylor' }
 );
}
```

```css
.greeting {
 color: darkgreen;
 font-family: Georgia;
}
```

</Sandpack>

And here is the same example written using JSX:

<Sandpack>

```js
function Greeting({ name }) {
 return (
 <h1 className="greeting">
 Hello <i>{name}</i>. Welcome!
 </h1>
 );
}

export default function App() {
 return <Greeting name="Taylor" />;
}
```

```css
.greeting {
 color: darkgreen;
 font-family: Georgia;
}
```

</Sandpack>

Both coding styles are fine, so you can use whichever one you prefer for your project. The main benefit of using JSX compared to `createElement` is that it's easy to see which closing tag corresponds to which opening tag.

<DeepDive>

#### What is a React element, exactly? {/*what-is-a-react-element-exactly*/}

An element is a lightweight description of a piece of the user interface. For example, both `<Greeting name="Taylor" />` and `createElement(Greeting, { name: 'Taylor' })` produce an object like this:

```js
// Slightly simplified
{
 type: Greeting,
 props: {
 name: 'Taylor'
 },
 key: null,
 ref: null,
}
```

**Note that creating this object does not render the `Greeting` component or create any DOM elements.**

A React element is more like a description--an instruction for React to later render the `Greeting` component. By returning this object from your `App` component, you tell React what to do next.

Creating elements is extremely cheap so you don't need to try to optimize or avoid it.

</DeepDive>

---
title: createRef
---

<Pitfall>

`createRef` is mostly used for [class components.](/reference/react/Component) Function components typically rely on [`useRef`](/reference/react/useRef) instead.

</Pitfall>

<Intro>

`createRef` creates a [ref](/learn/referencing-values-with-refs) object which can contain arbitrary value.

```js
class MyInput extends Component {
 inputRef = createRef();
 // ...
}
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `createRef()` {/*createref*/}

Call `createRef` to declare a [ref](/learn/referencing-values-with-refs) inside a [class component.](/reference/react/Component)

```js
import { createRef, Component } from 'react';

class MyComponent extends Component {
 intervalRef = createRef();
 inputRef = createRef();
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

`createRef` takes no parameters.

#### Returns {/*returns*/}

`createRef` returns an object with a single property:

* `current`: Initially, it's set to the `null`. You can later set it to something else. If you pass the ref object to React as a `ref` attribute to a JSX node, React will set its `current` property.

#### Caveats {/*caveats*/}

* `createRef` always returns a *different* object. It's equivalent to writing `{ current: null }` yourself.
* In a function component, you probably want [`useRef`](/reference/react/useRef) instead which always returns the same object.
* `const ref = useRef()` is equivalent to `const [ref, _] = useState(() => createRef(null))`.

---

## Usage {/*usage*/}

### Declaring a ref in a class component {/*declaring-a-ref-in-a-class-component*/}

To declare a ref inside a [class component,](/reference/react/Component) call `createRef` and assign its result to a class field:

```js {4}
import { Component, createRef } from 'react';

class Form extends Component {
 inputRef = createRef();

 // ...
}
```

If you now pass `ref={this.inputRef}` to an `<input>` in your JSX, React will populate `this.inputRef.current` with the input DOM node. For example, here is how you make a button that focuses the input:

<Sandpack>

```js
import { Component, createRef } from 'react';

export default class Form extends Component {
 inputRef = createRef();

 handleClick = () => {
 this.inputRef.current.focus();
 }

 render() {
 return (
 <>
 <input ref={this.inputRef} />
 <button onClick={this.handleClick}>
 Focus the input
 </button>
 </>
 );
 }
}
```

</Sandpack>

<Pitfall>

`createRef` is mostly used for [class components.](/reference/react/Component) Function components typically rely on [`useRef`](/reference/react/useRef) instead.

</Pitfall>

---

## Alternatives {/*alternatives*/}

### Migrating from a class with `createRef` to a function with `useRef` {/*migrating-from-a-class-with-createref-to-a-function-with-useref*/}

We recommend using function components instead of [class components](/reference/react/Component) in new code. If you have some existing class components using `createRef`, here is how you can convert them. This is the original code:

<Sandpack>

```js
import { Component, createRef } from 'react';

export default class Form extends Component {
 inputRef = createRef();

 handleClick = () => {
 this.inputRef.current.focus();
 }

 render() {
 return (
 <>
 <input ref={this.inputRef} />
 <button onClick={this.handleClick}>
 Focus the input
 </button>
 </>
 );
 }
}
```

</Sandpack>

When you [convert this component from a class to a function,](/reference/react/Component#alternatives) replace calls to `createRef` with calls to [`useRef`:](/reference/react/useRef)

<Sandpack>

```js
import { useRef } from 'react';

export default function Form() {
 const inputRef = useRef(null);

 function handleClick() {
 inputRef.current.focus();
 }

 return (
 <>
 <input ref={inputRef} />
 <button onClick={handleClick}>
 Focus the input
 </button>
 </>
 );
}
```

</Sandpack>

---
title: experimental_taintObjectReference
version: experimental
---

<Experimental>

**This API is experimental and is not available in a stable version of React yet.**

You can try it by upgrading React packages to the most recent experimental version:

- `react@experimental`
- `react-dom@experimental`
- `eslint-plugin-react-hooks@experimental`

Experimental versions of React may contain bugs. Don't use them in production.

This API is only available inside React Server Components.

</Experimental>

<Intro>

`taintObjectReference` lets you prevent a specific object instance from being passed to a Client Component like a `user` object.

```js
experimental_taintObjectReference(message, object);
```

To prevent passing a key, hash or token, see [`taintUniqueValue`](/reference/react/experimental_taintUniqueValue).

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `taintObjectReference(message, object)` {/*taintobjectreference*/}

Call `taintObjectReference` with an object to register it with React as something that should not be allowed to be passed to the Client as is:

```js
import {experimental_taintObjectReference} from 'react';

experimental_taintObjectReference(
 'Do not pass ALL environment variables to the client.',
 process.env
);
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `message`: The message you want to display if the object gets passed to a Client Component. This message will be displayed as a part of the Error that will be thrown if the object gets passed to a Client Component.

* `object`: The object to be tainted. Functions and class instances can be passed to `taintObjectReference` as `object`. Functions and classes are already blocked from being passed to Client Components but the React's default error message will be replaced by what you defined in `message`. When a specific instance of a Typed Array is passed to `taintObjectReference` as `object`, any other copies of the Typed Array will not be tainted.

#### Returns {/*returns*/}

`experimental_taintObjectReference` returns `undefined`.

#### Caveats {/*caveats*/}

- Recreating or cloning a tainted object creates a new untainted object which may contain sensitive data. For example, if you have a tainted `user` object, `const userInfo = {name: user.name, ssn: user.ssn}` or `{...user}` will create new objects which are not tainted. `taintObjectReference` only protects against simple mistakes when the object is passed through to a Client Component unchanged.

<Pitfall>

**Do not rely on just tainting for security.** Tainting an object doesn't prevent leaking of every possible derived value. For example, the clone of a tainted object will create a new untainted object. Using data from a tainted object (e.g. `{secret: taintedObj.secret}`) will create a new value or object that is not tainted. Tainting is a layer of protection; a secure app will have multiple layers of protection, well designed APIs, and isolation patterns.

</Pitfall>

---

## Usage {/*usage*/}

### Prevent user data from unintentionally reaching the client {/*prevent-user-data-from-unintentionally-reaching-the-client*/}

A Client Component should never accept objects that carry sensitive data. Ideally, the data fetching functions should not expose data that the current user should not have access to. Sometimes mistakes happen during refactoring. To protect against these mistakes happening down the line we can "taint" the user object in our data API.

```js
import {experimental_taintObjectReference} from 'react';

export async function getUser(id) {
 const user = await db`SELECT * FROM users WHERE id = ${id}`;
 experimental_taintObjectReference(
 'Do not pass the entire user object to the client. ' +
 'Instead, pick off the specific properties you need for this use case.',
 user,
 );
 return user;
}
```

Now whenever anyone tries to pass this object to a Client Component, an error will be thrown with the passed in error message instead.

<DeepDive>

#### Protecting against leaks in data fetching {/*protecting-against-leaks-in-data-fetching*/}

If you're running a Server Components environment that has access to sensitive data, you have to be careful not to pass objects straight through:

```js
// api.js
export async function getUser(id) {
 const user = await db`SELECT * FROM users WHERE id = ${id}`;
 return user;
}
```

```js
import { getUser } from 'api.js';
import { InfoCard } from 'components.js';

export async function Profile(props) {
 const user = await getUser(props.userId);
 // DO NOT DO THIS
 return <InfoCard user={user} />;
}
```

```js
// components.js
"use client";

export async function InfoCard({ user }) {
 return <div>{user.name}</div>;
}
```

Ideally, the `getUser` should not expose data that the current user should not have access to. To prevent passing the `user` object to a Client Component down the line we can "taint" the user object:

```js
// api.js
import {experimental_taintObjectReference} from 'react';

export async function getUser(id) {
 const user = await db`SELECT * FROM users WHERE id = ${id}`;
 experimental_taintObjectReference(
 'Do not pass the entire user object to the client. ' +
 'Instead, pick off the specific properties you need for this use case.',
 user,
 );
 return user;
}
```

Now if anyone tries to pass the `user` object to a Client Component, an error will be thrown with the passed in error message.

</DeepDive>

---
title: experimental_taintUniqueValue
version: experimental
---

<Experimental>

**This API is experimental and is not available in a stable version of React yet.**

You can try it by upgrading React packages to the most recent experimental version:

- `react@experimental`
- `react-dom@experimental`
- `eslint-plugin-react-hooks@experimental`

Experimental versions of React may contain bugs. Don't use them in production.

This API is only available inside [React Server Components](/reference/rsc/use-client).

</Experimental>

<Intro>

`taintUniqueValue` lets you prevent unique values from being passed to Client Components like passwords, keys, or tokens.

```js
taintUniqueValue(errMessage, lifetime, value)
```

To prevent passing an object containing sensitive data, see [`taintObjectReference`](/reference/react/experimental_taintObjectReference).

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `taintUniqueValue(message, lifetime, value)` {/*taintuniquevalue*/}

Call `taintUniqueValue` with a password, token, key or hash to register it with React as something that should not be allowed to be passed to the Client as is:

```js
import {experimental_taintUniqueValue} from 'react';

experimental_taintUniqueValue(
 'Do not pass secret keys to the client.',
 process,
 process.env.SECRET_KEY
);
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `message`: The message you want to display if `value` is passed to a Client Component. This message will be displayed as a part of the Error that will be thrown if `value` is passed to a Client Component.

* `lifetime`: Any object that indicates how long `value` should be tainted. `value` will be blocked from being sent to any Client Component while this object still exists. For example, passing `globalThis` blocks the value for the lifetime of an app. `lifetime` is typically an object whose properties contains `value`.

* `value`: A string, bigint or TypedArray. `value` must be a unique sequence of characters or bytes with high entropy such as a cryptographic token, private key, hash, or a long password. `value` will be blocked from being sent to any Client Component.

#### Returns {/*returns*/}

`experimental_taintUniqueValue` returns `undefined`.

#### Caveats {/*caveats*/}

* Deriving new values from tainted values can compromise tainting protection. New values created by uppercasing tainted values, concatenating tainted string values into a larger string, converting tainted values to base64, substringing tainted values, and other similar transformations are not tainted unless you explicitly call `taintUniqueValue` on these newly created values.
* Do not use `taintUniqueValue` to protect low-entropy values such as PIN codes or phone numbers. If any value in a request is controlled by an attacker, they could infer which value is tainted by enumerating all possible values of the secret.

---

## Usage {/*usage*/}

### Prevent a token from being passed to Client Components {/*prevent-a-token-from-being-passed-to-client-components*/}

To ensure that sensitive information such as passwords, session tokens, or other unique values do not inadvertently get passed to Client Components, the `taintUniqueValue` function provides a layer of protection. When a value is tainted, any attempt to pass it to a Client Component will result in an error.

The `lifetime` argument defines the duration for which the value remains tainted. For values that should remain tainted indefinitely, objects like [`globalThis`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/globalThis) or `process` can serve as the `lifetime` argument. These objects have a lifespan that spans the entire duration of your app's execution.

```js
import {experimental_taintUniqueValue} from 'react';

experimental_taintUniqueValue(
 'Do not pass a user password to the client.',
 globalThis,
 process.env.SECRET_KEY
);
```

If the tainted value's lifespan is tied to a object, the `lifetime` should be the object that encapsulates the value. This ensures the tainted value remains protected for the lifetime of the encapsulating object.

```js
import {experimental_taintUniqueValue} from 'react';

export async function getUser(id) {
 const user = await db`SELECT * FROM users WHERE id = ${id}`;
 experimental_taintUniqueValue(
 'Do not pass a user session token to the client.',
 user,
 user.session.token
 );
 return user;
}
```

In this example, the `user` object serves as the `lifetime` argument. If this object gets stored in a global cache or is accessible by another request, the session token remains tainted.

<Pitfall>

**Do not rely solely on tainting for security.** Tainting a value doesn't block every possible derived value. For example, creating a new value by upper casing a tainted string will not taint the new value.

```js
import {experimental_taintUniqueValue} from 'react';

const password = 'correct horse battery staple';

experimental_taintUniqueValue(
 'Do not pass the password to the client.',
 globalThis,
 password
);

const uppercasePassword = password.toUpperCase() // `uppercasePassword` is not tainted
```

In this example, the constant `password` is tainted. Then `password` is used to create a new value `uppercasePassword` by calling the `toUpperCase` method on `password`. The newly created `uppercasePassword` is not tainted.

Other similar ways of deriving new values from tainted values like concatenating it into a larger string, converting it to base64, or returning a substring create untained values.

Tainting only protects against simple mistakes like explicitly passing secret values to the client. Mistakes in calling the `taintUniqueValue` like using a global store outside of React, without the corresponding lifetime object, can cause the tainted value to become untainted. Tainting is a layer of protection; a secure app will have multiple layers of protection, well designed APIs, and isolation patterns.

</Pitfall>

<DeepDive>

#### Using `server-only` and `taintUniqueValue` to prevent leaking secrets {/*using-server-only-and-taintuniquevalue-to-prevent-leaking-secrets*/}

If you're running a Server Components environment that has access to private keys or passwords such as database passwords, you have to be careful not to pass that to a Client Component.

```js
export async function Dashboard(props) {
 // DO NOT DO THIS
 return <Overview password={process.env.API_PASSWORD} />;
}
```

```js
"use client";

import {useEffect} from '...'

export async function Overview({ password }) {
 useEffect(() => {
 const headers = { Authorization: password };
 fetch(url, { headers }).then(...);
 }, [password]);
 ...
}
```

This example would leak the secret API token to the client. If this API token can be used to access data this particular user shouldn't have access to, it could lead to a data breach.

[comment]: <> (TODO: Link to `server-only` docs once they are written)

Ideally, secrets like this are abstracted into a single helper file that can only be imported by trusted data utilities on the server. The helper can even be tagged with [`server-only`](https://www.npmjs.com/package/server-only) to ensure that this file isn't imported on the client.

```js
import "server-only";

export function fetchAPI(url) {
 const headers = { Authorization: process.env.API_PASSWORD };
 return fetch(url, { headers });
}
```

Sometimes mistakes happen during refactoring and not all of your colleagues might know about this.
To protect against this mistakes happening down the line we can "taint" the actual password:

```js
import "server-only";
import {experimental_taintUniqueValue} from 'react';

experimental_taintUniqueValue(
 'Do not pass the API token password to the client. ' +
 'Instead do all fetches on the server.'
 process,
 process.env.API_PASSWORD
);
```

Now whenever anyone tries to pass this password to a Client Component, or send the password to a Client Component with a Server Function, an error will be thrown with message you defined when you called `taintUniqueValue`.

</DeepDive>

---

---
title: forwardRef
---

<Deprecated>

In React 19, `forwardRef` is no longer necessary. Pass `ref` as a prop instead.

`forwardRef` will be deprecated in a future release. Learn more [here](/blog/2024/04/25/react-19#ref-as-a-prop).

</Deprecated>

<Intro>

`forwardRef` lets your component expose a DOM node to the parent component with a [ref.](/learn/manipulating-the-dom-with-refs)

```js
const SomeComponent = forwardRef(render)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `forwardRef(render)` {/*forwardref*/}

Call `forwardRef()` to let your component receive a ref and forward it to a child component:

```js
import { forwardRef } from 'react';

const MyInput = forwardRef(function MyInput(props, ref) {
 // ...
});
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `render`: The render function for your component. React calls this function with the props and `ref` that your component received from its parent. The JSX you return will be the output of your component.

#### Returns {/*returns*/}

`forwardRef` returns a React component that you can render in JSX. Unlike React components defined as plain functions, a component returned by `forwardRef` is also able to receive a `ref` prop.

#### Caveats {/*caveats*/}

* In Strict Mode, React will **call your render function twice** in order to [help you find accidental impurities.](/reference/react/useState#my-initializer-or-updater-function-runs-twice) This is development-only behavior and does not affect production. If your render function is pure (as it should be), this should not affect the logic of your component. The result from one of the calls will be ignored.

---

### `render` function {/*render-function*/}

`forwardRef` accepts a render function as an argument. React calls this function with `props` and `ref`:

```js
const MyInput = forwardRef(function MyInput(props, ref) {
 return (
 <label>
 {props.label}
 <input ref={ref} />
 </label>
 );
});
```

#### Parameters {/*render-parameters*/}

* `props`: The props passed by the parent component.

* `ref`: The `ref` attribute passed by the parent component. The `ref` can be an object or a function. If the parent component has not passed a ref, it will be `null`. You should either pass the `ref` you receive to another component, or pass it to [`useImperativeHandle`.](/reference/react/useImperativeHandle)

#### Returns {/*render-returns*/}

`forwardRef` returns a React component that you can render in JSX. Unlike React components defined as plain functions, the component returned by `forwardRef` is able to take a `ref` prop.

---

## Usage {/*usage*/}

### Exposing a DOM node to the parent component {/*exposing-a-dom-node-to-the-parent-component*/}

By default, each component's DOM nodes are private. However, sometimes it's useful to expose a DOM node to the parent--for example, to allow focusing it. To opt in, wrap your component definition into `forwardRef()`:

```js {3,11}
import { forwardRef } from 'react';

const MyInput = forwardRef(function MyInput(props, ref) {
 const { label, ...otherProps } = props;
 return (
 <label>
 {label}
 <input {...otherProps} />
 </label>
 );
});
```

You will receive a <CodeStep step={1}>ref</CodeStep> as the second argument after props. Pass it to the DOM node that you want to expose:

```js {8} [[1, 3, "ref"], [1, 8, "ref", 30]]
import { forwardRef } from 'react';

const MyInput = forwardRef(function MyInput(props, ref) {
 const { label, ...otherProps } = props;
 return (
 <label>
 {label}
 <input {...otherProps} ref={ref} />
 </label>
 );
});
```

This lets the parent `Form` component access the <CodeStep step={2}>`<input>` DOM node</CodeStep> exposed by `MyInput`:

```js [[1, 2, "ref"], [1, 10, "ref", 41], [2, 5, "ref.current"]]
function Form() {
 const ref = useRef(null);

 function handleClick() {
 ref.current.focus();
 }

 return (
 <form>
 <MyInput label="Enter your name:" ref={ref} />
 <button type="button" onClick={handleClick}>
 Edit
 </button>
 </form>
 );
}
```

This `Form` component [passes a ref](/reference/react/useRef#manipulating-the-dom-with-a-ref) to `MyInput`. The `MyInput` component *forwards* that ref to the `<input>` browser tag. As a result, the `Form` component can access that `<input>` DOM node and call [`focus()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/focus) on it.

Keep in mind that exposing a ref to the DOM node inside your component makes it harder to change your component's internals later. You will typically expose DOM nodes from reusable low-level components like buttons or text inputs, but you won't do it for application-level components like an avatar or a comment.

<Recipes titleText="Examples of forwarding a ref">

#### Focusing a text input {/*focusing-a-text-input*/}

Clicking the button will focus the input. The `Form` component defines a ref and passes it to the `MyInput` component. The `MyInput` component forwards that ref to the browser `<input>`. This lets the `Form` component focus the `<input>`.

<Sandpack>

```js
import { useRef } from 'react';
import MyInput from './MyInput.js';

export default function Form() {
 const ref = useRef(null);

 function handleClick() {
 ref.current.focus();
 }

 return (
 <form>
 <MyInput label="Enter your name:" ref={ref} />
 <button type="button" onClick={handleClick}>
 Edit
 </button>
 </form>
 );
}
```

```js src/MyInput.js
import { forwardRef } from 'react';

const MyInput = forwardRef(function MyInput(props, ref) {
 const { label, ...otherProps } = props;
 return (
 <label>
 {label}
 <input {...otherProps} ref={ref} />
 </label>
 );
});

export default MyInput;
```

```css
input {
 margin: 5px;
}
```

</Sandpack>

<Solution />

#### Playing and pausing a video {/*playing-and-pausing-a-video*/}

Clicking the button will call [`play()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/play) and [`pause()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/pause) on a `<video>` DOM node. The `App` component defines a ref and passes it to the `MyVideoPlayer` component. The `MyVideoPlayer` component forwards that ref to the browser `<video>` node. This lets the `App` component play and pause the `<video>`.

<Sandpack>

```js
import { useRef } from 'react';
import MyVideoPlayer from './MyVideoPlayer.js';

export default function App() {
 const ref = useRef(null);
 return (
 <>
 <button onClick={() => ref.current.play()}>
 Play
 </button>
 <button onClick={() => ref.current.pause()}>
 Pause
 </button>
 <br />
 <MyVideoPlayer
 ref={ref}
 src="https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4"
 type="video/mp4"
 width="250"
 />
 </>
 );
}
```

```js src/MyVideoPlayer.js
import { forwardRef } from 'react';

const VideoPlayer = forwardRef(function VideoPlayer({ src, type, width }, ref) {
 return (
 <video width={width} ref={ref}>
 <source
 src={src}
 type={type}
 />
 </video>
 );
});

export default VideoPlayer;
```

```css
button { margin-bottom: 10px; margin-right: 10px; }
```

</Sandpack>

<Solution />

</Recipes>

---

### Forwarding a ref through multiple components {/*forwarding-a-ref-through-multiple-components*/}

Instead of forwarding a `ref` to a DOM node, you can forward it to your own component like `MyInput`:

```js {1,5}
const FormField = forwardRef(function FormField(props, ref) {
 // ...
 return (
 <>
 <MyInput ref={ref} />
 ...
 </>
 );
});
```

If that `MyInput` component forwards a ref to its `<input>`, a ref to `FormField` will give you that `<input>`:

```js {2,5,10}
function Form() {
 const ref = useRef(null);

 function handleClick() {
 ref.current.focus();
 }

 return (
 <form>
 <FormField label="Enter your name:" ref={ref} isRequired={true} />
 <button type="button" onClick={handleClick}>
 Edit
 </button>
 </form>
 );
}
```

The `Form` component defines a ref and passes it to `FormField`. The `FormField` component forwards that ref to `MyInput`, which forwards it to a browser `<input>` DOM node. This is how `Form` accesses that DOM node.

<Sandpack>

```js
import { useRef } from 'react';
import FormField from './FormField.js';

export default function Form() {
 const ref = useRef(null);

 function handleClick() {
 ref.current.focus();
 }

 return (
 <form>
 <FormField label="Enter your name:" ref={ref} isRequired={true} />
 <button type="button" onClick={handleClick}>
 Edit
 </button>
 </form>
 );
}
```

```js src/FormField.js
import { forwardRef, useState } from 'react';
import MyInput from './MyInput.js';

const FormField = forwardRef(function FormField({ label, isRequired }, ref) {
 const [value, setValue] = useState('');
 return (
 <>
 <MyInput
 ref={ref}
 label={label}
 value={value}
 onChange={e => setValue(e.target.value)}
 />
 {(isRequired && value === '') &&
 <i>Required</i>
 }
 </>
 );
});

export default FormField;
```

```js src/MyInput.js
import { forwardRef } from 'react';

const MyInput = forwardRef((props, ref) => {
 const { label, ...otherProps } = props;
 return (
 <label>
 {label}
 <input {...otherProps} ref={ref} />
 </label>
 );
});

export default MyInput;
```

```css
input, button {
 margin: 5px;
}
```

</Sandpack>

---

### Exposing an imperative handle instead of a DOM node {/*exposing-an-imperative-handle-instead-of-a-dom-node*/}

Instead of exposing an entire DOM node, you can expose a custom object, called an *imperative handle,* with a more constrained set of methods. To do this, you'd need to define a separate ref to hold the DOM node:

```js {2,6}
const MyInput = forwardRef(function MyInput(props, ref) {
 const inputRef = useRef(null);

 // ...

 return <input {...props} ref={inputRef} />;
});
```

Pass the `ref` you received to [`useImperativeHandle`](/reference/react/useImperativeHandle) and specify the value you want to expose to the `ref`:

```js {6-15}
import { forwardRef, useRef, useImperativeHandle } from 'react';

const MyInput = forwardRef(function MyInput(props, ref) {
 const inputRef = useRef(null);

 useImperativeHandle(ref, () => {
 return {
 focus() {
 inputRef.current.focus();
 },
 scrollIntoView() {
 inputRef.current.scrollIntoView();
 },
 };
 }, []);

 return <input {...props} ref={inputRef} />;
});
```

If some component gets a ref to `MyInput`, it will only receive your `{ focus, scrollIntoView }` object instead of the DOM node. This lets you limit the information you expose about your DOM node to the minimum.

<Sandpack>

```js
import { useRef } from 'react';
import MyInput from './MyInput.js';

export default function Form() {
 const ref = useRef(null);

 function handleClick() {
 ref.current.focus();
 // This won't work because the DOM node isn't exposed:
 // ref.current.style.opacity = 0.5;
 }

 return (
 <form>
 <MyInput placeholder="Enter your name" ref={ref} />
 <button type="button" onClick={handleClick}>
 Edit
 </button>
 </form>
 );
}
```

```js src/MyInput.js
import { forwardRef, useRef, useImperativeHandle } from 'react';

const MyInput = forwardRef(function MyInput(props, ref) {
 const inputRef = useRef(null);

 useImperativeHandle(ref, () => {
 return {
 focus() {
 inputRef.current.focus();
 },
 scrollIntoView() {
 inputRef.current.scrollIntoView();
 },
 };
 }, []);

 return <input {...props} ref={inputRef} />;
});

export default MyInput;
```

```css
input {
 margin: 5px;
}
```

</Sandpack>

[Read more about using imperative handles.](/reference/react/useImperativeHandle)

<Pitfall>

**Do not overuse refs.** You should only use refs for *imperative* behaviors that you can't express as props: for example, scrolling to a node, focusing a node, triggering an animation, selecting text, and so on.

**If you can express something as a prop, you should not use a ref.** For example, instead of exposing an imperative handle like `{ open, close }` from a `Modal` component, it is better to take `isOpen` as a prop like `<Modal isOpen={isOpen} />`. [Effects](/learn/synchronizing-with-effects) can help you expose imperative behaviors via props.

</Pitfall>

---

## Troubleshooting {/*troubleshooting*/}

### My component is wrapped in `forwardRef`, but the `ref` to it is always `null` {/*my-component-is-wrapped-in-forwardref-but-the-ref-to-it-is-always-null*/}

This usually means that you forgot to actually use the `ref` that you received.

For example, this component doesn't do anything with its `ref`:

```js {1}
const MyInput = forwardRef(function MyInput({ label }, ref) {
 return (
 <label>
 {label}
 <input />
 </label>
 );
});
```

To fix it, pass the `ref` down to a DOM node or another component that can accept a ref:

```js {1,5}
const MyInput = forwardRef(function MyInput({ label }, ref) {
 return (
 <label>
 {label}
 <input ref={ref} />
 </label>
 );
});
```

The `ref` to `MyInput` could also be `null` if some of the logic is conditional:

```js {1,5}
const MyInput = forwardRef(function MyInput({ label, showInput }, ref) {
 return (
 <label>
 {label}
 {showInput && <input ref={ref} />}
 </label>
 );
});
```

If `showInput` is `false`, then the ref won't be forwarded to any node, and a ref to `MyInput` will remain empty. This is particularly easy to miss if the condition is hidden inside another component, like `Panel` in this example:

```js {5,7}
const MyInput = forwardRef(function MyInput({ label, showInput }, ref) {
 return (
 <label>
 {label}
 <Panel isExpanded={showInput}>
 <input ref={ref} />
 </Panel>
 </label>
 );
});
```

---
title: "Built-in React Hooks"
---

<Intro>

*Hooks* let you use different React features from your components. You can either use the built-in Hooks or combine them to build your own. This page lists all built-in Hooks in React.

</Intro>

---

## State Hooks {/*state-hooks*/}

*State* lets a component ["remember" information like user input.](/learn/state-a-components-memory) For example, a form component can use state to store the input value, while an image gallery component can use state to store the selected image index.

To add state to a component, use one of these Hooks:

* [`useState`](/reference/react/useState) declares a state variable that you can update directly.
* [`useReducer`](/reference/react/useReducer) declares a state variable with the update logic inside a [reducer function.](/learn/extracting-state-logic-into-a-reducer)

```js
function ImageGallery() {
 const [index, setIndex] = useState(0);
 // ...
```

---

## Context Hooks {/*context-hooks*/}

*Context* lets a component [receive information from distant parents without passing it as props.](/learn/passing-props-to-a-component) For example, your app's top-level component can pass the current UI theme to all components below, no matter how deep.

* [`useContext`](/reference/react/useContext) reads and subscribes to a context.

```js
function Button() {
 const theme = useContext(ThemeContext);
 // ...
```

---

## Ref Hooks {/*ref-hooks*/}

*Refs* let a component [hold some information that isn't used for rendering,](/learn/referencing-values-with-refs) like a DOM node or a timeout ID. Unlike with state, updating a ref does not re-render your component. Refs are an "escape hatch" from the React paradigm. They are useful when you need to work with non-React systems, such as the built-in browser APIs.

* [`useRef`](/reference/react/useRef) declares a ref. You can hold any value in it, but most often it's used to hold a DOM node.
* [`useImperativeHandle`](/reference/react/useImperativeHandle) lets you customize the ref exposed by your component. This is rarely used.

```js
function Form() {
 const inputRef = useRef(null);
 // ...
```

---

## Effect Hooks {/*effect-hooks*/}

*Effects* let a component [connect to and synchronize with external systems.](/learn/synchronizing-with-effects) This includes dealing with network, browser DOM, animations, widgets written using a different UI library, and other non-React code.

* [`useEffect`](/reference/react/useEffect) connects a component to an external system.

```js
function ChatRoom({ roomId }) {
 useEffect(() => {
 const connection = createConnection(roomId);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]);
 // ...
```

Effects are an "escape hatch" from the React paradigm. Don't use Effects to orchestrate the data flow of your application. If you're not interacting with an external system, [you might not need an Effect.](/learn/you-might-not-need-an-effect)

There are two rarely used variations of `useEffect` with differences in timing:

* [`useLayoutEffect`](/reference/react/useLayoutEffect) fires before the browser repaints the screen. You can measure layout here.
* [`useInsertionEffect`](/reference/react/useInsertionEffect) fires before React makes changes to the DOM. Libraries can insert dynamic CSS here.

You can also separate events from Effects:

- [`useEffectEvent`](/reference/react/useEffectEvent) creates a non-reactive event to fire from any Effect hook.
---

## Performance Hooks {/*performance-hooks*/}

A common way to optimize re-rendering performance is to skip unnecessary work. For example, you can tell React to reuse a cached calculation or to skip a re-render if the data has not changed since the previous render.

To skip calculations and unnecessary re-rendering, use one of these Hooks:

- [`useMemo`](/reference/react/useMemo) lets you cache the result of an expensive calculation.
- [`useCallback`](/reference/react/useCallback) lets you cache a function definition before passing it down to an optimized component.

```js
function TodoList({ todos, tab, theme }) {
 const visibleTodos = useMemo(() => filterTodos(todos, tab), [todos, tab]);
 // ...
}
```

Sometimes, you can't skip re-rendering because the screen actually needs to update. In that case, you can improve performance by separating blocking updates that must be synchronous (like typing into an input) from non-blocking updates which don't need to block the user interface (like updating a chart).

To prioritize rendering, use one of these Hooks:

- [`useTransition`](/reference/react/useTransition) lets you mark a state transition as non-blocking and allow other updates to interrupt it.
- [`useDeferredValue`](/reference/react/useDeferredValue) lets you defer updating a non-critical part of the UI and let other parts update first.

---

## Other Hooks {/*other-hooks*/}

These Hooks are mostly useful to library authors and aren't commonly used in the application code.

- [`useDebugValue`](/reference/react/useDebugValue) lets you customize the label React DevTools displays for your custom Hook.
- [`useId`](/reference/react/useId) lets a component associate a unique ID with itself. Typically used with accessibility APIs.
- [`useSyncExternalStore`](/reference/react/useSyncExternalStore) lets a component subscribe to an external store.
* [`useActionState`](/reference/react/useActionState) allows you to manage state of actions.

---

## Your own Hooks {/*your-own-hooks*/}

You can also [define your own custom Hooks](/learn/reusing-logic-with-custom-hooks#extracting-your-own-custom-hook-from-a-component) as JavaScript functions.

---
title: React Reference Overview
---

<Intro>

This section provides detailed reference documentation for working with React. For an introduction to React, please visit the [Learn](/learn) section.

</Intro>

The React reference documentation is broken down into functional subsections:

## React {/*react*/}

Programmatic React features:

* [Hooks](/reference/react/hooks) - Use different React features from your components.
* [Components](/reference/react/components) - Built-in components that you can use in your JSX.
* [APIs](/reference/react/apis) - APIs that are useful for defining components.
* [Directives](/reference/rsc/directives) - Provide instructions to bundlers compatible with React Server Components.

## React DOM {/*react-dom*/}

React DOM contains features that are only supported for web applications (which run in the browser DOM environment). This section is broken into the following:

* [Hooks](/reference/react-dom/hooks) - Hooks for web applications which run in the browser DOM environment.
* [Components](/reference/react-dom/components) - React supports all of the browser built-in HTML and SVG components.
* [APIs](/reference/react-dom) - The `react-dom` package contains methods supported only in web applications.
* [Client APIs](/reference/react-dom/client) - The `react-dom/client` APIs let you render React components on the client (in the browser).
* [Server APIs](/reference/react-dom/server) - The `react-dom/server` APIs let you render React components to HTML on the server.
* [Static APIs](/reference/react-dom/static) - The `react-dom/static` APIs let you generate static HTML for React components.

## React Compiler {/*react-compiler*/}

The React Compiler is a build-time optimization tool that automatically memoizes your React components and values:

* [Configuration](/reference/react-compiler/configuration) - Configuration options for React Compiler.
* [Directives](/reference/react-compiler/directives) - Function-level directives to control compilation.
* [Compiling Libraries](/reference/react-compiler/compiling-libraries) - Guide for shipping pre-compiled library code.

## ESLint Plugin React Hooks {/*eslint-plugin-react-hooks*/}

The [ESLint plugin for React Hooks](/reference/eslint-plugin-react-hooks) helps enforce the Rules of React:

* [Lints](/reference/eslint-plugin-react-hooks) - Detailed documentation for each lint with examples.

## Rules of React {/*rules-of-react*/}

React has idioms — or rules — for how to express patterns in a way that is easy to understand and yields high-quality applications:

* [Components and Hooks must be pure](/reference/rules/components-and-hooks-must-be-pure) – Purity makes your code easier to understand, debug, and allows React to automatically optimize your components and hooks correctly.
* [React calls Components and Hooks](/reference/rules/react-calls-components-and-hooks) – React is responsible for rendering components and hooks when necessary to optimize the user experience.
* [Rules of Hooks](/reference/rules/rules-of-hooks) – Hooks are defined using JavaScript functions, but they represent a special type of reusable UI logic with restrictions on where they can be called.

## Legacy APIs {/*legacy-apis*/}

* [Legacy APIs](/reference/react/legacy) - Exported from the `react` package, but not recommended for use in newly written code.

---
title: isValidElement
---

<Intro>

`isValidElement` checks whether a value is a React element.

```js
const isElement = isValidElement(value)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `isValidElement(value)` {/*isvalidelement*/}

Call `isValidElement(value)` to check whether `value` is a React element.

```js
import { isValidElement, createElement } from 'react';

// ✅ React elements
console.log(isValidElement(<p />)); // true
console.log(isValidElement(createElement('p'))); // true

// ❌ Not React elements
console.log(isValidElement(25)); // false
console.log(isValidElement('Hello')); // false
console.log(isValidElement({ age: 42 })); // false
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `value`: The `value` you want to check. It can be any a value of any type.

#### Returns {/*returns*/}

`isValidElement` returns `true` if the `value` is a React element. Otherwise, it returns `false`.

#### Caveats {/*caveats*/}

* **Only [JSX tags](/learn/writing-markup-with-jsx) and objects returned by [`createElement`](/reference/react/createElement) are considered to be React elements.** For example, even though a number like `42` is a valid React *node* (and can be returned from a component), it is not a valid React element. Arrays and portals created with [`createPortal`](/reference/react-dom/createPortal) are also *not* considered to be React elements.

---

## Usage {/*usage*/}

### Checking if something is a React element {/*checking-if-something-is-a-react-element*/}

Call `isValidElement` to check if some value is a *React element.*

React elements are:

- Values produced by writing a [JSX tag](/learn/writing-markup-with-jsx)
- Values produced by calling [`createElement`](/reference/react/createElement)

For React elements, `isValidElement` returns `true`:

```js
import { isValidElement, createElement } from 'react';

// ✅ JSX tags are React elements
console.log(isValidElement(<p />)); // true
console.log(isValidElement(<MyComponent />)); // true

// ✅ Values returned by createElement are React elements
console.log(isValidElement(createElement('p'))); // true
console.log(isValidElement(createElement(MyComponent))); // true
```

Any other values, such as strings, numbers, or arbitrary objects and arrays, are not React elements.

For them, `isValidElement` returns `false`:

```js
// ❌ These are *not* React elements
console.log(isValidElement(null)); // false
console.log(isValidElement(25)); // false
console.log(isValidElement('Hello')); // false
console.log(isValidElement({ age: 42 })); // false
console.log(isValidElement([<div />, <div />])); // false
console.log(isValidElement(MyComponent)); // false
```

It is very uncommon to need `isValidElement`. It's mostly useful if you're calling another API that *only* accepts elements (like [`cloneElement`](/reference/react/cloneElement) does) and you want to avoid an error when your argument is not a React element.

Unless you have some very specific reason to add an `isValidElement` check, you probably don't need it.

<DeepDive>

#### React elements vs React nodes {/*react-elements-vs-react-nodes*/}

When you write a component, you can return any kind of *React node* from it:

```js
function MyComponent() {
 // ... you can return any React node ...
}
```

A React node can be:

- A React element created like `<div />` or `createElement('div')`
- A portal created with [`createPortal`](/reference/react-dom/createPortal)
- A string
- A number
- `true`, `false`, `null`, or `undefined` (which are not displayed)
- An array of other React nodes

**Note `isValidElement` checks whether the argument is a *React element,* not whether it's a React node.** For example, `42` is not a valid React element. However, it is a perfectly valid React node:

```js
function MyComponent() {
 return 42; // It's ok to return a number from component
}
```

This is why you shouldn't use `isValidElement` as a way to check whether something can be rendered.

</DeepDive>

---
title: lazy
---

<Intro>

`lazy` lets you defer loading component's code until it is rendered for the first time.

```js
const SomeComponent = lazy(load)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `lazy(load)` {/*lazy*/}

Call `lazy` outside your components to declare a lazy-loaded React component:

```js
import { lazy } from 'react';

const MarkdownPreview = lazy(() => import('./MarkdownPreview.js'));
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `load`: A function that returns a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) or another *thenable* (a Promise-like object with a `then` method). React will not call `load` until the first time you attempt to render the returned component. After React first calls `load`, it will wait for it to resolve, and then render the resolved value's `.default` as a React component. Both the returned Promise and the Promise's resolved value will be cached, so React will not call `load` more than once. If the Promise rejects, React will `throw` the rejection reason for the nearest Error Boundary to handle.

#### Returns {/*returns*/}

`lazy` returns a React component you can render in your tree. While the code for the lazy component is still loading, attempting to render it will *suspend.* Use [`<Suspense>`](/reference/react/Suspense) to display a loading indicator while it's loading.

---

### `load` function {/*load*/}

#### Parameters {/*load-parameters*/}

`load` receives no parameters.

#### Returns {/*load-returns*/}

You need to return a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) or some other *thenable* (a Promise-like object with a `then` method). It needs to eventually resolve to an object whose `.default` property is a valid React component type, such as a function, [`memo`](/reference/react/memo), or a [`forwardRef`](/reference/react/forwardRef) component.

---

## Usage {/*usage*/}

### Lazy-loading components with Suspense {/*suspense-for-code-splitting*/}

Usually, you import components with the static [`import`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/import) declaration:

```js
import MarkdownPreview from './MarkdownPreview.js';
```

To defer loading this component's code until it's rendered for the first time, replace this import with:

```js
import { lazy } from 'react';

const MarkdownPreview = lazy(() => import('./MarkdownPreview.js'));
```

This code relies on [dynamic `import()`,](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/import) which might require support from your bundler or framework. Using this pattern requires that the lazy component you're importing was exported as the `default` export.

Now that your component's code loads on demand, you also need to specify what should be displayed while it is loading. You can do this by wrapping the lazy component or any of its parents into a [`<Suspense>`](/reference/react/Suspense) boundary:

```js {1,4}
<Suspense fallback={<Loading />}>
 <h2>Preview</h2>
 <MarkdownPreview />
</Suspense>
```

In this example, the code for `MarkdownPreview` won't be loaded until you attempt to render it. If `MarkdownPreview` hasn't loaded yet, `Loading` will be shown in its place. Try ticking the checkbox:

<Sandpack>

```js src/App.js
import { useState, Suspense, lazy } from 'react';
import Loading from './Loading.js';

const MarkdownPreview = lazy(() => delayForDemo(import('./MarkdownPreview.js')));

export default function MarkdownEditor() {
 const [showPreview, setShowPreview] = useState(false);
 const [markdown, setMarkdown] = useState('Hello, **world**!');
 return (
 <>
 <textarea value={markdown} onChange={e => setMarkdown(e.target.value)} />
 <label>
 <input type="checkbox" checked={showPreview} onChange={e => setShowPreview(e.target.checked)} />
 Show preview
 </label>
 <hr />
 {showPreview && (
 <Suspense fallback={<Loading />}>
 <h2>Preview</h2>
 <MarkdownPreview markdown={markdown} />
 </Suspense>
 )}
 </>
 );
}

// Add a fixed delay so you can see the loading state
function delayForDemo(promise) {
 return new Promise(resolve => {
 setTimeout(resolve, 2000);
 }).then(() => promise);
}
```

```js src/Loading.js
export default function Loading() {
 return <p><i>Loading...</i></p>;
}
```

```js src/MarkdownPreview.js
import { Remarkable } from 'remarkable';

const md = new Remarkable();

export default function MarkdownPreview({ markdown }) {
 return (
 <div
 className="content"
 dangerouslySetInnerHTML={{__html: md.render(markdown)}}
 />
 );
}
```

```json package.json hidden
{
 "dependencies": {
 "immer": "1.7.3",
 "react": "latest",
 "react-dom": "latest",
 "react-scripts": "latest",
 "remarkable": "2.0.1"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```css
label {
 display: block;
}

input, textarea {
 margin-bottom: 10px;
}

body {
 min-height: 200px;
}
```

</Sandpack>

This demo loads with an artificial delay. The next time you untick and tick the checkbox, `Preview` will be cached, so there will be no loading state. To see the loading state again, click "Reset" on the sandbox.

[Learn more about managing loading states with Suspense.](/reference/react/Suspense)

---

## Troubleshooting {/*troubleshooting*/}

### My `lazy` component's state gets reset unexpectedly {/*my-lazy-components-state-gets-reset-unexpectedly*/}

Do not declare `lazy` components *inside* other components:

```js {4-5}
import { lazy } from 'react';

function Editor() {
 // 🔴 Bad: This will cause all state to be reset on re-renders
 const MarkdownPreview = lazy(() => import('./MarkdownPreview.js'));
 // ...
}
```

Instead, always declare them at the top level of your module:

```js {3-4}
import { lazy } from 'react';

// ✅ Good: Declare lazy components outside of your components
const MarkdownPreview = lazy(() => import('./MarkdownPreview.js'));

function Editor() {
 // ...
}
```

---
title: "Legacy React APIs"
---

<Intro>

These APIs are exported from the `react` package, but they are not recommended for use in newly written code. See the linked individual API pages for the suggested alternatives.

</Intro>

---

## Legacy APIs {/*legacy-apis*/}

* [`Children`](/reference/react/Children) lets you manipulate and transform the JSX received as the `children` prop. [See alternatives.](/reference/react/Children#alternatives)
* [`cloneElement`](/reference/react/cloneElement) lets you create a React element using another element as a starting point. [See alternatives.](/reference/react/cloneElement#alternatives)
* [`Component`](/reference/react/Component) lets you define a React component as a JavaScript class. [See alternatives.](/reference/react/Component#alternatives)
* [`createElement`](/reference/react/createElement) lets you create a React element. Typically, you'll use JSX instead.
* [`createRef`](/reference/react/createRef) creates a ref object which can contain arbitrary value. [See alternatives.](/reference/react/createRef#alternatives)
* [`forwardRef`](/reference/react/forwardRef) lets your component expose a DOM node to parent component with a [ref.](/learn/manipulating-the-dom-with-refs)
* [`isValidElement`](/reference/react/isValidElement) checks whether a value is a React element. Typically used with [`cloneElement`.](/reference/react/cloneElement)
* [`PureComponent`](/reference/react/PureComponent) is similar to [`Component`,](/reference/react/Component) but it skip re-renders with same props. [See alternatives.](/reference/react/PureComponent#alternatives)

---

## Removed APIs {/*removed-apis*/}

These APIs were removed in React 19:

* [`createFactory`](https://18.react.dev/reference/react/createFactory): use JSX instead.
* Class Components: [`static contextTypes`](https://18.react.dev//reference/react/Component#static-contexttypes): use [`static contextType`](#static-contexttype) instead.
* Class Components: [`static childContextTypes`](https://18.react.dev//reference/react/Component#static-childcontexttypes): use [`static contextType`](#static-contexttype) instead.
* Class Components: [`static getChildContext`](https://18.react.dev//reference/react/Component#getchildcontext): use [`Context`](/reference/react/createContext#provider) instead.
* Class Components: [`static propTypes`](https://18.react.dev//reference/react/Component#static-proptypes): use a type system like [TypeScript](https://www.typescriptlang.org/) instead.
* Class Components: [`this.refs`](https://18.react.dev//reference/react/Component#refs): use [`createRef`](/reference/react/createRef) instead.

---
title: memo
---

<Intro>

`memo` lets you skip re-rendering a component when its props are unchanged.

```
const MemoizedComponent = memo(SomeComponent, arePropsEqual?)
```

</Intro>

<Note>

[React Compiler](/learn/react-compiler) automatically applies the equivalent of `memo` to all components, reducing the need for manual memoization. You can use the compiler to handle component memoization automatically.

</Note>

<InlineToc />

---

## Reference {/*reference*/}

### `memo(Component, arePropsEqual?)` {/*memo*/}

Wrap a component in `memo` to get a *memoized* version of that component. This memoized version of your component will usually not be re-rendered when its parent component is re-rendered as long as its props have not changed. But React may still re-render it: memoization is a performance optimization, not a guarantee.

```js
import { memo } from 'react';

const SomeComponent = memo(function SomeComponent(props) {
 // ...
});
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `Component`: The component that you want to memoize. The `memo` does not modify this component, but returns a new, memoized component instead. Any valid React component, including functions and [`forwardRef`](/reference/react/forwardRef) components, is accepted.

* **optional** `arePropsEqual`: A function that accepts two arguments: the component's previous props, and its new props. It should return `true` if the old and new props are equal: that is, if the component will render the same output and behave in the same way with the new props as with the old. Otherwise it should return `false`. Usually, you will not specify this function. By default, React will compare each prop with [`Object.is`.](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is)

#### Returns {/*returns*/}

`memo` returns a new React component. It behaves the same as the component provided to `memo` except that React will not always re-render it when its parent is being re-rendered unless its props have changed.

---

## Usage {/*usage*/}

### Skipping re-rendering when props are unchanged {/*skipping-re-rendering-when-props-are-unchanged*/}

React normally re-renders a component whenever its parent re-renders. With `memo`, you can create a component that React will not re-render when its parent re-renders so long as its new props are the same as the old props. Such a component is said to be *memoized*.

To memoize a component, wrap it in `memo` and use the value that it returns in place of your original component:

```js
const Greeting = memo(function Greeting({ name }) {
 return <h1>Hello, {name}!</h1>;
});

export default Greeting;
```

A React component should always have [pure rendering logic.](/learn/keeping-components-pure) This means that it must return the same output if its props, state, and context haven't changed. By using `memo`, you are telling React that your component complies with this requirement, so React doesn't need to re-render as long as its props haven't changed. Even with `memo`, your component will re-render if its own state changes or if a context that it's using changes.

In this example, notice that the `Greeting` component re-renders whenever `name` is changed (because that's one of its props), but not when `address` is changed (because it's not passed to `Greeting` as a prop):

<Sandpack>

```js
import { memo, useState } from 'react';

export default function MyApp() {
 const [name, setName] = useState('');
 const [address, setAddress] = useState('');
 return (
 <>
 <label>
 Name{': '}
 <input value={name} onChange={e => setName(e.target.value)} />
 </label>
 <label>
 Address{': '}
 <input value={address} onChange={e => setAddress(e.target.value)} />
 </label>
 <Greeting name={name} />
 </>
 );
}

const Greeting = memo(function Greeting({ name }) {
 console.log("Greeting was rendered at", new Date().toLocaleTimeString());
 return <h3>Hello{name && ', '}{name}!</h3>;
});
```

```css
label {
 display: block;
 margin-bottom: 16px;
}
```

</Sandpack>

<Note>

**You should only rely on `memo` as a performance optimization.** If your code doesn't work without it, find the underlying problem and fix it first. Then you may add `memo` to improve performance.

</Note>

<DeepDive>

#### Should you add memo everywhere? {/*should-you-add-memo-everywhere*/}

If your app is like this site, and most interactions are coarse (like replacing a page or an entire section), memoization is usually unnecessary. On the other hand, if your app is more like a drawing editor, and most interactions are granular (like moving shapes), then you might find memoization very helpful.

Optimizing with `memo` is only valuable when your component re-renders often with the same exact props, and its re-rendering logic is expensive. If there is no perceptible lag when your component re-renders, `memo` is unnecessary. Keep in mind that `memo` is completely useless if the props passed to your component are *always different,* such as if you pass an object or a plain function defined during rendering. This is why you will often need [`useMemo`](/reference/react/useMemo#skipping-re-rendering-of-components) and [`useCallback`](/reference/react/useCallback#skipping-re-rendering-of-components) together with `memo`.

There is no benefit to wrapping a component in `memo` in other cases. There is no significant harm to doing that either, so some teams choose to not think about individual cases, and memoize as much as possible. The downside of this approach is that code becomes less readable. Also, not all memoization is effective: a single value that's "always new" is enough to break memoization for an entire component.

**In practice, you can make a lot of memoization unnecessary by following a few principles:**

1. When a component visually wraps other components, let it [accept JSX as children.](/learn/passing-props-to-a-component#passing-jsx-as-children) This way, when the wrapper component updates its own state, React knows that its children don't need to re-render.
1. Prefer local state and don't [lift state up](/learn/sharing-state-between-components) any further than necessary. For example, don't keep transient state like forms and whether an item is hovered at the top of your tree or in a global state library.
1. Keep your [rendering logic pure.](/learn/keeping-components-pure) If re-rendering a component causes a problem or produces some noticeable visual artifact, it's a bug in your component! Fix the bug instead of adding memoization.
1. Avoid [unnecessary Effects that update state.](/learn/you-might-not-need-an-effect) Most performance problems in React apps are caused by chains of updates originating from Effects that cause your components to render over and over.
1. Try to [remove unnecessary dependencies from your Effects.](/learn/removing-effect-dependencies) For example, instead of memoization, it's often simpler to move some object or a function inside an Effect or outside the component.

If a specific interaction still feels laggy, [use the React Developer Tools profiler](https://legacy.reactjs.org/blog/2018/09/10/introducing-the-react-profiler.html) to see which components would benefit the most from memoization, and add memoization where needed. These principles make your components easier to debug and understand, so it's good to follow them in any case. In the long term, we're researching [doing granular memoization automatically](https://www.youtube.com/watch?v=lGEMwh32soc) to solve this once and for all.

</DeepDive>

---

### Updating a memoized component using state {/*updating-a-memoized-component-using-state*/}

Even when a component is memoized, it will still re-render when its own state changes. Memoization only has to do with props that are passed to the component from its parent.

<Sandpack>

```js
import { memo, useState } from 'react';

export default function MyApp() {
 const [name, setName] = useState('');
 const [address, setAddress] = useState('');
 return (
 <>
 <label>
 Name{': '}
 <input value={name} onChange={e => setName(e.target.value)} />
 </label>
 <label>
 Address{': '}
 <input value={address} onChange={e => setAddress(e.target.value)} />
 </label>
 <Greeting name={name} />
 </>
 );
}

const Greeting = memo(function Greeting({ name }) {
 console.log('Greeting was rendered at', new Date().toLocaleTimeString());
 const [greeting, setGreeting] = useState('Hello');
 return (
 <>
 <h3>{greeting}{name && ', '}{name}!</h3>
 <GreetingSelector value={greeting} onChange={setGreeting} />
 </>
 );
});

function GreetingSelector({ value, onChange }) {
 return (
 <>
 <label>
 <input
 type="radio"
 checked={value === 'Hello'}
 onChange={e => onChange('Hello')}
 />
 Regular greeting
 </label>
 <label>
 <input
 type="radio"
 checked={value === 'Hello and welcome'}
 onChange={e => onChange('Hello and welcome')}
 />
 Enthusiastic greeting
 </label>
 </>
 );
}
```

```css
label {
 display: block;
 margin-bottom: 16px;
}
```

</Sandpack>

If you set a state variable to its current value, React will skip re-rendering your component even without `memo`. You may still see your component function being called an extra time, but the result will be discarded.

---

### Updating a memoized component using a context {/*updating-a-memoized-component-using-a-context*/}

Even when a component is memoized, it will still re-render when a context that it's using changes. Memoization only has to do with props that are passed to the component from its parent.

<Sandpack>

```js
import { createContext, memo, useContext, useState } from 'react';

const ThemeContext = createContext(null);

export default function MyApp() {
 const [theme, setTheme] = useState('dark');

 function handleClick() {
 setTheme(theme === 'dark' ? 'light' : 'dark');
 }

 return (
 <ThemeContext value={theme}>
 <button onClick={handleClick}>
 Switch theme
 </button>
 <Greeting name="Taylor" />
 </ThemeContext>
 );
}

const Greeting = memo(function Greeting({ name }) {
 console.log("Greeting was rendered at", new Date().toLocaleTimeString());
 const theme = useContext(ThemeContext);
 return (
 <h3 className={theme}>Hello, {name}!</h3>
 );
});
```

```css
label {
 display: block;
 margin-bottom: 16px;
}

.light {
 color: black;
 background-color: white;
}

.dark {
 color: white;
 background-color: black;
}
```

</Sandpack>

To make your component re-render only when a _part_ of some context changes, split your component in two. Read what you need from the context in the outer component, and pass it down to a memoized child as a prop.

---

### Minimizing props changes {/*minimizing-props-changes*/}

When you use `memo`, your component re-renders whenever any prop is not *shallowly equal* to what it was previously. This means that React compares every prop in your component with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. Note that `Object.is(3, 3)` is `true`, but `Object.is({}, {})` is `false`.

To get the most out of `memo`, minimize the times that the props change. For example, if the prop is an object, prevent the parent component from re-creating that object every time by using [`useMemo`:](/reference/react/useMemo)

```js {5-8}
function Page() {
 const [name, setName] = useState('Taylor');
 const [age, setAge] = useState(42);

 const person = useMemo(
 () => ({ name, age }),
 [name, age]
 );

 return <Profile person={person} />;
}

const Profile = memo(function Profile({ person }) {
 // ...
});
```

A better way to minimize props changes is to make sure the component accepts the minimum necessary information in its props. For example, it could accept individual values instead of a whole object:

```js {4,7}
function Page() {
 const [name, setName] = useState('Taylor');
 const [age, setAge] = useState(42);
 return <Profile name={name} age={age} />;
}

const Profile = memo(function Profile({ name, age }) {
 // ...
});
```

Even individual values can sometimes be projected to ones that change less frequently. For example, here a component accepts a boolean indicating the presence of a value rather than the value itself:

```js {3}
function GroupsLanding({ person }) {
 const hasGroups = person.groups !== null;
 return <CallToAction hasGroups={hasGroups} />;
}

const CallToAction = memo(function CallToAction({ hasGroups }) {
 // ...
});
```

When you need to pass a function to memoized component, either declare it outside your component so that it never changes, or [`useCallback`](/reference/react/useCallback#skipping-re-rendering-of-components) to cache its definition between re-renders.

---

### Specifying a custom comparison function {/*specifying-a-custom-comparison-function*/}

In rare cases it may be infeasible to minimize the props changes of a memoized component. In that case, you can provide a custom comparison function, which React will use to compare the old and new props instead of using shallow equality. This function is passed as a second argument to `memo`. It should return `true` only if the new props would result in the same output as the old props; otherwise it should return `false`.

```js {3}
const Chart = memo(function Chart({ dataPoints }) {
 // ...
}, arePropsEqual);

function arePropsEqual(oldProps, newProps) {
 return (
 oldProps.dataPoints.length === newProps.dataPoints.length &&
 oldProps.dataPoints.every((oldPoint, index) => {
 const newPoint = newProps.dataPoints[index];
 return oldPoint.x === newPoint.x && oldPoint.y === newPoint.y;
 })
 );
}
```

If you do this, use the Performance panel in your browser developer tools to make sure that your comparison function is actually faster than re-rendering the component. You might be surprised.

When you do performance measurements, make sure that React is running in the production mode.

<Pitfall>

If you provide a custom `arePropsEqual` implementation, **you must compare every prop, including functions.** Functions often [close over](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Closures) the props and state of parent components. If you return `true` when `oldProps.onClick !== newProps.onClick`, your component will keep "seeing" the props and state from a previous render inside its `onClick` handler, leading to very confusing bugs.

Avoid doing deep equality checks inside `arePropsEqual` unless you are 100% sure that the data structure you're working with has a known limited depth. **Deep equality checks can become incredibly slow** and can freeze your app for many seconds if someone changes the data structure later.

</Pitfall>

---

### Do I still need React.memo if I use React Compiler? {/*react-compiler-memo*/}

When you enable [React Compiler](/learn/react-compiler), you typically don't need `React.memo` anymore. The compiler automatically optimizes component re-rendering for you.

Here's how it works:

**Without React Compiler**, you need `React.memo` to prevent unnecessary re-renders:

```js
// Parent re-renders every second
function Parent() {
 const [seconds, setSeconds] = useState(0);

 useEffect(() => {
 const interval = setInterval(() => {
 setSeconds(s => s + 1);
 }, 1000);
 return () => clearInterval(interval);
 }, []);

 return (
 <>
 <h1>Seconds: {seconds}</h1>
 <ExpensiveChild name="John" />
 </>
 );
}

// Without memo, this re-renders every second even though props don't change
const ExpensiveChild = memo(function ExpensiveChild({ name }) {
 console.log('ExpensiveChild rendered');
 return <div>Hello, {name}!</div>;
});
```

**With React Compiler enabled**, the same optimization happens automatically:

```js
// No memo needed - compiler prevents re-renders automatically
function ExpensiveChild({ name }) {
 console.log('ExpensiveChild rendered');
 return <div>Hello, {name}!</div>;
}
```

Here's the key part of what the React Compiler generates:

```js {6-12}
function Parent() {
 const $ = _c(7);
 const [seconds, setSeconds] = useState(0);
 // ... other code ...

 let t3;
 if ($[4] === Symbol.for("react.memo_cache_sentinel")) {
 t3 = <ExpensiveChild name="John" />;
 $[4] = t3;
 } else {
 t3 = $[4];
 }
 // ... return statement ...
}
```

Notice the highlighted lines: The compiler wraps `<ExpensiveChild name="John" />` in a cache check. Since the `name` prop is always `"John"`, this JSX is created once and reused on every parent re-render. This is exactly what `React.memo` does - it prevents the child from re-rendering when its props haven't changed.

The React Compiler automatically:
1. Tracks that the `name` prop passed to `ExpensiveChild` hasn't changed
2. Reuses the previously created JSX for `<ExpensiveChild name="John" />`
3. Skips re-rendering `ExpensiveChild` entirely

This means **you can safely remove `React.memo` from your components when using React Compiler**. The compiler provides the same optimization automatically, making your code cleaner and easier to maintain.

<Note>

The compiler's optimization is actually more comprehensive than `React.memo`. It also memoizes intermediate values and expensive computations within your components, similar to combining `React.memo` with `useMemo` throughout your component tree.

</Note>

---

## Troubleshooting {/*troubleshooting*/}
### My component re-renders when a prop is an object, array, or function {/*my-component-rerenders-when-a-prop-is-an-object-or-array*/}

React compares old and new props by shallow equality: that is, it considers whether each new prop is reference-equal to the old prop. If you create a new object or array each time the parent is re-rendered, even if the individual elements are each the same, React will still consider it to be changed. Similarly, if you create a new function when rendering the parent component, React will consider it to have changed even if the function has the same definition. To avoid this, [simplify props or memoize props in the parent component](#minimizing-props-changes).

---
title: startTransition
---

<Intro>

`startTransition` lets you render a part of the UI in the background.

```js
startTransition(action)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `startTransition(action)` {/*starttransition*/}

The `startTransition` function lets you mark a state update as a Transition.

```js {7,9}
import { startTransition } from 'react';

function TabContainer() {
 const [tab, setTab] = useState('about');

 function selectTab(nextTab) {
 startTransition(() => {
 setTab(nextTab);
 });
 }
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `action`: A function that updates some state by calling one or more [`set` functions](/reference/react/useState#setstate). React calls `action` immediately with no parameters and marks all state updates scheduled synchronously during the `action` function call as Transitions. Any async calls awaited in the `action` will be included in the transition, but currently require wrapping any `set` functions after the `await` in an additional `startTransition` (see [Troubleshooting](/reference/react/useTransition#react-doesnt-treat-my-state-update-after-await-as-a-transition)). State updates marked as Transitions will be [non-blocking](#marking-a-state-update-as-a-non-blocking-transition) and [will not display unwanted loading indicators.](/reference/react/useTransition#preventing-unwanted-loading-indicators).

#### Returns {/*returns*/}

`startTransition` does not return anything.

#### Caveats {/*caveats*/}

* `startTransition` does not provide a way to track whether a Transition is pending. To show a pending indicator while the Transition is ongoing, you need [`useTransition`](/reference/react/useTransition) instead.

* You can wrap an update into a Transition only if you have access to the `set` function of that state. If you want to start a Transition in response to some prop or a custom Hook return value, try [`useDeferredValue`](/reference/react/useDeferredValue) instead.

* The function you pass to `startTransition` is called immediately, marking all state updates that happen while it executes as Transitions. If you try to perform state updates in a `setTimeout`, for example, they won't be marked as Transitions.

* You must wrap any state updates after any async requests in another `startTransition` to mark them as Transitions. This is a known limitation that we will fix in the future (see [Troubleshooting](/reference/react/useTransition#react-doesnt-treat-my-state-update-after-await-as-a-transition)).

* A state update marked as a Transition will be interrupted by other state updates. For example, if you update a chart component inside a Transition, but then start typing into an input while the chart is in the middle of a re-render, React will restart the rendering work on the chart component after handling the input state update.

* Transition updates can't be used to control text inputs.

* If there are multiple ongoing Transitions, React currently batches them together. This is a limitation that may be removed in a future release.

---

## Usage {/*usage*/}

### Marking a state update as a non-blocking Transition {/*marking-a-state-update-as-a-non-blocking-transition*/}

You can mark a state update as a *Transition* by wrapping it in a `startTransition` call:

```js {7,9}
import { startTransition } from 'react';

function TabContainer() {
 const [tab, setTab] = useState('about');

 function selectTab(nextTab) {
 startTransition(() => {
 setTab(nextTab);
 });
 }
 // ...
}
```

Transitions let you keep the user interface updates responsive even on slow devices.

With a Transition, your UI stays responsive in the middle of a re-render. For example, if the user clicks a tab but then change their mind and click another tab, they can do that without waiting for the first re-render to finish.

<Note>

`startTransition` is very similar to [`useTransition`](/reference/react/useTransition), except that it does not provide the `isPending` flag to track whether a Transition is ongoing. You can call `startTransition` when `useTransition` is not available. For example, `startTransition` works outside components, such as from a data library.

[Learn about Transitions and see examples on the `useTransition` page.](/reference/react/useTransition)

</Note>

---
title: use
---

<Intro>

`use` is a React API that lets you read the value of a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) or [context](/learn/passing-data-deeply-with-context).

```js
const value = use(resource);
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `use(context)` {/*use-context*/}

Call `use` with a [context](/learn/passing-data-deeply-with-context) to read its value. Unlike [`useContext`](/reference/react/useContext), `use` can be called within loops and conditional statements like `if`.

```js
import { use } from 'react';

function Button() {
 const theme = use(ThemeContext);
 // ...
```

[See more examples below.](#usage-context)

#### Parameters {/*context-parameters*/}

* `context`: A [context](/learn/passing-data-deeply-with-context) created with [`createContext`](/reference/react/createContext).

#### Returns {/*context-returns*/}

The context value for the passed context, determined by the closest context provider above the calling component. If there is no provider, the returned value is the `defaultValue` passed to [`createContext`](/reference/react/createContext).

#### Caveats {/*context-caveats*/}

* `use` must be called inside a Component or a Hook.
* Reading context with `use` is not supported in [Server Components](/reference/rsc/server-components).

---

### `use(promise)` {/*use-promise*/}

Call `use` with a [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) to read its resolved value. The component calling `use` *suspends* while the Promise is pending. Despite its name, `use` is not a Hook. Unlike Hooks, it can be called inside loops and conditional statements like `if`.

```js
import { use } from 'react';

function MessageComponent({ messagePromise }) {
 const message = use(messagePromise);
 // ...
```

If the component that calls `use` is wrapped in a [Suspense](/reference/react/Suspense) boundary, the fallback will be displayed while the Promise is pending. Once the Promise is resolved, the Suspense fallback is replaced by the rendered components using the data returned by `use`. If the Promise is rejected, the fallback of the nearest [Error Boundary](/reference/react/Component#catching-rendering-errors-with-an-error-boundary) will be displayed.

[See more examples below.](#usage-promises)

#### Parameters {/*promise-parameters*/}

* `promise`: A [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) whose resolved value you want to read. The Promise must be [cached](#caching-promises-for-client-components) so that the same instance is reused across re-renders.

#### Returns {/*promise-returns*/}

The resolved value of the Promise.

#### Caveats {/*promise-caveats*/}

* `use` must be called inside a Component or a Hook.
* `use` cannot be called inside a try-catch block. Instead, wrap your component in an [Error Boundary](#displaying-an-error-with-an-error-boundary) to catch the error and display a fallback.
* Promises passed to `use` must be cached so the same Promise instance is reused across re-renders. [See caching Promises below.](#caching-promises-for-client-components)
* When passing a Promise from a Server Component to a Client Component, its resolved value must be [serializable](/reference/rsc/use-client#serializable-types).

---

## Usage (Context) {/*usage-context*/}

### Reading context with `use` {/*reading-context-with-use*/}

When a [context](/learn/passing-data-deeply-with-context) is passed to `use`, it works similarly to [`useContext`](/reference/react/useContext). While `useContext` must be called at the top level of your component, `use` can be called inside conditionals like `if` and loops like `for`.

```js [[2, 4, "theme"], [1, 4, "ThemeContext"]]
import { use } from 'react';

function Button() {
 const theme = use(ThemeContext);
 // ...
```

`use` returns the <CodeStep step={2}>context value</CodeStep> for the <CodeStep step={1}>context</CodeStep> you passed. To determine the context value, React searches the component tree and finds **the closest context provider above** for that particular context.

To pass context to a `Button`, wrap it or one of its parent components into the corresponding context provider.

```js [[1, 3, "ThemeContext"], [2, 3, "\\"dark\\""], [1, 5, "ThemeContext"]]
function MyPage() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 );
}

function Form() {
 // ... renders buttons inside ...
}
```

It doesn't matter how many layers of components there are between the provider and the `Button`. When a `Button` *anywhere* inside of `Form` calls `use(ThemeContext)`, it will receive `"dark"` as the value.

Unlike [`useContext`](/reference/react/useContext), <CodeStep step={2}>`use`</CodeStep> can be called in conditionals and loops like <CodeStep step={1}>`if`</CodeStep>.

```js [[1, 2, "if"], [2, 3, "use"]]
function HorizontalRule({ show }) {
 if (show) {
 const theme = use(ThemeContext);
 return <hr className={theme} />;
 }
 return false;
}
```

<CodeStep step={2}>`use`</CodeStep> is called from inside a <CodeStep step={1}>`if`</CodeStep> statement, allowing you to conditionally read values from a Context.

<Pitfall>

Like `useContext`, `use(context)` always looks for the closest context provider *above* the component that calls it. It searches upwards and **does not** consider context providers in the component from which you're calling `use(context)`.

</Pitfall>

<Sandpack>

```js
import { createContext, use } from 'react';

const ThemeContext = createContext(null);

export default function MyApp() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 )
}

function Form() {
 return (
 <Panel title="Welcome">
 <Button show={true}>Sign up</Button>
 <Button show={false}>Log in</Button>
 </Panel>
 );
}

function Panel({ title, children }) {
 const theme = use(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ show, children }) {
 if (show) {
 const theme = use(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {children}
 </button>
 );
 }
 return false
}
```

```css
.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

### Reading a Promise from context {/*reading-a-promise-from-context*/}

To share asynchronous data without prop drilling, set a Promise as a context value, then read it with `use(context)` and resolve it with `use(promise)`:

```js
import { use } from 'react';
import { UserContext } from './UserContext';

function Profile() {
 const userPromise = use(UserContext);
 const user = use(userPromise);
 return <h1>{user.name}</h1>;
}
```

Reading the value requires two `use` calls because the context value itself isn't awaited. See [Before you use context](/learn/passing-data-deeply-with-context#before-you-use-context) for alternatives to consider before reaching for context.

Wrap the components that read the Promise in a [Suspense](/reference/react/Suspense) boundary so only that subtree suspends while the Promise is pending. See [Usage (Promises)](#usage-promises) below for more on reading Promises with `use`.

<Pitfall>

When this pattern is used with [Server Components](/reference/rsc/server-components), refetching the Promise requires refetching the Server Component that sets the Promise in context. Avoid setting the Promise in context high in the tree, since that would refetch large parts of the app unnecessarily.

</Pitfall>

---

## Usage (Promises) {/*usage-promises*/}

### Reading a Promise with `use` {/*reading-a-promise-with-use*/}

Call `use` with a Promise to read its resolved value. The component will [suspend](/reference/react/Suspense) while the Promise is pending.

```js [[1, 4, "use(albumsPromise)"]]
import { use } from 'react';

function Albums({ albumsPromise }) {
 const albums = use(albumsPromise);
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

Wrap the component that calls <CodeStep step={1}>`use`</CodeStep> in a [Suspense](/reference/react/Suspense) boundary so React can show a fallback while the Promise is pending. The closest Suspense boundary above the suspending component shows its fallback. Once the Promise resolves, React reads the value with `use` and replaces the fallback with the rendered component.

<Recipes titleText="Reading a Promise with use vs fetching in an Effect" titleId="examples-promise">

#### Fetching data with `use` {/*fetching-data-with-use*/}

In this example, `Albums` calls `use` with a cached Promise. The component suspends while the Promise is pending, and React displays the nearest Suspense fallback. Rejected Promises propagate to the nearest [Error Boundary](/reference/react/Component#catching-rendering-errors-with-an-error-boundary).

<Sandpack>

```js src/App.js active
import { use, Suspense } from 'react';
import { ErrorBoundary } from 'react-error-boundary';
import { fetchData } from './data.js';

export default function App() {
 return (
 <ErrorBoundary fallback={<p>Could not fetch albums.</p>}>
 <Suspense fallback={<Loading />}>
 <Albums />
 </Suspense>
 </ErrorBoundary>
 );
}

function Albums() {
 const albums = use(fetchData('/albums'));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}

function Loading() {
 return <h2>Loading...</h2>;
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/albums') {
 return await getAlbums();
 } else {
 throw Error('Not implemented');
 }
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }];
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "19.0.0",
 "react-dom": "19.0.0",
 "react-scripts": "^5.0.0",
 "react-error-boundary": "4.0.3"
 },
 "main": "/index.js"
}
```

</Sandpack>

<Solution />

#### Fetching data with `useEffect` {/*fetching-data-with-useeffect*/}

Before `use`, a common approach was to fetch data in an Effect and update state when the data arrives. Compared to `use`, this approach requires managing loading and error states manually. For more details on why fetching in an Effect is discouraged, see [You Might Not Need an Effect](/learn/you-might-not-need-an-effect#fetching-data).

<Sandpack>

```js src/App.js active
import { useState, useEffect } from 'react';
import { fetchAlbums } from './data.js';

export default function App() {
 const [albums, setAlbums] = useState(null);
 const [isLoading, setIsLoading] = useState(true);
 const [error, setError] = useState(null);

 useEffect(() => {
 fetchAlbums()
 .then(data => {
 setAlbums(data);
 setIsLoading(false);
 })
 .catch(err => {
 setError(err);
 setIsLoading(false);
 });
 }, []);

 if (isLoading) {
 return <h2>Loading...</h2>;
 }

 if (error) {
 return <p>Error: {error.message}</p>;
 }

 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
export async function fetchAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }];
}
```

</Sandpack>

<Solution />

</Recipes>

<Pitfall>

##### Promises passed to `use` must be cached {/*promises-must-cached*/}

Promises created during render are recreated on every render, which causes React to show the Suspense fallback repeatedly and prevents content from appearing.

```js
function Albums() {
 // 🔴 `fetch` creates a new Promise on every render.
 const albums = use(fetch('/albums'));
 // ...
}
```

Instead, pass a Promise from a cache, a [Suspense-enabled framework](/reference/react/Suspense#suspense-enabled-frameworks), or a Server Component:

```js
// ✅ fetchData reads the Promise from a cache.
const albums = use(fetchData('/albums'));
```

</Pitfall>

<DeepDive>

#### Why are Promises recreated on every render? {/*why-promises-recreated*/}

[React doesn't preserve state for renders that suspended before mounting](/reference/react/Suspense#caveats). After each suspension, React retries rendering from scratch, so any Promise created during render is recreated.

Common ways a Promise can be unintentionally recreated during render:

```js
function Albums() {
 // 🔴 `fetch` creates a new Promise on every render.
 const albums = use(fetch('/albums'));

 // 🔴 Uncached `async` function calls create a new Promise on every render.
 const albums = use((async () => {
 const res = await fetch('/albums');
 return res.json();
 })());

 // 🔴 Adding `.then` returns a new Promise on every render,
 // even if `fetchData` is cached.
 const albums = use(fetchData('/albums').then(res => res.json()));
 // ...
}
```

Ideally, Promises are created before rendering, such as in an event handler, a route loader, or a Server Component, and passed to the component that calls `use`. Fetching lazily in render delays network requests and can create waterfalls.

```js
// ✅ fetchData reads the Promise from a cache.
const albums = use(fetchData('/albums'));
```

</DeepDive>

---

### Caching Promises for Client Components {/*caching-promises-for-client-components*/}

Promises passed to `use` in Client Components must be cached so the same Promise instance is reused across re-renders. If a new Promise is created directly in render, React will display the Suspense fallback on every re-render.

```js
// ✅ Cache the Promise so the same one is reused across renders
let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}
```

The `fetchData` function returns the same Promise each time it's called with the same URL. When `use` receives the same Promise on a re-render, it reads the already-resolved value synchronously without suspending.

<Note>

The way you cache Promises depends on the framework you use with Suspense. Frameworks typically provide built-in caching mechanisms. If you don't use a framework, you can use a simple module-level cache like the one above, or a [Suspense-enabled data source](/reference/react/Suspense#what-activates-a-suspense-boundary).

</Note>

In the example below, clicking "Re-render" updates state in `App` and triggers a re-render. Because `fetchData` returns the same cached Promise, `Albums` reads the value synchronously instead of showing the Suspense fallback again.

<Sandpack>

```js src/App.js active
import { use, Suspense, useState } from 'react';
import { fetchData } from './data.js';

export default function App() {
 const [count, setCount] = useState(0);
 return (
 <>
 <button onClick={() => setCount(count + 1)}>
 Re-render
 </button>
 <p>Render count: {count}</p>
 <Suspense fallback={<p>Loading...</p>}>
 <Albums />
 </Suspense>
 </>
 );
}

function Albums() {
 const albums = use(fetchData('/albums'));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/albums') {
 return await getAlbums();
 } else {
 throw Error('Not implemented');
 }
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }];
}
```

</Sandpack>

<DeepDive>

#### How to implement a promise cache {/*how-to-implement-a-promise-cache*/}

A basic cache stores the Promise keyed by URL so the same instance is reused across renders. To also avoid unnecessary Suspense fallbacks when data is already available, you can set `status` and `value` (or `reason`) fields on the Promise. React checks these fields when `use` is called: if `status` is `'fulfilled'`, it reads `value` synchronously without suspending. If `status` is `'rejected'`, it throws `reason`. If the field is missing or `'pending'`, it suspends.

```js
let cache = new Map();

function fetchData(url) {
 if (!cache.has(url)) {
 const promise = getData(url);
 promise.status = 'pending';
 promise.then(
 value => {
 promise.status = 'fulfilled';
 promise.value = value;
 },
 reason => {
 promise.status = 'rejected';
 promise.reason = reason;
 },
 );
 cache.set(url, promise);
 }
 return cache.get(url);
}
```

This is primarily useful for library authors building Suspense-compatible data layers. React will set the `status` field itself on Promises that don't have it, but setting it yourself avoids an extra render when the data is already available.

This cache pattern is the foundation for [re-fetching data](#re-fetching-data-in-client-components) (where changing the cache key triggers a new fetch) and [preloading data on hover](#preloading-data-on-hover) (where calling `fetchData` early means the Promise may already be resolved by the time `use` reads it).

</DeepDive>

<Pitfall>

##### Don't skip calling `use` based on whether a Promise is already settled. {/*conditional-use*/}

Unlike other hooks, `use` can be called inside conditions and loops — but it must always be called for the Promise itself. Never read `promise.status` or `promise.value` directly to bypass `use`; always pass the Promise to `use` and let React handle it.

```js
// 🔴 Don't bypass `use` by reading promise status directly
if (promise.status === 'fulfilled') {
 return promise.value;
}
const value = use(promise);
```

```js
// ✅ Pass the promise to `use` and let React track the promise
const value = use(promise);
```

Bypassing `use` this way can break React Suspense optimizations and Suspense features for React DevTools. You can `use(promise)` conditionally, but don't conditionally `use(promise)` based on the promise itself.

</Pitfall>

---

### Re-fetching data in Client Components {/*re-fetching-data-in-client-components*/}

To refresh data at the same URL (for example, with a "Refresh" button), invalidate the cache entry and start a new fetch inside a [`startTransition`](/reference/react/startTransition). Store the resulting Promise in state to trigger a re-render. While the new Promise is pending, React keeps showing the existing content because the update is inside a Transition.

```js
function App() {
 const [albumsPromise, setAlbumsPromise] = useState(fetchData('/albums'));
 const [isPending, startTransition] = useTransition();

 function handleRefresh() {
 startTransition(() => {
 setAlbumsPromise(refetchData('/albums'));
 });
 }
 // ...
}
```

`refetchData` clears the old cache entry and starts a new fetch at the same URL. Storing the resulting Promise in state triggers a re-render inside the Transition. On re-render, `Albums` receives the new Promise and `use` suspends on it while React keeps showing the old content.

<Sandpack>

```js src/App.js active
import { Suspense, useState, useTransition } from 'react';
import { use } from 'react';
import { fetchData, refetchData } from './data.js';

export default function App() {
 const [albumsPromise, setAlbumsPromise] = useState(
 () => fetchData('/the-beatles/albums')
 );
 const [isPending, startTransition] = useTransition();

 function handleRefresh() {
 startTransition(() => {
 setAlbumsPromise(refetchData('/the-beatles/albums'));
 });
 }

 return (
 <>
 <button
 onClick={handleRefresh}
 disabled={isPending}
 >
 {isPending ? 'Refreshing...' : 'Refresh'}
 </button>
 <div style={{ opacity: isPending ? 0.6 : 1 }}>
 <Suspense fallback={<Loading />}>
 <Albums albumsPromise={albumsPromise} />
 </Suspense>
 </div>
 </>
 );
}

function Albums({ albumsPromise }) {
 const albums = use(albumsPromise);
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}

function Loading() {
 return <h2>Loading...</h2>;
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

export function refetchData(url) {
 cache.delete(url);
 return fetchData(url);
}

async function getData(url) {
 if (url.startsWith('/the-beatles/albums')) {
 return await getAlbums();
 } else {
 throw Error('Not implemented');
 }
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }];
}
```

```css
button { margin-bottom: 10px; }
```

</Sandpack>

<Note>

Frameworks that support Suspense typically provide their own caching and invalidation mechanisms. The custom cache above is useful for understanding the pattern, but in practice prefer your framework's data fetching solution.

</Note>

---

### Preloading data on hover {/*preloading-data-on-hover*/}

You can start loading data before it's needed by calling `fetchData` during a hover event. Since `fetchData` caches the Promise, the data may already be available by the time the user clicks. If the Promise has resolved by the time `use` reads it, React renders the component immediately without showing a Suspense fallback.

```js
<button
 onMouseEnter={() => fetchData(`/${id}/albums`)}
 onClick={() => {
 startTransition(() => {
 setArtistId(id);
 });
 }}
>
```

In this example, hovering over an artist button starts fetching their albums in the background. Without hovering first, clicking shows a loading fallback. Try hovering over a button for a moment before clicking to see the difference.

<Sandpack>

```js src/App.js active
import { Suspense, useState, useTransition } from 'react';
import Albums from './Albums.js';
import { fetchData } from './data.js';

export default function App() {
 const [artistId, setArtistId] = useState('the-beatles');
 const [isPending, startTransition] = useTransition();

 return (
 <>
 <div>
 {['the-beatles', 'led-zeppelin', 'pink-floyd'].map(id => (
 <button
 key={id}
 onMouseEnter={() => {
 fetchData(`/${id}/albums`);
 }}
 onClick={() => {
 startTransition(() => {
 setArtistId(id);
 });
 }}
 >
 {id === 'the-beatles' ? 'The Beatles' :
 id === 'led-zeppelin' ? 'Led Zeppelin' :
 'Pink Floyd'}
 </button>
 ))}
 </div>
 <Suspense key={artistId} fallback={<Loading />}>
 <Albums artistId={artistId} />
 </Suspense>
 </>
 );
}

function Loading() {
 return <h2>Loading...</h2>;
}
```

```js src/Albums.js
import { use } from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 const promise = getData(url);
 // Set status fields so React can read the value
 // synchronously if the Promise resolves before
 // `use` is called (e.g. when preloading on hover).
 promise.status = 'pending';
 promise.then(
 value => {
 promise.status = 'fulfilled';
 promise.value = value;
 },
 reason => {
 promise.status = 'rejected';
 promise.reason = reason;
 },
 );
 cache.set(url, promise);
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/the-beatles/albums')) {
 return await getAlbums('the-beatles');
 } else if (url.startsWith('/led-zeppelin/albums')) {
 return await getAlbums('led-zeppelin');
 } else if (url.startsWith('/pink-floyd/albums')) {
 return await getAlbums('pink-floyd');
 } else {
 throw Error('Not implemented');
 }
}

async function getAlbums(artistId) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 800);
 });

 if (artistId === 'the-beatles') {
 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }];
 } else if (artistId === 'led-zeppelin') {
 return [{
 id: 10,
 title: 'Coda',
 year: 1982
 }, {
 id: 9,
 title: 'In Through the Out Door',
 year: 1979
 }, {
 id: 8,
 title: 'Presence',
 year: 1976
 }];
 } else {
 return [{
 id: 7,
 title: 'The Wall',
 year: 1979
 }, {
 id: 6,
 title: 'Animals',
 year: 1977
 }, {
 id: 5,
 title: 'Wish You Were Here',
 year: 1975
 }];
 }
}
```

```css
button { margin-right: 10px; }
```

</Sandpack>

---

### Streaming data from server to client {/*streaming-data-from-server-to-client*/}

Data can be streamed from the server to the client by passing a Promise as a prop from a Server Component to a Client Component.

```js
import { fetchMessage } from './lib.js';
import { Message } from './message.js';

export default function App() {
 const messagePromise = fetchMessage();
 return (
 <Suspense fallback={<p>waiting for message...</p>}>
 <Message messagePromise={messagePromise} />
 </Suspense>
 );
}
```

The Client Component then takes the Promise it received as a prop and passes it to the `use` API. This allows the Client Component to read the value from the Promise that was initially created by the Server Component.

```js
// message.js
'use client';

import { use } from 'react';

export function Message({ messagePromise }) {
 const messageContent = use(messagePromise);
 return <p>Here is the message: {messageContent}</p>;
}
```
Because `Message` is wrapped in a [Suspense](/reference/react/Suspense) boundary, the fallback will be displayed until the Promise is resolved. When the Promise is resolved, the value will be read by the `use` API and the `Message` component will replace the Suspense fallback.

<Sandpack>

```js src/message.js active
"use client";

import { use, Suspense } from "react";

function Message({ messagePromise }) {
 const messageContent = use(messagePromise);
 return <p>Here is the message: {messageContent}</p>;
}

export function MessageContainer({ messagePromise }) {
 return (
 <Suspense fallback={<p>⌛Downloading message...</p>}>
 <Message messagePromise={messagePromise} />
 </Suspense>
 );
}
```

```js src/App.js hidden
import { useState } from "react";
import { MessageContainer } from "./message.js";

function fetchMessage() {
 return new Promise((resolve) => setTimeout(resolve, 1000, "⚛️"));
}

export default function App() {
 const [messagePromise, setMessagePromise] = useState(null);
 const [show, setShow] = useState(false);
 function download() {
 setMessagePromise(fetchMessage());
 setShow(true);
 }

 if (show) {
 return <MessageContainer messagePromise={messagePromise} />;
 } else {
 return <button onClick={download}>Download message</button>;
 }
}
```

```js src/index.js hidden
import React, { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';

// TODO: update this example to use
// the Codesandbox Server Component
// demo environment once it is created
import App from './App';

const root = createRoot(document.getElementById('root'));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

</Sandpack>

<DeepDive>

#### Should I resolve a Promise in a Server or Client Component? {/*resolve-promise-in-server-or-client-component*/}

If you have a Promise, at some point you need to unwrap it to read its value. You unwrap it with `await` in a Server Component, and with `use` in a Client Component.

Usually, the simplest option is to `await` the Promise where you create it. The Server Component suspends until the data is ready, and everything below it waits too:

```js
// Server Component
export default async function App() {
 const messageContent = await fetchMessage();
 return <Message messageContent={messageContent} />;
}
```

However, you don't have to unwrap it right away. You can pass the Promise down as a prop, and unwrap it deeper in the tree. The component that reads the Promise still suspends, but only that part of the tree waits for the data. Wrap that component in a [`<Suspense>`](/reference/react/Suspense) boundary to show a fallback while the rest of the page renders immediately.

For example, a deeper Server Component can `await` the Promise it receives:

```js
import { Suspense } from 'react';

// Server Component
export default function App() {
 const messagePromise = fetchMessage();
 return (
 <Suspense fallback={<p>⌛Downloading message...</p>}>
 <Message messagePromise={messagePromise} />
 </Suspense>
 );
}

// Server Component
async function Message({ messagePromise }) {
 const messageContent = await messagePromise;
 return <p>{messageContent}</p>;
}
```

Or, in a separate file, a Client Component can unwrap the same Promise with `use`:

```js
// Client Component
'use client';

import { use } from 'react';

export function Message({ messagePromise }) {
 const messageContent = use(messagePromise);
 return <p>{messageContent}</p>;
}
```

Passing the Promise down works the same way in both cases. Both suspend where the Promise is read, and both unblock the UI above. The only difference is that Client Components can't `await` during render, so they unwrap the Promise with `use` instead. A common case is interactive content like popovers and tooltips, where the data is only needed after a hover or click.

See [Revealing content together at once](/reference/react/Suspense#revealing-content-together-at-once) for guidance on where to place Suspense boundaries.

</DeepDive>

---

### Displaying an error with an Error Boundary {/*displaying-an-error-with-an-error-boundary*/}

If the Promise passed to `use` is rejected, the error propagates to the nearest [Error Boundary](/reference/react/Component#catching-rendering-errors-with-an-error-boundary). Wrap the component that calls `use` in an Error Boundary to display a fallback when the Promise is rejected.

In the example below, `fetchData` rejects on the first attempt and succeeds on retry. The Error Boundary catches the rejection and shows a fallback with a "Try again" button.

<Sandpack>

```js src/App.js active
import { use, Suspense, useState, startTransition } from "react";
import { ErrorBoundary } from "react-error-boundary";
import { fetchData, refetchData } from "./data.js";

export default function App() {
 const [albumsPromise, setAlbumsPromise] = useState(
 () => fetchData('/the-beatles/albums')
 );

 function handleRetry() {
 startTransition(() => {
 setAlbumsPromise(refetchData('/the-beatles/albums'));
 });
 }

 return (
 <ErrorBoundary
 resetKeys={[albumsPromise]}
 fallbackRender={() => (
 <>
 <p>⚠️ Something went wrong loading the albums.</p>
 <button onClick={handleRetry}>Try again</button>
 </>
 )}
 >
 <Suspense fallback={<p>Loading...</p>}>
 <Albums albumsPromise={albumsPromise} />
 </Suspense>
 </ErrorBoundary>
 );
}

function Albums({ albumsPromise }) {
 const albums = use(albumsPromise);
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();
let retried = false;

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

export function refetchData(url) {
 cache.delete(url);
 retried = true;
 return fetchData(url);
}

async function getData(url) {
 // Add a fake delay to make the loading state visible.
 await new Promise(resolve => setTimeout(resolve, 1000));
 if (url === '/the-beatles/albums') {
 // Fail the first attempt to demonstrate the Error Boundary,
 // then succeed on retry.
 if (!retried) {
 throw new Error('Example Error: Failed to fetch albums');
 }
 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }];
 }
 throw new Error('Not implemented');
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "19.0.0",
 "react-dom": "19.0.0",
 "react-scripts": "^5.0.0",
 "react-error-boundary": "4.0.3"
 },
 "main": "/index.js"
}
```
</Sandpack>

---

## Troubleshooting {/*troubleshooting*/}

### I'm getting an error: "Suspense Exception: This is not a real error!" {/*suspense-exception-error*/}

You are calling `use` inside a try-catch block. `use` throws internally to integrate with Suspense, so it cannot be wrapped in try-catch. Instead, wrap the component that calls `use` in an [Error Boundary](#displaying-an-error-with-an-error-boundary) to handle errors.

```jsx
function Albums({ albumsPromise }) {
 try {
 // ❌ Don't wrap `use` in try-catch
 const albums = use(albumsPromise);
 } catch (e) {
 return <p>Error</p>;
 }
 // ...
```

Instead, wrap the component in an Error Boundary:

```jsx
function Albums({ albumsPromise }) {
 // ✅ Call `use` without try-catch
 const albums = use(albumsPromise);
 // ...
```

```jsx
// ✅ Use an Error Boundary to handle errors
<ErrorBoundary fallback={<p>Error</p>}>
 <Albums albumsPromise={albumsPromise} />
</ErrorBoundary>
```

---

### I'm getting a warning: "A component was suspended by an uncached promise" {/*uncached-promise-error*/}

The Promise passed to `use` is not cached, so React cannot reuse it across re-renders.

This commonly happens when calling `fetch` or an `async` function directly in render:

```js
function Albums() {
 // 🔴 This creates a new Promise on every render
 const albums = use(fetch('/albums'));
 // ...
}
```

To fix this, cache the Promise so the same instance is reused:

```js
// ✅ fetchData returns the same Promise for the same URL
const albums = use(fetchData('/albums'));
```

See [caching Promises for Client Components](#caching-promises-for-client-components) for more details.

---
title: useActionState
---

<Intro>

`useActionState` is a React Hook that lets you update state with side effects using [Actions](/reference/react/useTransition#functions-called-in-starttransition-are-called-actions).

```js
const [state, dispatchAction, isPending] = useActionState(reducerAction, initialState, permalink?);
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useActionState(reducerAction, initialState, permalink?)` {/*useactionstate*/}

Call `useActionState` at the top level of your component to create state for the result of an Action.

```js
import { useActionState } from 'react';

function reducerAction(previousState, actionPayload) {
 // ...
}

function MyCart({initialState}) {
 const [state, dispatchAction, isPending] = useActionState(reducerAction, initialState);
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `reducerAction`: The function to be called when the Action is triggered. When called, it receives the previous state (initially the `initialState` you provided, then its previous return value) as its first argument, followed by the `actionPayload` passed to `dispatchAction`.
* `initialState`: The value you want the state to be initially. React ignores this argument after `dispatchAction` is invoked for the first time.
* **optional** `permalink`: A string containing the unique page URL that this form modifies.
 * For use on pages with [React Server Components](/reference/rsc/server-components) with progressive enhancement.
 * If `reducerAction` is a [Server Function](/reference/rsc/server-functions) and the form is submitted before the JavaScript bundle loads, the browser will navigate to the specified permalink URL rather than the current page's URL.

#### Returns {/*returns*/}

`useActionState` returns an array with exactly three values:

1. The current state. During the first render, it will match the `initialState` you passed. After `dispatchAction` is invoked, it will match the value returned by the `reducerAction`.
2. A `dispatchAction` function that you call inside [Actions](/reference/react/useTransition#functions-called-in-starttransition-are-called-actions).
3. The `isPending` flag that tells you if any dispatched Actions for this Hook are pending.

#### Caveats {/*caveats*/}

* `useActionState` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.
* React queues and executes multiple calls to `dispatchAction` sequentially. Each call to `reducerAction` receives the result of the previous call.
* The `dispatchAction` function has a stable identity, so you will often see it omitted from Effect dependencies, but including it will not cause the Effect to fire. If the linter lets you omit a dependency without errors, it is safe to do. [Learn more about removing Effect dependencies.](/learn/removing-effect-dependencies#move-dynamic-objects-and-functions-inside-your-effect)
* When using the `permalink` option, ensure the same form component is rendered on the destination page (including the same `reducerAction` and `permalink`) so React knows how to pass the state through. Once the page becomes interactive, this parameter has no effect.
* When using Server Functions, `initialState` needs to be [serializable](/reference/rsc/use-server#serializable-parameters-and-return-values) (values like plain objects, arrays, strings, and numbers).
* If `dispatchAction` throws an error, React cancels all queued actions and shows the nearest [Error Boundary](/reference/react/Component#catching-rendering-errors-with-an-error-boundary).
* If there are multiple ongoing Actions, React batches them together. This is a limitation that may be removed in a future release.

<Note>

`dispatchAction` must be called from an Action.

You can wrap it in [`startTransition`](/reference/react/startTransition), or pass it to an [Action prop](/reference/react/useTransition#exposing-action-props-from-components). Calls outside that scope won’t be treated as part of the Transition and [log an error](#async-function-outside-transition) on development mode.

</Note>

---

### `reducerAction` function {/*reduceraction*/}

The `reducerAction` function passed to `useActionState` receives the previous state and returns a new state.

Unlike reducers in `useReducer`, the `reducerAction` can be async and perform side effects:

```js
async function reducerAction(previousState, actionPayload) {
 const newState = await post(actionPayload);
 return newState;
}
```

Each time you call `dispatchAction`, React calls the `reducerAction` with the `actionPayload`. The reducer will perform side effects such as posting data, and return the new state. If `dispatchAction` is called multiple times, React queues and executes them in order so the result of the previous call is passed as `previousState` for the current call.

#### Parameters {/*reduceraction-parameters*/}

* `previousState`: The last state. Initially this is equal to the `initialState`. After the first call to `dispatchAction`, it's equal to the last state returned.

* **optional** `actionPayload`: The argument passed to `dispatchAction`. It can be a value of any type. Similar to `useReducer` conventions, it is usually an object with a `type` property identifying it and, optionally, other properties with additional information.

#### Returns {/*reduceraction-returns*/}

`reducerAction` returns the new state, and triggers a Transition to re-render with that state.

#### Caveats {/*reduceraction-caveats*/}

* `reducerAction` can be sync or async. It can perform sync actions like showing a notification, or async actions like posting updates to a server.
* `reducerAction` is not invoked twice in `<StrictMode>` since `reducerAction` is designed to allow side effects.
* The return type of `reducerAction` must match the type of `initialState`. If TypeScript infers a mismatch, you may need to explicitly annotate your state type.
* If you set state after `await` in the `reducerAction` you currently need to wrap the state update in an additional `startTransition`. See the [startTransition](/reference/react/useTransition#react-doesnt-treat-my-state-update-after-await-as-a-transition) docs for more info.
* When using Server Functions, `actionPayload` needs to be [serializable](/reference/rsc/use-server#serializable-parameters-and-return-values) (values like plain objects, arrays, strings, and numbers).

<DeepDive>

#### Why is it called `reducerAction`? {/*why-is-it-called-reduceraction*/}

The function passed to `useActionState` is called a *reducer action* because:

- It *reduces* the previous state into a new state, like `useReducer`.
- It's an *Action* because it's called inside a Transition and can perform side effects.

Conceptually, `useActionState` is like `useReducer`, but you can do side effects in the reducer.

</DeepDive>

---

## Usage {/*usage*/}

### Adding state to an Action {/*adding-state-to-an-action*/}

Call `useActionState` at the top level of your component to create state for the result of an Action.

```js [[1, 7, "count"], [2, 7, "dispatchAction"], [3, 7, "isPending"]]
import { useActionState } from 'react';

async function addToCartAction(prevCount) {
 // ...
}
function Counter() {
 const [count, dispatchAction, isPending] = useActionState(addToCartAction, 0);

 // ...
}
```

`useActionState` returns an array with exactly three items:

1. The <CodeStep step={1}>current state</CodeStep>, initially set to the initial state you provided.
2. The <CodeStep step={2}>action dispatcher</CodeStep> that lets you trigger `reducerAction`.
3. A <CodeStep step={3}>pending state</CodeStep> that tells you whether the Action is in progress.

To call `addToCartAction`, call the <CodeStep step={2}>action dispatcher</CodeStep>. React will queue calls to `addToCartAction` with the previous count.

<Sandpack>

```js src/App.js
import { useActionState, startTransition } from 'react';
import { addToCart } from './api';
import Total from './Total';

export default function Checkout() {
 const [count, dispatchAction, isPending] = useActionState(async (prevCount) => {
 return await addToCart(prevCount)
 }, 0);

 function handleClick() {
 startTransition(() => {
 dispatchAction();
 });
 }

 return (
 <div className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <span>Qty: {count}</span>
 </div>
 <div className="row">
 <button onClick={handleClick}>Add Ticket{isPending ? ' 🌀' : ' '}</button>
 </div>
 <hr />
 <Total quantity={count} />
 </div>
 );
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity}) {
 return (
 <div className="row total">
 <span>Total</span>
 <span>{formatter.format(quantity * 9999)}</span>
 </div>
 );
}
```

```js src/api.js
export async function addToCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return count + 1;
}

export async function removeFromCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return Math.max(0, count - 1);
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.row button {
 margin-left: auto;
 min-width: 150px;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}

button {
 padding: 8px 16px;
 cursor: pointer;
}
```

</Sandpack>

Every time you click "Add Ticket," React queues a call to `addToCartAction`. React shows the pending state until all the tickets are added, and then re-renders with the final state.

<DeepDive>

#### How `useActionState` queuing works {/*how-useactionstate-queuing-works*/}

Try clicking "Add Ticket" multiple times. Every time you click, a new `addToCartAction` is queued. Since there's an artificial 1 second delay, that means 4 clicks will take ~4 seconds to complete.

**This is intentional in the design of `useActionState`.**

We have to wait for the previous result of `addToCartAction` in order to pass the `prevCount` to the next call to `addToCartAction`. That means React has to wait for the previous Action to finish before calling the next Action.

You can typically solve this by [using with useOptimistic](/reference/react/useActionState#using-with-useoptimistic) but for more complex cases you may want to consider [cancelling queued actions](#cancelling-queued-actions) or not using `useActionState`.

</DeepDive>

---

### Using multiple Action types {/*using-multiple-action-types*/}

To handle multiple types, you can pass an argument to `dispatchAction`.

By convention, it is common to write it as a switch statement. For each case in the switch, calculate and return some next state. The argument can have any shape, but it is common to pass objects with a `type` property identifying the action.

<Sandpack>

```js src/App.js
import { useActionState, startTransition } from 'react';
import { addToCart, removeFromCart } from './api';
import Total from './Total';

export default function Checkout() {
 const [count, dispatchAction, isPending] = useActionState(updateCartAction, 0);

 function handleAdd() {
 startTransition(() => {
 dispatchAction({ type: 'ADD' });
 });
 }

 function handleRemove() {
 startTransition(() => {
 dispatchAction({ type: 'REMOVE' });
 });
 }

 return (
 <div className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <span className="stepper">
 <span className="qty">{isPending ? '🌀' : count}</span>
 <span className="buttons">
 <button onClick={handleAdd}>▲</button>
 <button onClick={handleRemove}>▼</button>
 </span>
 </span>
 </div>
 <hr />
 <Total quantity={count} isPending={isPending}/>
 </div>
 );
}

async function updateCartAction(prevCount, actionPayload) {
 switch (actionPayload.type) {
 case 'ADD': {
 return await addToCart(prevCount);
 }
 case 'REMOVE': {
 return await removeFromCart(prevCount);
 }
 }
 return prevCount;
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity, isPending}) {
 return (
 <div className="row total">
 <span>Total</span>
 {isPending ? '🌀 Updating...' : formatter.format(quantity * 9999)}
 </div>
 );
}
```

```js src/api.js hidden
export async function addToCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return count + 1;
}

export async function removeFromCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return Math.max(0, count - 1);
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.stepper {
 display: flex;
 align-items: center;
 gap: 8px;
}

.qty {
 min-width: 20px;
 text-align: center;
}

.buttons {
 display: flex;
 flex-direction: column;
 gap: 2px;
}

.buttons button {
 padding: 0 8px;
 font-size: 10px;
 line-height: 1.2;
 cursor: pointer;
}

.pending {
 width: 20px;
 text-align: center;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}
```

</Sandpack>

When you click to increase or decrease the quantity, an `"ADD"` or `"REMOVE"` is dispatched. In the `reducerAction`, different APIs are called to update the quantity.

In this example, we use the pending state of the Actions to replace both the quantity and the total. If you want to provide immediate feedback, such as immediately updating the quantity, you can use `useOptimistic`.

<DeepDive>

#### How is `useActionState` different from `useReducer`? {/*useactionstate-vs-usereducer*/}

You might notice this example looks a lot like `useReducer`, but they serve different purposes:

- **Use `useReducer`** to manage state of your UI. The reducer must be pure.

- **Use `useActionState`** to manage state of your Actions. The reducer can perform side effects.

You can think of `useActionState` as `useReducer` for side effects from user Actions. Since it computes the next Action to take based on the previous Action, it has to [order the calls sequentially](/reference/react/useActionState#how-useactionstate-queuing-works). If you want to perform Actions in parallel, use `useState` and `useTransition` directly.

</DeepDive>

---

### Using with `useOptimistic` {/*using-with-useoptimistic*/}

You can combine `useActionState` with [`useOptimistic`](/reference/react/useOptimistic) to show immediate UI feedback:

<Sandpack>

```js src/App.js
import { useActionState, startTransition, useOptimistic } from 'react';
import { addToCart, removeFromCart } from './api';
import Total from './Total';

export default function Checkout() {
 const [count, dispatchAction, isPending] = useActionState(updateCartAction, 0);
 const [optimisticCount, setOptimisticCount] = useOptimistic(count);

 function handleAdd() {
 startTransition(() => {
 setOptimisticCount(c => c + 1);
 dispatchAction({ type: 'ADD' });
 });
 }

 function handleRemove() {
 startTransition(() => {
 setOptimisticCount(c => c - 1);
 dispatchAction({ type: 'REMOVE' });
 });
 }

 return (
 <div className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <span className="stepper">
 <span className="pending">{isPending && '🌀'}</span>
 <span className="qty">{optimisticCount}</span>
 <span className="buttons">
 <button onClick={handleAdd}>▲</button>
 <button onClick={handleRemove}>▼</button>
 </span>
 </span>
 </div>
 <hr />
 <Total quantity={optimisticCount} isPending={isPending}/>
 </div>
 );
}

async function updateCartAction(prevCount, actionPayload) {
 switch (actionPayload.type) {
 case 'ADD': {
 return await addToCart(prevCount);
 }
 case 'REMOVE': {
 return await removeFromCart(prevCount);
 }
 }
 return prevCount;
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity, isPending}) {
 return (
 <div className="row total">
 <span>Total</span>
 <span>{isPending ? '🌀 Updating...' : formatter.format(quantity * 9999)}</span>
 </div>
 );
}
```

```js src/api.js hidden
export async function addToCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return count + 1;
}

export async function removeFromCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return Math.max(0, count - 1);
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.stepper {
 display: flex;
 align-items: center;
 gap: 8px;
}

.qty {
 min-width: 20px;
 text-align: center;
}

.buttons {
 display: flex;
 flex-direction: column;
 gap: 2px;
}

.buttons button {
 padding: 0 8px;
 font-size: 10px;
 line-height: 1.2;
 cursor: pointer;
}

.pending {
 width: 20px;
 text-align: center;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}
```

</Sandpack>

`setOptimisticCount` immediately updates the quantity, and `dispatchAction()` queues the `updateCartAction`. A pending indicator appears on both the quantity and total to give the user feedback that their update is still being applied.

---

### Using with Action props {/*using-with-action-props*/}

When you pass the `dispatchAction` function to a component that exposes an [Action prop](/reference/react/useTransition#exposing-action-props-from-components), you don't need to call `startTransition` or `useOptimistic` yourself.

This example shows using the `increaseAction` and `decreaseAction` props of a QuantityStepper component:

<Sandpack>

```js src/App.js
import { useActionState } from 'react';
import { addToCart, removeFromCart } from './api';
import QuantityStepper from './QuantityStepper';
import Total from './Total';

export default function Checkout() {
 const [count, dispatchAction, isPending] = useActionState(updateCartAction, 0);

 function addAction() {
 dispatchAction({type: 'ADD'});
 }

 function removeAction() {
 dispatchAction({type: 'REMOVE'});
 }

 return (
 <div className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <QuantityStepper
 value={count}
 increaseAction={addAction}
 decreaseAction={removeAction}
 />
 </div>
 <hr />
 <Total quantity={count} isPending={isPending} />
 </div>
 );
}

async function updateCartAction(prevCount, actionPayload) {
 switch (actionPayload.type) {
 case 'ADD': {
 return await addToCart(prevCount);
 }
 case 'REMOVE': {
 return await removeFromCart(prevCount);
 }
 }
 return prevCount;
}
```

```js src/QuantityStepper.js
import { startTransition, useOptimistic } from 'react';

export default function QuantityStepper({value, increaseAction, decreaseAction}) {
 const [optimisticValue, setOptimisticValue] = useOptimistic(value);
 const isPending = value !== optimisticValue;
 function handleIncrease() {
 startTransition(async () => {
 setOptimisticValue(c => c + 1);
 await increaseAction();
 });
 }

 function handleDecrease() {
 startTransition(async () => {
 setOptimisticValue(c => Math.max(0, c - 1));
 await decreaseAction();
 });
 }

 return (
 <span className="stepper">
 <span className="pending">{isPending && '🌀'}</span>
 <span className="qty">{optimisticValue}</span>
 <span className="buttons">
 <button onClick={handleIncrease}>▲</button>
 <button onClick={handleDecrease}>▼</button>
 </span>
 </span>
 );
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity, isPending}) {
 return (
 <div className="row total">
 <span>Total</span>
 {isPending ? '🌀 Updating...' : formatter.format(quantity * 9999)}
 </div>
 );
}
```

```js src/api.js hidden
export async function addToCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return count + 1;
}

export async function removeFromCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return Math.max(0, count - 1);
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.stepper {
 display: flex;
 align-items: center;
 gap: 8px;
}

.qty {
 min-width: 20px;
 text-align: center;
}

.buttons {
 display: flex;
 flex-direction: column;
 gap: 2px;
}

.buttons button {
 padding: 0 8px;
 font-size: 10px;
 line-height: 1.2;
 cursor: pointer;
}

.pending {
 width: 20px;
 text-align: center;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}
```

</Sandpack>

Since `<QuantityStepper>` has built-in support for transitions, pending state, and optimistically updating the count, you just need to tell the Action _what_ to change, and _how_ to change it is handled for you.

---

### Cancelling queued Actions {/*cancelling-queued-actions*/}

You can use an `AbortController` to cancel pending Actions:

<Sandpack>

```js src/App.js
import { useActionState, useRef } from 'react';
import { addToCart, removeFromCart } from './api';
import QuantityStepper from './QuantityStepper';
import Total from './Total';

export default function Checkout() {
 const abortRef = useRef(null);
 const [count, dispatchAction, isPending] = useActionState(updateCartAction, 0);

 async function addAction() {
 if (abortRef.current) {
 abortRef.current.abort();
 }
 abortRef.current = new AbortController();
 await dispatchAction({ type: 'ADD', signal: abortRef.current.signal });
 }

 async function removeAction() {
 if (abortRef.current) {
 abortRef.current.abort();
 }
 abortRef.current = new AbortController();
 await dispatchAction({ type: 'REMOVE', signal: abortRef.current.signal });
 }

 return (
 <div className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <QuantityStepper
 value={count}
 increaseAction={addAction}
 decreaseAction={removeAction}
 />
 </div>
 <hr />
 <Total quantity={count} isPending={isPending} />
 </div>
 );
}

async function updateCartAction(prevCount, actionPayload) {
 switch (actionPayload.type) {
 case 'ADD': {
 try {
 return await addToCart(prevCount, { signal: actionPayload.signal });
 } catch (e) {
 return prevCount + 1;
 }
 }
 case 'REMOVE': {
 try {
 return await removeFromCart(prevCount, { signal: actionPayload.signal });
 } catch (e) {
 return Math.max(0, prevCount - 1);
 }
 }
 }
 return prevCount;
}
```

```js src/QuantityStepper.js
import { startTransition, useOptimistic } from 'react';

export default function QuantityStepper({value, increaseAction, decreaseAction}) {
 const [optimisticValue, setOptimisticValue] = useOptimistic(value);
 const isPending = value !== optimisticValue;
 function handleIncrease() {
 startTransition(async () => {
 setOptimisticValue(c => c + 1);
 await increaseAction();
 });
 }

 function handleDecrease() {
 startTransition(async () => {
 setOptimisticValue(c => Math.max(0, c - 1));
 await decreaseAction();
 });
 }

 return (
 <span className="stepper">
 <span className="pending">{isPending && '🌀'}</span>
 <span className="qty">{optimisticValue}</span>
 <span className="buttons">
 <button onClick={handleIncrease}>▲</button>
 <button onClick={handleDecrease}>▼</button>
 </span>
 </span>
 );
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity, isPending}) {
 return (
 <div className="row total">
 <span>Total</span>
 {isPending ? '🌀 Updating...' : formatter.format(quantity * 9999)}
 </div>
 );
}
```

```js src/api.js hidden
class AbortError extends Error {
 name = 'AbortError';
 constructor(message = 'The operation was aborted') {
 super(message);
 }
}

function sleep(ms, signal) {
 if (!signal) return new Promise((resolve) => setTimeout(resolve, ms));
 if (signal.aborted) return Promise.reject(new AbortError());

 return new Promise((resolve, reject) => {
 const id = setTimeout(() => {
 signal.removeEventListener('abort', onAbort);
 resolve();
 }, ms);

 const onAbort = () => {
 clearTimeout(id);
 reject(new AbortError());
 };

 signal.addEventListener('abort', onAbort, { once: true });
 });
}
export async function addToCart(count, opts) {
 await sleep(1000, opts?.signal);
 return count + 1;
}

export async function removeFromCart(count, opts) {
 await sleep(1000, opts?.signal);
 return Math.max(0, count - 1);
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.stepper {
 display: flex;
 align-items: center;
 gap: 8px;
}

.qty {
 min-width: 20px;
 text-align: center;
}

.buttons {
 display: flex;
 flex-direction: column;
 gap: 2px;
}

.buttons button {
 padding: 0 8px;
 font-size: 10px;
 line-height: 1.2;
 cursor: pointer;
}

.pending {
 width: 20px;
 text-align: center;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}
```

</Sandpack>

Try clicking increase or decrease multiple times, and notice that the total updates within 1 second no matter how many times you click. This works because it uses an `AbortController` to "complete" the previous Action so the next Action can proceed.

<Pitfall>

Aborting an Action isn't always safe.

For example, if the Action performs a mutation (like writing to a database), aborting the network request doesn't undo the server-side change. This is why `useActionState` doesn't abort by default. It's only safe when you know the side effect can be safely ignored or retried.

</Pitfall>

---

### Using with `<form>` Action props {/*use-with-a-form*/}

You can pass the `dispatchAction` function as the `action` prop to a `<form>`.

When used this way, React automatically wraps the submission in a Transition, so you don't need to call `startTransition` yourself. The `reducerAction` receives the previous state and the submitted `FormData`:

<Sandpack>

```js src/App.js
import { useActionState, useOptimistic } from 'react';
import { addToCart, removeFromCart } from './api';
import Total from './Total';

export default function Checkout() {
 const [count, dispatchAction, isPending] = useActionState(updateCartAction, 0);
 const [optimisticCount, setOptimisticCount] = useOptimistic(count);

 async function formAction(formData) {
 const type = formData.get('type');
 if (type === 'ADD') {
 setOptimisticCount(c => c + 1);
 } else {
 setOptimisticCount(c => Math.max(0, c - 1));
 }
 return dispatchAction(formData);
 }

 return (
 <form action={formAction} className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <span className="stepper">
 <span className="pending">{isPending && '🌀'}</span>
 <span className="qty">{optimisticCount}</span>
 <span className="buttons">
 <button type="submit" name="type" value="ADD">▲</button>
 <button type="submit" name="type" value="REMOVE">▼</button>
 </span>
 </span>
 </div>
 <hr />
 <Total quantity={count} isPending={isPending} />
 </form>
 );
}

async function updateCartAction(prevCount, formData) {
 const type = formData.get('type');
 switch (type) {
 case 'ADD': {
 return await addToCart(prevCount);
 }
 case 'REMOVE': {
 return await removeFromCart(prevCount);
 }
 }
 return prevCount;
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity, isPending}) {
 return (
 <div className="row total">
 <span>Total</span>
 {isPending ? '🌀 Updating...' : formatter.format(quantity * 9999)}
 </div>
 );
}
```

```js src/api.js hidden
export async function addToCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return count + 1;
}

export async function removeFromCart(count) {
 await new Promise(resolve => setTimeout(resolve, 1000));
 return Math.max(0, count - 1);
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.stepper {
 display: flex;
 align-items: center;
 gap: 8px;
}

.qty {
 min-width: 20px;
 text-align: center;
}

.buttons {
 display: flex;
 flex-direction: column;
 gap: 2px;
}

.buttons button {
 padding: 0 8px;
 font-size: 10px;
 line-height: 1.2;
 cursor: pointer;
}

.pending {
 width: 20px;
 text-align: center;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}
```

</Sandpack>

In this example, when the user clicks the stepper arrows, the button submits the form and `useActionState` calls `updateCartAction` with the form data. The example uses `useOptimistic` to immediately show the new quantity while the server confirms the update.

<RSC>

When used with a [Server Function](/reference/rsc/server-functions), `useActionState` allows the server's response to be shown before hydration (when React attaches to server-rendered HTML) completes. You can also use the optional `permalink` parameter for progressive enhancement (allowing the form to work before JavaScript loads) on pages with dynamic content. This is typically handled by your framework for you.

</RSC>

See the [`<form>`](/reference/react-dom/components/form#handle-form-submission-with-a-server-function) docs for more information on using Actions with forms.

---

### Handling errors {/*handling-errors*/}

There are two ways to handle errors with `useActionState`.

For known errors, such as "quantity not available" validation errors from your backend, you can return it as part of your `reducerAction` state and display it in the UI.

For unknown errors, such as `undefined is not a function`, you can throw an error. React will cancel all queued Actions and shows the nearest [Error Boundary](/reference/react/Component#catching-rendering-errors-with-an-error-boundary) by rethrowing the error from the `useActionState` hook.

<Sandpack>

```js src/App.js
import {useActionState, startTransition} from 'react';
import {ErrorBoundary} from 'react-error-boundary';
import {addToCart} from './api';
import Total from './Total';

function Checkout() {
 const [state, dispatchAction, isPending] = useActionState(
 async (prevState, quantity) => {
 const result = await addToCart(prevState.count, quantity);
 if (result.error) {
 // Return the error from the API as state
 return {...prevState, error: `Could not add quanitiy ${quantity}: ${result.error}`};
 }

 if (!isPending) {
 // Clear the error state for the first dispatch.
 return {count: result.count, error: null};
 }

 // Return the new count, and any errors that happened.
 return {count: result.count, error: prevState.error};

 },
 {
 count: 0,
 error: null,
 }
 );

 function handleAdd(quantity) {
 startTransition(() => {
 dispatchAction(quantity);
 });
 }

 return (
 <div className="checkout">
 <h2>Checkout</h2>
 <div className="row">
 <span>Eras Tour Tickets</span>
 <span>
 {isPending && '🌀 '}Qty: {state.count}
 </span>
 </div>
 <div className="buttons">
 <button onClick={() => handleAdd(1)}>Add 1</button>
 <button onClick={() => handleAdd(10)}>Add 10</button>
 <button onClick={() => handleAdd(NaN)}>Add NaN</button>
 </div>
 {state.error && <div className="error">{state.error}</div>}
 <hr />
 <Total quantity={state.count} isPending={isPending} />
 </div>
 );
}

export default function App() {
 return (
 <ErrorBoundary
 fallbackRender={({resetErrorBoundary}) => (
 <div className="checkout">
 <h2>Something went wrong</h2>
 <p>The action could not be completed.</p>
 <button onClick={resetErrorBoundary}>Try again</button>
 </div>
 )}>
 <Checkout />
 </ErrorBoundary>
 );
}
```

```js src/Total.js
const formatter = new Intl.NumberFormat('en-US', {
 style: 'currency',
 currency: 'USD',
 minimumFractionDigits: 0,
});

export default function Total({quantity, isPending}) {
 return (
 <div className="row total">
 <span>Total</span>
 <span>
 {isPending ? '🌀 Updating...' : formatter.format(quantity * 9999)}
 </span>
 </div>
 );
}
```

```js src/api.js hidden
export async function addToCart(count, quantity) {
 await new Promise((resolve) => setTimeout(resolve, 1000));
 if (quantity > 5) {
 return {error: 'Quantity not available'};
 } else if (isNaN(quantity)) {
 throw new Error('Quantity must be a number');
 }
 return {count: count + quantity};
}
```

```css
.checkout {
 display: flex;
 flex-direction: column;
 gap: 12px;
 padding: 16px;
 border: 1px solid #ccc;
 border-radius: 8px;
 font-family: system-ui;
}

.checkout h2 {
 margin: 0 0 8px 0;
}

.row {
 display: flex;
 justify-content: space-between;
 align-items: center;
}

.total {
 font-weight: bold;
}

hr {
 width: 100%;
 border: none;
 border-top: 1px solid #ccc;
 margin: 4px 0;
}

button {
 padding: 8px 16px;
 cursor: pointer;
}

.buttons {
 display: flex;
 gap: 8px;
}

.error {
 color: red;
 font-size: 14px;
}
```

```json package.json hidden
{
 "dependencies": {
 "react": "19.0.0",
 "react-dom": "19.0.0",
 "react-scripts": "^5.0.0",
 "react-error-boundary": "4.0.3"
 },
 "main": "/index.js"
}
```

</Sandpack>

In this example, "Add 10" simulates an API that returns a validation error, which `updateCartAction` stores in state and displays inline. "Add NaN" results in an invalid count, so `updateCartAction` throws, which propagates through `useActionState` to the `ErrorBoundary` and shows a reset UI.

---

## Troubleshooting {/*troubleshooting*/}

### My `isPending` flag is not updating {/*ispending-not-updating*/}

If you're calling `dispatchAction` manually (not through an Action prop), make sure you wrap the call in [`startTransition`](/reference/react/startTransition):

```js
import { useActionState, startTransition } from 'react';

function MyComponent() {
 const [state, dispatchAction, isPending] = useActionState(myAction, null);

 function handleClick() {
 // ✅ Correct: wrap in startTransition
 startTransition(() => {
 dispatchAction();
 });
 }

 // ...
}
```

When `dispatchAction` is passed to an Action prop, React automatically wraps it in a Transition.

---

### My Action cannot read form data {/*action-cannot-read-form-data*/}

When you use `useActionState`, the `reducerAction` receives an extra argument as its first argument: the previous or initial state. The submitted form data is therefore its second argument instead of its first.

```js {2,7}
// Without useActionState
function action(formData) {
 const name = formData.get('name');
}

// With useActionState
function action(prevState, formData) {
 const name = formData.get('name');
}
```

---

### My actions are being skipped {/*actions-skipped*/}

If you call `dispatchAction` multiple times and some of them don't run, it may be because an earlier `dispatchAction` call threw an error.

When a `reducerAction` throws, React skips all subsequently queued `dispatchAction` calls.

To handle this, catch errors within your `reducerAction` and return an error state instead of throwing:

```js
async function myReducerAction(prevState, data) {
 try {
 const result = await submitData(data);
 return { success: true, data: result };
 } catch (error) {
 // ✅ Return error state instead of throwing
 return { success: false, error: error.message };
 }
}
```

---

### My state doesn't reset {/*reset-state*/}

`useActionState` doesn't provide a built-in reset function. To reset the state, you can design your `reducerAction` to handle a reset signal:

```js
const initialState = { name: '', error: null };

async function formAction(prevState, payload) {
 // Handle reset
 if (payload === null) {
 return initialState;
 }
 // Normal action logic
 const result = await submitData(payload);
 return result;
}

function MyComponent() {
 const [state, dispatchAction, isPending] = useActionState(formAction, initialState);

 function handleReset() {
 startTransition(() => {
 dispatchAction(null); // Pass null to trigger reset
 });
 }

 // ...
}
```

Alternatively, you can add a `key` prop to the component using `useActionState` to force it to remount with fresh state, or a `<form>` `action` prop, which resets automatically after submission.

---

### I'm getting an error: "An async function with useActionState was called outside of a transition." {/*async-function-outside-transition*/}

A common mistake is to forget to call `dispatchAction` from inside a Transition:

<ConsoleBlockMulti>
<ConsoleLogLine level="error">

An async function with useActionState was called outside of a transition. This is likely not what you intended (for example, isPending will not update correctly). Either call the returned function inside startTransition, or pass it to an `action` or `formAction` prop.

</ConsoleLogLine>
</ConsoleBlockMulti>

This error happens because `dispatchAction` must run inside a Transition:

```js
function MyComponent() {
 const [state, dispatchAction, isPending] = useActionState(myAsyncAction, null);

 function handleClick() {
 // ❌ Wrong: calling dispatchAction outside a Transition
 dispatchAction();
 }

 // ...
}
```

To fix, either wrap the call in [`startTransition`](/reference/react/startTransition):

```js
import { useActionState, startTransition } from 'react';

function MyComponent() {
 const [state, dispatchAction, isPending] = useActionState(myAsyncAction, null);

 function handleClick() {
 // ✅ Correct: wrap in startTransition
 startTransition(() => {
 dispatchAction();
 });
 }

 // ...
}
```

Or pass `dispatchAction` to an Action prop, is call in a Transition:

```js
function MyComponent() {
 const [state, dispatchAction, isPending] = useActionState(myAsyncAction, null);

 // ✅ Correct: action prop wraps in a Transition for you
 return <Button action={dispatchAction}>...</Button>;
}
```

---

### I'm getting an error: "Cannot update action state while rendering" {/*cannot-update-during-render*/}

You cannot call `dispatchAction` during render:

<ConsoleBlock level="error">

Cannot update action state while rendering.

</ConsoleBlock>

This causes an infinite loop because calling `dispatchAction` schedules a state update, which triggers a re-render, which calls `dispatchAction` again.

```js
function MyComponent() {
 const [state, dispatchAction, isPending] = useActionState(myAction, null);

 // ❌ Wrong: calling dispatchAction during render
 dispatchAction();

 // ...
}
```

To fix, only call `dispatchAction` in response to user events (like form submissions or button clicks).

---
title: useCallback
---

<Intro>

`useCallback` is a React Hook that lets you cache a function definition between re-renders.

```js
const cachedFn = useCallback(fn, dependencies)
```

</Intro>

<Note>

[React Compiler](/learn/react-compiler) automatically memoizes values and functions, reducing the need for manual `useCallback` calls. You can use the compiler to handle memoization automatically.

</Note>

<InlineToc />

---

## Reference {/*reference*/}

### `useCallback(fn, dependencies)` {/*usecallback*/}

Call `useCallback` at the top level of your component to cache a function definition between re-renders:

```js {4,9}
import { useCallback } from 'react';

export default function ProductPage({ productId, referrer, theme }) {
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }, [productId, referrer]);
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `fn`: The function value that you want to cache. It can take any arguments and return any values. React will return (not call!) your function back to you during the initial render. On next renders, React will give you the same function again if the `dependencies` have not changed since the last render. Otherwise, it will give you the function that you have passed during the current render, and store it in case it can be reused later. React will not call your function. The function is returned to you so you can decide when and whether to call it.

* `dependencies`: The list of all reactive values referenced inside of the `fn` code. Reactive values include props, state, and all the variables and functions declared directly inside your component body. If your linter is [configured for React](/learn/editor-setup#linting), it will verify that every reactive value is correctly specified as a dependency. The list of dependencies must have a constant number of items and be written inline like `[dep1, dep2, dep3]`. React will compare each dependency with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison algorithm.

#### Returns {/*returns*/}

On the initial render, `useCallback` returns the `fn` function you have passed.

During subsequent renders, it will either return an already stored `fn` function from the last render (if the dependencies haven't changed), or return the `fn` function you have passed during this render.

#### Caveats {/*caveats*/}

* `useCallback` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.
* React **will not throw away the cached function unless there is a specific reason to do that.** For example, in development, React throws away the cache when you edit the file of your component. Both in development and in production, React will throw away the cache if your component suspends during the initial mount. In the future, React may add more features that take advantage of throwing away the cache--for example, if React adds built-in support for virtualized lists in the future, it would make sense to throw away the cache for items that scroll out of the virtualized table viewport. This should match your expectations if you rely on `useCallback` as a performance optimization. Otherwise, a [state variable](/reference/react/useState#im-trying-to-set-state-to-a-function-but-it-gets-called-instead) or a [ref](/reference/react/useRef#avoiding-recreating-the-ref-contents) may be more appropriate.

---

## Usage {/*usage*/}

### Skipping re-rendering of components {/*skipping-re-rendering-of-components*/}

When you optimize rendering performance, you will sometimes need to cache the functions that you pass to child components. Let's first look at the syntax for how to do this, and then see in which cases it's useful.

To cache a function between re-renders of your component, wrap its definition into the `useCallback` Hook:

```js [[3, 4, "handleSubmit"], [2, 9, "[productId, referrer]"]]
import { useCallback } from 'react';

function ProductPage({ productId, referrer, theme }) {
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }, [productId, referrer]);
 // ...
```

You need to pass two things to `useCallback`:

1. A function definition that you want to cache between re-renders.
2. A <CodeStep step={2}>list of dependencies</CodeStep> including every value within your component that's used inside your function.

On the initial render, the <CodeStep step={3}>returned function</CodeStep> you'll get from `useCallback` will be the function you passed.

On the following renders, React will compare the <CodeStep step={2}>dependencies</CodeStep> with the dependencies you passed during the previous render. If none of the dependencies have changed (compared with [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is)), `useCallback` will return the same function as before. Otherwise, `useCallback` will return the function you passed on *this* render.

In other words, `useCallback` caches a function between re-renders until its dependencies change.

**Let's walk through an example to see when this is useful.**

Say you're passing a `handleSubmit` function down from the `ProductPage` to the `ShippingForm` component:

```js {5}
function ProductPage({ productId, referrer, theme }) {
 // ...
 return (
 <div className={theme}>
 <ShippingForm onSubmit={handleSubmit} />
 </div>
 );
```

You've noticed that toggling the `theme` prop freezes the app for a moment, but if you remove `<ShippingForm />` from your JSX, it feels fast. This tells you that it's worth trying to optimize the `ShippingForm` component.

**By default, when a component re-renders, React re-renders all of its children recursively.** This is why, when `ProductPage` re-renders with a different `theme`, the `ShippingForm` component *also* re-renders. This is fine for components that don't require much calculation to re-render. But if you verified a re-render is slow, you can tell `ShippingForm` to skip re-rendering when its props are the same as on last render by wrapping it in [`memo`:](/reference/react/memo)

```js {3,5}
import { memo } from 'react';

const ShippingForm = memo(function ShippingForm({ onSubmit }) {
 // ...
});
```

**With this change, `ShippingForm` will skip re-rendering if all of its props are the *same* as on the last render.** This is when caching a function becomes important! Let's say you defined `handleSubmit` without `useCallback`:

```js {2,3,8,12-13}
function ProductPage({ productId, referrer, theme }) {
 // Every time the theme changes, this will be a different function...
 function handleSubmit(orderDetails) {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }

 return (
 <div className={theme}>
 {/* ... so ShippingForm's props will never be the same, and it will re-render every time */}
 <ShippingForm onSubmit={handleSubmit} />
 </div>
 );
}
```

**In JavaScript, a `function () {}` or `() => {}` always creates a _different_ function,** similar to how the `{}` object literal always creates a new object. Normally, this wouldn't be a problem, but it means that `ShippingForm` props will never be the same, and your [`memo`](/reference/react/memo) optimization won't work. This is where `useCallback` comes in handy:

```js {2,3,8,12-13}
function ProductPage({ productId, referrer, theme }) {
 // Tell React to cache your function between re-renders...
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }, [productId, referrer]); // ...so as long as these dependencies don't change...

 return (
 <div className={theme}>
 {/* ...ShippingForm will receive the same props and can skip re-rendering */}
 <ShippingForm onSubmit={handleSubmit} />
 </div>
 );
}
```

**By wrapping `handleSubmit` in `useCallback`, you ensure that it's the *same* function between the re-renders** (until dependencies change). You don't *have to* wrap a function in `useCallback` unless you do it for some specific reason. In this example, the reason is that you pass it to a component wrapped in [`memo`,](/reference/react/memo) and this lets it skip re-rendering. There are other reasons you might need `useCallback` which are described further on this page.

<Note>

**You should only rely on `useCallback` as a performance optimization.** If your code doesn't work without it, find the underlying problem and fix it first. Then you may add `useCallback` back.

</Note>

<DeepDive>

#### How is useCallback related to useMemo? {/*how-is-usecallback-related-to-usememo*/}

You will often see [`useMemo`](/reference/react/useMemo) alongside `useCallback`. They are both useful when you're trying to optimize a child component. They let you [memoize](https://en.wikipedia.org/wiki/Memoization) (or, in other words, cache) something you're passing down:

```js {6-8,10-15,19}
import { useMemo, useCallback } from 'react';

function ProductPage({ productId, referrer }) {
 const product = useData('/product/' + productId);

 const requirements = useMemo(() => { // Calls your function and caches its result
 return computeRequirements(product);
 }, [product]);

 const handleSubmit = useCallback((orderDetails) => { // Caches your function itself
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }, [productId, referrer]);

 return (
 <div className={theme}>
 <ShippingForm requirements={requirements} onSubmit={handleSubmit} />
 </div>
 );
}
```

The difference is in *what* they're letting you cache:

* **[`useMemo`](/reference/react/useMemo) caches the *result* of calling your function.** In this example, it caches the result of calling `computeRequirements(product)` so that it doesn't change unless `product` has changed. This lets you pass the `requirements` object down without unnecessarily re-rendering `ShippingForm`. When necessary, React will call the function you've passed during rendering to calculate the result.
* **`useCallback` caches *the function itself.*** Unlike `useMemo`, it does not call the function you provide. Instead, it caches the function you provided so that `handleSubmit` *itself* doesn't change unless `productId` or `referrer` has changed. This lets you pass the `handleSubmit` function down without unnecessarily re-rendering `ShippingForm`. Your code won't run until the user submits the form.

If you're already familiar with [`useMemo`,](/reference/react/useMemo) you might find it helpful to think of `useCallback` as this:

```js {expectedErrors: {'react-compiler': [3]}}
// Simplified implementation (inside React)
function useCallback(fn, dependencies) {
 return useMemo(() => fn, dependencies);
}
```

[Read more about the difference between `useMemo` and `useCallback`.](/reference/react/useMemo#memoizing-a-function)

</DeepDive>

<DeepDive>

#### Should you add useCallback everywhere? {/*should-you-add-usecallback-everywhere*/}

If your app is like this site, and most interactions are coarse (like replacing a page or an entire section), memoization is usually unnecessary. On the other hand, if your app is more like a drawing editor, and most interactions are granular (like moving shapes), then you might find memoization very helpful.

Caching a function with `useCallback` is only valuable in a few cases:

- You pass it as a prop to a component wrapped in [`memo`.](/reference/react/memo) You want to skip re-rendering if the value hasn't changed. Memoization lets your component re-render only if dependencies changed.
- The function you're passing is later used as a dependency of some Hook. For example, another function wrapped in `useCallback` depends on it, or you depend on this function from [`useEffect.`](/reference/react/useEffect)

There is no benefit to wrapping a function in `useCallback` in other cases. There is no significant harm to doing that either, so some teams choose to not think about individual cases, and memoize as much as possible. The downside is that code becomes less readable. Also, not all memoization is effective: a single value that's "always new" is enough to break memoization for an entire component.

Note that `useCallback` does not prevent *creating* the function. You're always creating a function (and that's fine!), but React ignores it and gives you back a cached function if nothing changed.

**In practice, you can make a lot of memoization unnecessary by following a few principles:**

1. When a component visually wraps other components, let it [accept JSX as children.](/learn/passing-props-to-a-component#passing-jsx-as-children) Then, if the wrapper component updates its own state, React knows that its children don't need to re-render.
2. Prefer local state and don't [lift state up](/learn/sharing-state-between-components) any further than necessary. Don't keep transient state like forms and whether an item is hovered at the top of your tree or in a global state library.
3. Keep your [rendering logic pure.](/learn/keeping-components-pure) If re-rendering a component causes a problem or produces some noticeable visual artifact, it's a bug in your component! Fix the bug instead of adding memoization.
4. Avoid [unnecessary Effects that update state.](/learn/you-might-not-need-an-effect) Most performance problems in React apps are caused by chains of updates originating from Effects that cause your components to render over and over.
5. Try to [remove unnecessary dependencies from your Effects.](/learn/removing-effect-dependencies) For example, instead of memoization, it's often simpler to move some object or a function inside an Effect or outside the component.

If a specific interaction still feels laggy, [use the React Developer Tools profiler](https://legacy.reactjs.org/blog/2018/09/10/introducing-the-react-profiler.html) to see which components benefit the most from memoization, and add memoization where needed. These principles make your components easier to debug and understand, so it's good to follow them in any case. In long term, we're researching [doing memoization automatically](https://www.youtube.com/watch?v=lGEMwh32soc) to solve this once and for all.

</DeepDive>

<Recipes titleText="The difference between useCallback and declaring a function directly" titleId="examples-rerendering">

#### Skipping re-rendering with `useCallback` and `memo` {/*skipping-re-rendering-with-usecallback-and-memo*/}

In this example, the `ShippingForm` component is **artificially slowed down** so that you can see what happens when a React component you're rendering is genuinely slow. Try incrementing the counter and toggling the theme.

Incrementing the counter feels slow because it forces the slowed down `ShippingForm` to re-render. That's expected because the counter has changed, and so you need to reflect the user's new choice on the screen.

Next, try toggling the theme. **Thanks to `useCallback` together with [`memo`](/reference/react/memo), it’s fast despite the artificial slowdown!** `ShippingForm` skipped re-rendering because the `handleSubmit` function has not changed. The `handleSubmit` function has not changed because both `productId` and `referrer` (your `useCallback` dependencies) haven't changed since last render.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import ProductPage from './ProductPage.js';

export default function App() {
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <ProductPage
 referrerId="wizard_of_oz"
 productId={123}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}
```

```js src/ProductPage.js active
import { useCallback } from 'react';
import ShippingForm from './ShippingForm.js';

export default function ProductPage({ productId, referrer, theme }) {
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }, [productId, referrer]);

 return (
 <div className={theme}>
 <ShippingForm onSubmit={handleSubmit} />
 </div>
 );
}

function post(url, data) {
 // Imagine this sends a request...
 console.log('POST /' + url);
 console.log(data);
}
```

```js {expectedErrors: {'react-compiler': [7, 8]}} src/ShippingForm.js
import { memo, useState } from 'react';

const ShippingForm = memo(function ShippingForm({ onSubmit }) {
 const [count, setCount] = useState(1);

 console.log('[ARTIFICIALLY SLOW] Rendering <ShippingForm />');
 let startTime = performance.now();
 while (performance.now() - startTime < 500) {
 // Do nothing for 500 ms to emulate extremely slow code
 }

 function handleSubmit(e) {
 e.preventDefault();
 const formData = new FormData(e.target);
 const orderDetails = {
 ...Object.fromEntries(formData),
 count
 };
 onSubmit(orderDetails);
 }

 return (
 <form onSubmit={handleSubmit}>
 <p><b>Note: <code>ShippingForm</code> is artificially slowed down!</b></p>
 <label>
 Number of items:
 <button type="button" onClick={() => setCount(count - 1)}>–</button>
 {count}
 <button type="button" onClick={() => setCount(count + 1)}>+</button>
 </label>
 <label>
 Street:
 <input name="street" />
 </label>
 <label>
 City:
 <input name="city" />
 </label>
 <label>
 Postal code:
 <input name="zipCode" />
 </label>
 <button type="submit">Submit</button>
 </form>
 );
});

export default ShippingForm;
```

```css
label {
 display: block; margin-top: 10px;
}

input {
 margin-left: 5px;
}

button[type="button"] {
 margin: 5px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

<Solution />

#### Always re-rendering a component {/*always-re-rendering-a-component*/}

In this example, the `ShippingForm` implementation is also **artificially slowed down** so that you can see what happens when some React component you're rendering is genuinely slow. Try incrementing the counter and toggling the theme.

Unlike in the previous example, toggling the theme is also slow now! This is because **there is no `useCallback` call in this version,** so `handleSubmit` is always a new function, and the slowed down `ShippingForm` component can't skip re-rendering.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import ProductPage from './ProductPage.js';

export default function App() {
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <ProductPage
 referrerId="wizard_of_oz"
 productId={123}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}
```

```js src/ProductPage.js active
import ShippingForm from './ShippingForm.js';

export default function ProductPage({ productId, referrer, theme }) {
 function handleSubmit(orderDetails) {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }

 return (
 <div className={theme}>
 <ShippingForm onSubmit={handleSubmit} />
 </div>
 );
}

function post(url, data) {
 // Imagine this sends a request...
 console.log('POST /' + url);
 console.log(data);
}
```

```js {expectedErrors: {'react-compiler': [7, 8]}} src/ShippingForm.js
import { memo, useState } from 'react';

const ShippingForm = memo(function ShippingForm({ onSubmit }) {
 const [count, setCount] = useState(1);

 console.log('[ARTIFICIALLY SLOW] Rendering <ShippingForm />');
 let startTime = performance.now();
 while (performance.now() - startTime < 500) {
 // Do nothing for 500 ms to emulate extremely slow code
 }

 function handleSubmit(e) {
 e.preventDefault();
 const formData = new FormData(e.target);
 const orderDetails = {
 ...Object.fromEntries(formData),
 count
 };
 onSubmit(orderDetails);
 }

 return (
 <form onSubmit={handleSubmit}>
 <p><b>Note: <code>ShippingForm</code> is artificially slowed down!</b></p>
 <label>
 Number of items:
 <button type="button" onClick={() => setCount(count - 1)}>–</button>
 {count}
 <button type="button" onClick={() => setCount(count + 1)}>+</button>
 </label>
 <label>
 Street:
 <input name="street" />
 </label>
 <label>
 City:
 <input name="city" />
 </label>
 <label>
 Postal code:
 <input name="zipCode" />
 </label>
 <button type="submit">Submit</button>
 </form>
 );
});

export default ShippingForm;
```

```css
label {
 display: block; margin-top: 10px;
}

input {
 margin-left: 5px;
}

button[type="button"] {
 margin: 5px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

However, here is the same code **with the artificial slowdown removed.** Does the lack of `useCallback` feel noticeable or not?

<Sandpack>

```js src/App.js
import { useState } from 'react';
import ProductPage from './ProductPage.js';

export default function App() {
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <ProductPage
 referrerId="wizard_of_oz"
 productId={123}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}
```

```js src/ProductPage.js active
import ShippingForm from './ShippingForm.js';

export default function ProductPage({ productId, referrer, theme }) {
 function handleSubmit(orderDetails) {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }

 return (
 <div className={theme}>
 <ShippingForm onSubmit={handleSubmit} />
 </div>
 );
}

function post(url, data) {
 // Imagine this sends a request...
 console.log('POST /' + url);
 console.log(data);
}
```

```js src/ShippingForm.js
import { memo, useState } from 'react';

const ShippingForm = memo(function ShippingForm({ onSubmit }) {
 const [count, setCount] = useState(1);

 console.log('Rendering <ShippingForm />');

 function handleSubmit(e) {
 e.preventDefault();
 const formData = new FormData(e.target);
 const orderDetails = {
 ...Object.fromEntries(formData),
 count
 };
 onSubmit(orderDetails);
 }

 return (
 <form onSubmit={handleSubmit}>
 <label>
 Number of items:
 <button type="button" onClick={() => setCount(count - 1)}>–</button>
 {count}
 <button type="button" onClick={() => setCount(count + 1)}>+</button>
 </label>
 <label>
 Street:
 <input name="street" />
 </label>
 <label>
 City:
 <input name="city" />
 </label>
 <label>
 Postal code:
 <input name="zipCode" />
 </label>
 <button type="submit">Submit</button>
 </form>
 );
});

export default ShippingForm;
```

```css
label {
 display: block; margin-top: 10px;
}

input {
 margin-left: 5px;
}

button[type="button"] {
 margin: 5px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

Quite often, code without memoization works fine. If your interactions are fast enough, you don't need memoization.

Keep in mind that you need to run React in production mode, disable [React Developer Tools](/learn/react-developer-tools), and use devices similar to the ones your app's users have in order to get a realistic sense of what's actually slowing down your app.

<Solution />

</Recipes>

---

### Updating state from a memoized callback {/*updating-state-from-a-memoized-callback*/}

Sometimes, you might need to update state based on previous state from a memoized callback.

This `handleAddTodo` function specifies `todos` as a dependency because it computes the next todos from it:

```js {6,7}
function TodoList() {
 const [todos, setTodos] = useState([]);

 const handleAddTodo = useCallback((text) => {
 const newTodo = { id: nextId++, text };
 setTodos([...todos, newTodo]);
 }, [todos]);
 // ...
```

You'll usually want memoized functions to have as few dependencies as possible. When you read some state only to calculate the next state, you can remove that dependency by passing an [updater function](/reference/react/useState#updating-state-based-on-the-previous-state) instead:

```js {6,7}
function TodoList() {
 const [todos, setTodos] = useState([]);

 const handleAddTodo = useCallback((text) => {
 const newTodo = { id: nextId++, text };
 setTodos(todos => [...todos, newTodo]);
 }, []); // ✅ No need for the todos dependency
 // ...
```

Here, instead of making `todos` a dependency and reading it inside, you pass an instruction about *how* to update the state (`todos => [...todos, newTodo]`) to React. [Read more about updater functions.](/reference/react/useState#updating-state-based-on-the-previous-state)

---

### Preventing an Effect from firing too often {/*preventing-an-effect-from-firing-too-often*/}

Sometimes, you might want to call a function from inside an [Effect:](/learn/synchronizing-with-effects)

```js {4-9,12}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 function createOptions() {
 return {
 serverUrl: 'https://localhost:1234',
 roomId: roomId
 };
 }

 useEffect(() => {
 const options = createOptions();
 const connection = createConnection(options);
 connection.connect();
 // ...
```

This creates a problem. [Every reactive value must be declared as a dependency of your Effect.](/learn/lifecycle-of-reactive-effects#react-verifies-that-you-specified-every-reactive-value-as-a-dependency) However, if you declare `createOptions` as a dependency, it will cause your Effect to constantly reconnect to the chat room:

```js {6}
 useEffect(() => {
 const options = createOptions();
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [createOptions]); // 🔴 Problem: This dependency changes on every render
 // ...
```

To solve this, you can wrap the function you need to call from an Effect into `useCallback`:

```js {4-9,16}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 const createOptions = useCallback(() => {
 return {
 serverUrl: 'https://localhost:1234',
 roomId: roomId
 };
 }, [roomId]); // ✅ Only changes when roomId changes

 useEffect(() => {
 const options = createOptions();
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [createOptions]); // ✅ Only changes when createOptions changes
 // ...
```

This ensures that the `createOptions` function is the same between re-renders if the `roomId` is the same. **However, it's even better to remove the need for a function dependency.** Move your function *inside* the Effect:

```js {5-10,16}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 useEffect(() => {
 function createOptions() { // ✅ No need for useCallback or function dependencies!
 return {
 serverUrl: 'https://localhost:1234',
 roomId: roomId
 };
 }

 const options = createOptions();
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]); // ✅ Only changes when roomId changes
 // ...
```

Now your code is simpler and doesn't need `useCallback`. [Learn more about removing Effect dependencies.](/learn/removing-effect-dependencies#move-dynamic-objects-and-functions-inside-your-effect)

---

### Optimizing a custom Hook {/*optimizing-a-custom-hook*/}

If you're writing a [custom Hook,](/learn/reusing-logic-with-custom-hooks) it's recommended to wrap any functions that it returns into `useCallback`:

```js {4-6,8-10}
function useRouter() {
 const { dispatch } = useContext(RouterStateContext);

 const navigate = useCallback((url) => {
 dispatch({ type: 'navigate', url });
 }, [dispatch]);

 const goBack = useCallback(() => {
 dispatch({ type: 'back' });
 }, [dispatch]);

 return {
 navigate,
 goBack,
 };
}
```

This ensures that the consumers of your Hook can optimize their own code when needed.

---

## Troubleshooting {/*troubleshooting*/}

### Every time my component renders, `useCallback` returns a different function {/*every-time-my-component-renders-usecallback-returns-a-different-function*/}

Make sure you've specified the dependency array as a second argument!

If you forget the dependency array, `useCallback` will return a new function every time:

```js {7}
function ProductPage({ productId, referrer }) {
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }); // 🔴 Returns a new function every time: no dependency array
 // ...
```

This is the corrected version passing the dependency array as a second argument:

```js {7}
function ProductPage({ productId, referrer }) {
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails,
 });
 }, [productId, referrer]); // ✅ Does not return a new function unnecessarily
 // ...
```

If this doesn't help, then the problem is that at least one of your dependencies is different from the previous render. You can debug this problem by manually logging your dependencies to the console:

```js {5}
 const handleSubmit = useCallback((orderDetails) => {
 // ..
 }, [productId, referrer]);

 console.log([productId, referrer]);
```

You can then right-click on the arrays from different re-renders in the console and select "Store as a global variable" for both of them. Assuming the first one got saved as `temp1` and the second one got saved as `temp2`, you can then use the browser console to check whether each dependency in both arrays is the same:

```js
Object.is(temp1[0], temp2[0]); // Is the first dependency the same between the arrays?
Object.is(temp1[1], temp2[1]); // Is the second dependency the same between the arrays?
Object.is(temp1[2], temp2[2]); // ... and so on for every dependency ...
```

When you find which dependency is breaking memoization, either find a way to remove it, or [memoize it as well.](/reference/react/useMemo#memoizing-a-dependency-of-another-hook)

---

### I need to call `useCallback` for each list item in a loop, but it's not allowed {/*i-need-to-call-usememo-for-each-list-item-in-a-loop-but-its-not-allowed*/}

Suppose the `Chart` component is wrapped in [`memo`](/reference/react/memo). You want to skip re-rendering every `Chart` in the list when the `ReportList` component re-renders. However, you can't call `useCallback` in a loop:

```js {expectedErrors: {'react-compiler': [6]}} {5-14}
function ReportList({ items }) {
 return (
 <article>
 {items.map(item => {
 // 🔴 You can't call useCallback in a loop like this:
 const handleClick = useCallback(() => {
 sendReport(item)
 }, [item]);

 return (
 <figure key={item.id}>
 <Chart onClick={handleClick} />
 </figure>
 );
 })}
 </article>
 );
}
```

Instead, extract a component for an individual item, and put `useCallback` there:

```js {5,12-21}
function ReportList({ items }) {
 return (
 <article>
 {items.map(item =>
 <Report key={item.id} item={item} />
 )}
 </article>
 );
}

function Report({ item }) {
 // ✅ Call useCallback at the top level:
 const handleClick = useCallback(() => {
 sendReport(item)
 }, [item]);

 return (
 <figure>
 <Chart onClick={handleClick} />
 </figure>
 );
}
```

Alternatively, you could remove `useCallback` in the last snippet and instead wrap `Report` itself in [`memo`.](/reference/react/memo) If the `item` prop does not change, `Report` will skip re-rendering, so `Chart` will skip re-rendering too:

```js {5,6-8,15}
function ReportList({ items }) {
 // ...
}

const Report = memo(function Report({ item }) {
 function handleClick() {
 sendReport(item);
 }

 return (
 <figure>
 <Chart onClick={handleClick} />
 </figure>
 );
});
```

---
title: useContext
---

<Intro>

`useContext` is a React Hook that lets you read and subscribe to [context](/learn/passing-data-deeply-with-context) from your component.

```js
const value = useContext(SomeContext)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useContext(SomeContext)` {/*usecontext*/}

Call `useContext` at the top level of your component to read and subscribe to [context.](/learn/passing-data-deeply-with-context)

```js
import { useContext } from 'react';

function MyComponent() {
 const theme = useContext(ThemeContext);
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `SomeContext`: The context that you've previously created with [`createContext`](/reference/react/createContext). The context itself does not hold the information, it only represents the kind of information you can provide or read from components.

#### Returns {/*returns*/}

`useContext` returns the context value for the calling component. It is determined as the `value` passed to the closest `SomeContext` above the calling component in the tree. If there is no such provider, then the returned value will be the `defaultValue` you have passed to [`createContext`](/reference/react/createContext) for that context. The returned value is always up-to-date. React automatically re-renders components that read some context if it changes.

#### Caveats {/*caveats*/}

* `useContext()` call in a component is not affected by providers returned from the *same* component. The corresponding `<Context>` **needs to be *above*** the component doing the `useContext()` call.
* React **automatically re-renders** all the children that use a particular context starting from the provider that receives a different `value`. The previous and the next values are compared with the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. Skipping re-renders with [`memo`](/reference/react/memo) does not prevent the children receiving fresh context values.
* If your build system produces duplicates modules in the output (which can happen with symlinks), this can break context. Passing something via context only works if `SomeContext` that you use to provide context and `SomeContext` that you use to read it are ***exactly* the same object**, as determined by a `===` comparison.

---

## Usage {/*usage*/}

### Passing data deeply into the tree {/*passing-data-deeply-into-the-tree*/}

Call `useContext` at the top level of your component to read and subscribe to [context.](/learn/passing-data-deeply-with-context)

```js [[2, 4, "theme"], [1, 4, "ThemeContext"]]
import { useContext } from 'react';

function Button() {
 const theme = useContext(ThemeContext);
 // ...
```

`useContext` returns the <CodeStep step={2}>context value</CodeStep> for the <CodeStep step={1}>context</CodeStep> you passed. To determine the context value, React searches the component tree and finds **the closest context provider above** for that particular context.

To pass context to a `Button`, wrap it or one of its parent components into the corresponding context provider:

```js [[1, 3, "ThemeContext"], [2, 3, "\\"dark\\""], [1, 5, "ThemeContext"]]
function MyPage() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 );
}

function Form() {
 // ... renders buttons inside ...
}
```

It doesn't matter how many layers of components there are between the provider and the `Button`. When a `Button` *anywhere* inside of `Form` calls `useContext(ThemeContext)`, it will receive `"dark"` as the value.

<Pitfall>

`useContext()` always looks for the closest provider *above* the component that calls it. It searches upwards and **does not** consider providers in the component from which you're calling `useContext()`.

</Pitfall>

<Sandpack>

```js
import { createContext, useContext } from 'react';

const ThemeContext = createContext(null);

export default function MyApp() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 )
}

function Form() {
 return (
 <Panel title="Welcome">
 <Button>Sign up</Button>
 <Button>Log in</Button>
 </Panel>
 );
}

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {children}
 </button>
 );
}
```

```css
.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

---

### Updating data passed via context {/*updating-data-passed-via-context*/}

Often, you'll want the context to change over time. To update context, combine it with [state.](/reference/react/useState) Declare a state variable in the parent component, and pass the current state down as the <CodeStep step={2}>context value</CodeStep> to the provider.

```js {2} [[1, 4, "ThemeContext"], [2, 4, "theme"], [1, 11, "ThemeContext"]]
function MyPage() {
 const [theme, setTheme] = useState('dark');
 return (
 <ThemeContext value={theme}>
 <Form />
 <Button onClick={() => {
 setTheme('light');
 }}>
 Switch to light theme
 </Button>
 </ThemeContext>
 );
}
```

Now any `Button` inside of the provider will receive the current `theme` value. If you call `setTheme` to update the `theme` value that you pass to the provider, all `Button` components will re-render with the new `'light'` value.

<Recipes titleText="Examples of updating context" titleId="examples-basic">

#### Updating a value via context {/*updating-a-value-via-context*/}

In this example, the `MyApp` component holds a state variable which is then passed to the `ThemeContext` provider. Checking the "Dark mode" checkbox updates the state. Changing the provided value re-renders all the components using that context.

<Sandpack>

```js
import { createContext, useContext, useState } from 'react';

const ThemeContext = createContext(null);

export default function MyApp() {
 const [theme, setTheme] = useState('light');
 return (
 <ThemeContext value={theme}>
 <Form />
 <label>
 <input
 type="checkbox"
 checked={theme === 'dark'}
 onChange={(e) => {
 setTheme(e.target.checked ? 'dark' : 'light')
 }}
 />
 Use dark mode
 </label>
 </ThemeContext>
 )
}

function Form({ children }) {
 return (
 <Panel title="Welcome">
 <Button>Sign up</Button>
 <Button>Log in</Button>
 </Panel>
 );
}

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {children}
 </button>
 );
}
```

```css
.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
 margin-bottom: 10px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

Note that `value="dark"` passes the `"dark"` string, but `value={theme}` passes the value of the JavaScript `theme` variable with [JSX curly braces.](/learn/javascript-in-jsx-with-curly-braces) Curly braces also let you pass context values that aren't strings.

<Solution />

#### Updating an object via context {/*updating-an-object-via-context*/}

In this example, there is a `currentUser` state variable which holds an object. You combine `{ currentUser, setCurrentUser }` into a single object and pass it down through the context inside the `value={}`. This lets any component below, such as `LoginButton`, read both `currentUser` and `setCurrentUser`, and then call `setCurrentUser` when needed.

<Sandpack>

```js
import { createContext, useContext, useState } from 'react';

const CurrentUserContext = createContext(null);

export default function MyApp() {
 const [currentUser, setCurrentUser] = useState(null);
 return (
 <CurrentUserContext
 value={{
 currentUser,
 setCurrentUser
 }}
 >
 <Form />
 </CurrentUserContext>
 );
}

function Form({ children }) {
 return (
 <Panel title="Welcome">
 <LoginButton />
 </Panel>
 );
}

function LoginButton() {
 const {
 currentUser,
 setCurrentUser
 } = useContext(CurrentUserContext);

 if (currentUser !== null) {
 return <p>You logged in as {currentUser.name}.</p>;
 }

 return (
 <Button onClick={() => {
 setCurrentUser({ name: 'Advika' })
 }}>Log in as Advika</Button>
 );
}

function Panel({ title, children }) {
 return (
 <section className="panel">
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children, onClick }) {
 return (
 <button className="button" onClick={onClick}>
 {children}
 </button>
 );
}
```

```css
label {
 display: block;
}

.panel {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
 margin-bottom: 10px;
}

.button {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}
```

</Sandpack>

<Solution />

#### Multiple contexts {/*multiple-contexts*/}

In this example, there are two independent contexts. `ThemeContext` provides the current theme, which is a string, while `CurrentUserContext` holds the object representing the current user.

<Sandpack>

```js
import { createContext, useContext, useState } from 'react';

const ThemeContext = createContext(null);
const CurrentUserContext = createContext(null);

export default function MyApp() {
 const [theme, setTheme] = useState('light');
 const [currentUser, setCurrentUser] = useState(null);
 return (
 <ThemeContext value={theme}>
 <CurrentUserContext
 value={{
 currentUser,
 setCurrentUser
 }}
 >
 <WelcomePanel />
 <label>
 <input
 type="checkbox"
 checked={theme === 'dark'}
 onChange={(e) => {
 setTheme(e.target.checked ? 'dark' : 'light')
 }}
 />
 Use dark mode
 </label>
 </CurrentUserContext>
 </ThemeContext>
 )
}

function WelcomePanel({ children }) {
 const {currentUser} = useContext(CurrentUserContext);
 return (
 <Panel title="Welcome">
 {currentUser !== null ?
 <Greeting /> :
 <LoginForm />
 }
 </Panel>
 );
}

function Greeting() {
 const {currentUser} = useContext(CurrentUserContext);
 return (
 <p>You logged in as {currentUser.name}.</p>
 )
}

function LoginForm() {
 const {setCurrentUser} = useContext(CurrentUserContext);
 const [firstName, setFirstName] = useState('');
 const [lastName, setLastName] = useState('');
 const canLogin = firstName.trim() !== '' && lastName.trim() !== '';
 return (
 <>
 <label>
 First name{': '}
 <input
 required
 value={firstName}
 onChange={e => setFirstName(e.target.value)}
 />
 </label>
 <label>
 Last name{': '}
 <input
 required
 value={lastName}
 onChange={e => setLastName(e.target.value)}
 />
 </label>
 <Button
 disabled={!canLogin}
 onClick={() => {
 setCurrentUser({
 name: firstName + ' ' + lastName
 });
 }}
 >
 Log in
 </Button>
 {!canLogin && <i>Fill in both fields.</i>}
 </>
 );
}

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children, disabled, onClick }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button
 className={className}
 disabled={disabled}
 onClick={onClick}
 >
 {children}
 </button>
 );
}
```

```css
label {
 display: block;
}

.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
 margin-bottom: 10px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

<Solution />

#### Extracting providers to a component {/*extracting-providers-to-a-component*/}

As your app grows, it is expected that you'll have a "pyramid" of contexts closer to the root of your app. There is nothing wrong with that. However, if you dislike the nesting aesthetically, you can extract the providers into a single component. In this example, `MyProviders` hides the "plumbing" and renders the children passed to it inside the necessary providers. Note that the `theme` and `setTheme` state is needed in `MyApp` itself, so `MyApp` still owns that piece of the state.

<Sandpack>

```js
import { createContext, useContext, useState } from 'react';

const ThemeContext = createContext(null);
const CurrentUserContext = createContext(null);

export default function MyApp() {
 const [theme, setTheme] = useState('light');
 return (
 <MyProviders theme={theme} setTheme={setTheme}>
 <WelcomePanel />
 <label>
 <input
 type="checkbox"
 checked={theme === 'dark'}
 onChange={(e) => {
 setTheme(e.target.checked ? 'dark' : 'light')
 }}
 />
 Use dark mode
 </label>
 </MyProviders>
 );
}

function MyProviders({ children, theme, setTheme }) {
 const [currentUser, setCurrentUser] = useState(null);
 return (
 <ThemeContext value={theme}>
 <CurrentUserContext
 value={{
 currentUser,
 setCurrentUser
 }}
 >
 {children}
 </CurrentUserContext>
 </ThemeContext>
 );
}

function WelcomePanel({ children }) {
 const {currentUser} = useContext(CurrentUserContext);
 return (
 <Panel title="Welcome">
 {currentUser !== null ?
 <Greeting /> :
 <LoginForm />
 }
 </Panel>
 );
}

function Greeting() {
 const {currentUser} = useContext(CurrentUserContext);
 return (
 <p>You logged in as {currentUser.name}.</p>
 )
}

function LoginForm() {
 const {setCurrentUser} = useContext(CurrentUserContext);
 const [firstName, setFirstName] = useState('');
 const [lastName, setLastName] = useState('');
 const canLogin = firstName !== '' && lastName !== '';
 return (
 <>
 <label>
 First name{': '}
 <input
 required
 value={firstName}
 onChange={e => setFirstName(e.target.value)}
 />
 </label>
 <label>
 Last name{': '}
 <input
 required
 value={lastName}
 onChange={e => setLastName(e.target.value)}
 />
 </label>
 <Button
 disabled={!canLogin}
 onClick={() => {
 setCurrentUser({
 name: firstName + ' ' + lastName
 });
 }}
 >
 Log in
 </Button>
 {!canLogin && <i>Fill in both fields.</i>}
 </>
 );
}

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children, disabled, onClick }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button
 className={className}
 disabled={disabled}
 onClick={onClick}
 >
 {children}
 </button>
 );
}
```

```css
label {
 display: block;
}

.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
 margin-bottom: 10px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

<Solution />

#### Scaling up with context and a reducer {/*scaling-up-with-context-and-a-reducer*/}

In larger apps, it is common to combine context with a [reducer](/reference/react/useReducer) to extract the logic related to some state out of components. In this example, all the "wiring" is hidden in the `TasksContext.js`, which contains a reducer and two separate contexts.

Read a [full walkthrough](/learn/scaling-up-with-reducer-and-context) of this example.

<Sandpack>

```js src/App.js
import AddTask from './AddTask.js';
import TaskList from './TaskList.js';
import { TasksProvider } from './TasksContext.js';

export default function TaskApp() {
 return (
 <TasksProvider>
 <h1>Day off in Kyoto</h1>
 <AddTask />
 <TaskList />
 </TasksProvider>
 );
}
```

```js src/TasksContext.js
import { createContext, useContext, useReducer } from 'react';

const TasksContext = createContext(null);

const TasksDispatchContext = createContext(null);

export function TasksProvider({ children }) {
 const [tasks, dispatch] = useReducer(
 tasksReducer,
 initialTasks
 );

 return (
 <TasksContext value={tasks}>
 <TasksDispatchContext value={dispatch}>
 {children}
 </TasksDispatchContext>
 </TasksContext>
 );
}

export function useTasks() {
 return useContext(TasksContext);
}

export function useTasksDispatch() {
 return useContext(TasksDispatchContext);
}

function tasksReducer(tasks, action) {
 switch (action.type) {
 case 'added': {
 return [...tasks, {
 id: action.id,
 text: action.text,
 done: false
 }];
 }
 case 'changed': {
 return tasks.map(t => {
 if (t.id === action.task.id) {
 return action.task;
 } else {
 return t;
 }
 });
 }
 case 'deleted': {
 return tasks.filter(t => t.id !== action.id);
 }
 default: {
 throw Error('Unknown action: ' + action.type);
 }
 }
}

const initialTasks = [
 { id: 0, text: 'Philosopher’s Path', done: true },
 { id: 1, text: 'Visit the temple', done: false },
 { id: 2, text: 'Drink matcha', done: false }
];
```

```js src/AddTask.js
import { useState } from 'react';
import { useTasksDispatch } from './TasksContext.js';

export default function AddTask() {
 const [text, setText] = useState('');
 const dispatch = useTasksDispatch();
 return (
 <>
 <input
 placeholder="Add task"
 value={text}
 onChange={e => setText(e.target.value)}
 />
 <button onClick={() => {
 setText('');
 dispatch({
 type: 'added',
 id: nextId++,
 text: text,
 });
 }}>Add</button>
 </>
 );
}

let nextId = 3;
```

```js src/TaskList.js
import { useState } from 'react';
import { useTasks, useTasksDispatch } from './TasksContext.js';

export default function TaskList() {
 const tasks = useTasks();
 return (
 <ul>
 {tasks.map(task => (
 <li key={task.id}>
 <Task task={task} />
 </li>
 ))}
 </ul>
 );
}

function Task({ task }) {
 const [isEditing, setIsEditing] = useState(false);
 const dispatch = useTasksDispatch();
 let taskContent;
 if (isEditing) {
 taskContent = (
 <>
 <input
 value={task.text}
 onChange={e => {
 dispatch({
 type: 'changed',
 task: {
 ...task,
 text: e.target.value
 }
 });
 }} />
 <button onClick={() => setIsEditing(false)}>
 Save
 </button>
 </>
 );
 } else {
 taskContent = (
 <>
 {task.text}
 <button onClick={() => setIsEditing(true)}>
 Edit
 </button>
 </>
 );
 }
 return (
 <label>
 <input
 type="checkbox"
 checked={task.done}
 onChange={e => {
 dispatch({
 type: 'changed',
 task: {
 ...task,
 done: e.target.checked
 }
 });
 }}
 />
 {taskContent}
 <button onClick={() => {
 dispatch({
 type: 'deleted',
 id: task.id
 });
 }}>
 Delete
 </button>
 </label>
 );
}
```

```css
button { margin: 5px; }
li { list-style-type: none; }
ul, li { margin: 0; padding: 0; }
```

</Sandpack>

<Solution />

</Recipes>

---

### Specifying a fallback default value {/*specifying-a-fallback-default-value*/}

If React can't find any providers of that particular <CodeStep step={1}>context</CodeStep> in the parent tree, the context value returned by `useContext()` will be equal to the <CodeStep step={3}>default value</CodeStep> that you specified when you [created that context](/reference/react/createContext):

```js [[1, 1, "ThemeContext"], [3, 1, "null"]]
const ThemeContext = createContext(null);
```

The default value **never changes**. If you want to update context, use it with state as [described above.](#updating-data-passed-via-context)

Often, instead of `null`, there is some more meaningful value you can use as a default, for example:

```js [[1, 1, "ThemeContext"], [3, 1, "light"]]
const ThemeContext = createContext('light');
```

This way, if you accidentally render some component without a corresponding provider, it won't break. This also helps your components work well in a test environment without setting up a lot of providers in the tests.

In the example below, the "Toggle theme" button is always light because it's **outside any theme context provider** and the default context theme value is `'light'`. Try editing the default theme to be `'dark'`.

<Sandpack>

```js
import { createContext, useContext, useState } from 'react';

const ThemeContext = createContext('light');

export default function MyApp() {
 const [theme, setTheme] = useState('light');
 return (
 <>
 <ThemeContext value={theme}>
 <Form />
 </ThemeContext>
 <Button onClick={() => {
 setTheme(theme === 'dark' ? 'light' : 'dark');
 }}>
 Toggle theme
 </Button>
 </>
 )
}

function Form({ children }) {
 return (
 <Panel title="Welcome">
 <Button>Sign up</Button>
 <Button>Log in</Button>
 </Panel>
 );
}

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 <h1>{title}</h1>
 {children}
 </section>
 )
}

function Button({ children, onClick }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button className={className} onClick={onClick}>
 {children}
 </button>
 );
}
```

```css
.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
 margin-bottom: 10px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

---

### Overriding context for a part of the tree {/*overriding-context-for-a-part-of-the-tree*/}

You can override the context for a part of the tree by wrapping that part in a provider with a different value.

```js {3,5}
<ThemeContext value="dark">
 ...
 <ThemeContext value="light">
 <Footer />
 </ThemeContext>
 ...
</ThemeContext>
```

You can nest and override providers as many times as you need.

<Recipes titleText="Examples of overriding context">

#### Overriding a theme {/*overriding-a-theme*/}

Here, the button *inside* the `Footer` receives a different context value (`"light"`) than the buttons outside (`"dark"`).

<Sandpack>

```js
import { createContext, useContext } from 'react';

const ThemeContext = createContext(null);

export default function MyApp() {
 return (
 <ThemeContext value="dark">
 <Form />
 </ThemeContext>
 )
}

function Form() {
 return (
 <Panel title="Welcome">
 <Button>Sign up</Button>
 <Button>Log in</Button>
 <ThemeContext value="light">
 <Footer />
 </ThemeContext>
 </Panel>
 );
}

function Footer() {
 return (
 <footer>
 <Button>Settings</Button>
 </footer>
 );
}

function Panel({ title, children }) {
 const theme = useContext(ThemeContext);
 const className = 'panel-' + theme;
 return (
 <section className={className}>
 {title && <h1>{title}</h1>}
 {children}
 </section>
 )
}

function Button({ children }) {
 const theme = useContext(ThemeContext);
 const className = 'button-' + theme;
 return (
 <button className={className}>
 {children}
 </button>
 );
}
```

```css
footer {
 margin-top: 20px;
 border-top: 1px solid #aaa;
}

.panel-light,
.panel-dark {
 border: 1px solid black;
 border-radius: 4px;
 padding: 20px;
}
.panel-light {
 color: #222;
 background: #fff;
}

.panel-dark {
 color: #fff;
 background: rgb(23, 32, 42);
}

.button-light,
.button-dark {
 border: 1px solid #777;
 padding: 5px;
 margin-right: 10px;
 margin-top: 10px;
}

.button-dark {
 background: #222;
 color: #fff;
}

.button-light {
 background: #fff;
 color: #222;
}
```

</Sandpack>

<Solution />

#### Automatically nested headings {/*automatically-nested-headings*/}

You can "accumulate" information when you nest context providers. In this example, the `Section` component keeps track of the `LevelContext` which specifies the depth of the section nesting. It reads the `LevelContext` from the parent section, and provides the `LevelContext` number increased by one to its children. As a result, the `Heading` component can automatically decide which of the `<h1>`, `<h2>`, `<h3>`, ..., tags to use based on how many `Section` components it is nested inside of.

Read a [detailed walkthrough](/learn/passing-data-deeply-with-context) of this example.

<Sandpack>

```js
import Heading from './Heading.js';
import Section from './Section.js';

export default function Page() {
 return (
 <Section>
 <Heading>Title</Heading>
 <Section>
 <Heading>Heading</Heading>
 <Heading>Heading</Heading>
 <Heading>Heading</Heading>
 <Section>
 <Heading>Sub-heading</Heading>
 <Heading>Sub-heading</Heading>
 <Heading>Sub-heading</Heading>
 <Section>
 <Heading>Sub-sub-heading</Heading>
 <Heading>Sub-sub-heading</Heading>
 <Heading>Sub-sub-heading</Heading>
 </Section>
 </Section>
 </Section>
 </Section>
 );
}
```

```js src/Section.js
import { useContext } from 'react';
import { LevelContext } from './LevelContext.js';

export default function Section({ children }) {
 const level = useContext(LevelContext);
 return (
 <section className="section">
 <LevelContext value={level + 1}>
 {children}
 </LevelContext>
 </section>
 );
}
```

```js src/Heading.js
import { useContext } from 'react';
import { LevelContext } from './LevelContext.js';

export default function Heading({ children }) {
 const level = useContext(LevelContext);
 switch (level) {
 case 0:
 throw Error('Heading must be inside a Section!');
 case 1:
 return <h1>{children}</h1>;
 case 2:
 return <h2>{children}</h2>;
 case 3:
 return <h3>{children}</h3>;
 case 4:
 return <h4>{children}</h4>;
 case 5:
 return <h5>{children}</h5>;
 case 6:
 return <h6>{children}</h6>;
 default:
 throw Error('Unknown level: ' + level);
 }
}
```

```js src/LevelContext.js
import { createContext } from 'react';

export const LevelContext = createContext(0);
```

```css
.section {
 padding: 10px;
 margin: 5px;
 border-radius: 5px;
 border: 1px solid #aaa;
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Optimizing re-renders when passing objects and functions {/*optimizing-re-renders-when-passing-objects-and-functions*/}

You can pass any values via context, including objects and functions.

```js [[2, 10, "{ currentUser, login }"]]
function MyApp() {
 const [currentUser, setCurrentUser] = useState(null);

 function login(response) {
 storeCredentials(response.credentials);
 setCurrentUser(response.user);
 }

 return (
 <AuthContext value={{ currentUser, login }}>
 <Page />
 </AuthContext>
 );
}
```

Here, the <CodeStep step={2}>context value</CodeStep> is a JavaScript object with two properties, one of which is a function. Whenever `MyApp` re-renders (for example, on a route update), this will be a *different* object pointing at a *different* function, so React will also have to re-render all components deep in the tree that call `useContext(AuthContext)`.

In smaller apps, this is not a problem. However, there is no need to re-render them if the underlying data, like `currentUser`, has not changed. To help React take advantage of that fact, you may wrap the `login` function with [`useCallback`](/reference/react/useCallback) and wrap the object creation into [`useMemo`](/reference/react/useMemo). This is a performance optimization:

```js {6,9,11,14,17}
import { useCallback, useMemo } from 'react';

function MyApp() {
 const [currentUser, setCurrentUser] = useState(null);

 const login = useCallback((response) => {
 storeCredentials(response.credentials);
 setCurrentUser(response.user);
 }, []);

 const contextValue = useMemo(() => ({
 currentUser,
 login
 }), [currentUser, login]);

 return (
 <AuthContext value={contextValue}>
 <Page />
 </AuthContext>
 );
}
```

As a result of this change, even if `MyApp` needs to re-render, the components calling `useContext(AuthContext)` won't need to re-render unless `currentUser` has changed.

Read more about [`useMemo`](/reference/react/useMemo#skipping-re-rendering-of-components) and [`useCallback`.](/reference/react/useCallback#skipping-re-rendering-of-components)

---

## Troubleshooting {/*troubleshooting*/}

### My component doesn't see the value from my provider {/*my-component-doesnt-see-the-value-from-my-provider*/}

There are a few common ways that this can happen:

1. You're rendering `<SomeContext>` in the same component (or below) as where you're calling `useContext()`. Move `<SomeContext>` *above and outside* the component calling `useContext()`.
2. You may have forgotten to wrap your component with `<SomeContext>`, or you might have put it in a different part of the tree than you thought. Check whether the hierarchy is right using [React DevTools.](/learn/react-developer-tools)
3. You might be running into some build issue with your tooling that causes `SomeContext` as seen from the providing component and `SomeContext` as seen by the reading component to be two different objects. This can happen if you use symlinks, for example. You can verify this by assigning them to globals like `window.SomeContext1` and `window.SomeContext2` and then checking whether `window.SomeContext1 === window.SomeContext2` in the console. If they're not the same, fix that issue on the build tool level.

### I am always getting `undefined` from my context although the default value is different {/*i-am-always-getting-undefined-from-my-context-although-the-default-value-is-different*/}

You might have a provider without a `value` in the tree:

```js {1,2}
// 🚩 Doesn't work: no value prop
<ThemeContext>
 <Button />
</ThemeContext>
```

If you forget to specify `value`, it's like passing `value={undefined}`.

You may have also mistakingly used a different prop name by mistake:

```js {1,2}
// 🚩 Doesn't work: prop should be called "value"
<ThemeContext theme={theme}>
 <Button />
</ThemeContext>
```

In both of these cases you should see a warning from React in the console. To fix them, call the prop `value`:

```js {1,2}
// ✅ Passing the value prop
<ThemeContext value={theme}>
 <Button />
</ThemeContext>
```

Note that the [default value from your `createContext(defaultValue)` call](#specifying-a-fallback-default-value) is only used **if there is no matching provider above at all.** If there is a `<SomeContext value={undefined}>` component somewhere in the parent tree, the component calling `useContext(SomeContext)` *will* receive `undefined` as the context value.

---
title: useDebugValue
---

<Intro>

`useDebugValue` is a React Hook that lets you add a label to a custom Hook in [React DevTools.](/learn/react-developer-tools)

```js
useDebugValue(value, format?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useDebugValue(value, format?)` {/*usedebugvalue*/}

Call `useDebugValue` at the top level of your [custom Hook](/learn/reusing-logic-with-custom-hooks) to display a readable debug value:

```js
import { useDebugValue } from 'react';

function useOnlineStatus() {
 // ...
 useDebugValue(isOnline ? 'Online' : 'Offline');
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `value`: The value you want to display in React DevTools. It can have any type.
* **optional** `format`: A formatting function. When the component is inspected, React DevTools will call the formatting function with the `value` as the argument, and then display the returned formatted value (which may have any type). If you don't specify the formatting function, the original `value` itself will be displayed.

#### Returns {/*returns*/}

`useDebugValue` does not return anything.

## Usage {/*usage*/}

### Adding a label to a custom Hook {/*adding-a-label-to-a-custom-hook*/}

Call `useDebugValue` at the top level of your [custom Hook](/learn/reusing-logic-with-custom-hooks) to display a readable <CodeStep step={1}>debug value</CodeStep> for [React DevTools.](/learn/react-developer-tools)

```js [[1, 5, "isOnline ? 'Online' : 'Offline'"]]
import { useDebugValue } from 'react';

function useOnlineStatus() {
 // ...
 useDebugValue(isOnline ? 'Online' : 'Offline');
 // ...
}
```

This gives components calling `useOnlineStatus` a label like `OnlineStatus: "Online"` when you inspect them:

![A screenshot of React DevTools showing the debug value](/images/docs/react-devtools-usedebugvalue.png)

Without the `useDebugValue` call, only the underlying data (in this example, `true`) would be displayed.

<Sandpack>

```js
import { useOnlineStatus } from './useOnlineStatus.js';

function StatusBar() {
 const isOnline = useOnlineStatus();
 return <h1>{isOnline ? '✅ Online' : '❌ Disconnected'}</h1>;
}

export default function App() {
 return <StatusBar />;
}
```

```js src/useOnlineStatus.js active
import { useSyncExternalStore, useDebugValue } from 'react';

export function useOnlineStatus() {
 const isOnline = useSyncExternalStore(subscribe, () => navigator.onLine, () => true);
 useDebugValue(isOnline ? 'Online' : 'Offline');
 return isOnline;
}

function subscribe(callback) {
 window.addEventListener('online', callback);
 window.addEventListener('offline', callback);
 return () => {
 window.removeEventListener('online', callback);
 window.removeEventListener('offline', callback);
 };
}
```

</Sandpack>

<Note>

Don't add debug values to every custom Hook. It's most valuable for custom Hooks that are part of shared libraries and that have a complex internal data structure that's difficult to inspect.

</Note>

---

### Deferring formatting of a debug value {/*deferring-formatting-of-a-debug-value*/}

You can also pass a formatting function as the second argument to `useDebugValue`:

```js [[1, 1, "date", 18], [2, 1, "date.toDateString()"]]
useDebugValue(date, date => date.toDateString());
```

Your formatting function will receive the <CodeStep step={1}>debug value</CodeStep> as a parameter and should return a <CodeStep step={2}>formatted display value</CodeStep>. When your component is inspected, React DevTools will call this function and display its result.

This lets you avoid running potentially expensive formatting logic unless the component is actually inspected. For example, if `date` is a Date value, this avoids calling `toDateString()` on it for every render.

---
title: useDeferredValue
---

<Intro>

`useDeferredValue` is a React Hook that lets you defer updating a part of the UI.

```js
const deferredValue = useDeferredValue(value)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useDeferredValue(value, initialValue?)` {/*usedeferredvalue*/}

Call `useDeferredValue` at the top level of your component to get a deferred version of that value.

```js
import { useState, useDeferredValue } from 'react';

function SearchPage() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `value`: The value you want to defer. It can have any type.
* **optional** `initialValue`: A value to use during the initial render of a component. If this option is omitted, `useDeferredValue` will not defer during the initial render, because there's no previous version of `value` that it can render instead.

#### Returns {/*returns*/}

- `currentValue`: During the initial render, the returned deferred value will be the `initialValue`, or the same as the value you provided. During updates, React will first attempt a re-render with the old value (so it will return the old value), and then try another re-render in the background with the new value (so it will return the updated value).

#### Caveats {/*caveats*/}

- When an update is inside a Transition, `useDeferredValue` always returns the new `value` and does not spawn a deferred render, since the update is already deferred.

- The values you pass to `useDeferredValue` should either be primitive values (like strings and numbers) or objects created outside of rendering. If you create a new object during rendering and immediately pass it to `useDeferredValue`, it will be different on every render, causing unnecessary background re-renders.

- When `useDeferredValue` receives a different value (compared with [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is)), in addition to the current render (when it still uses the previous value), it schedules a re-render in the background with the new value. The background re-render is interruptible: if there's another update to the `value`, React will restart the background re-render from scratch. For example, if the user is typing into an input faster than a chart receiving its deferred value can re-render, the chart will only re-render after the user stops typing.

- `useDeferredValue` is integrated with [`<Suspense>`.](/reference/react/Suspense) If the background update caused by a new value suspends the UI, the user will not see the fallback. They will see the old deferred value until the data loads.

- `useDeferredValue` does not by itself prevent extra network requests.

- There is no fixed delay caused by `useDeferredValue` itself. As soon as React finishes the original re-render, React will immediately start working on the background re-render with the new deferred value. Any updates caused by events (like typing) will interrupt the background re-render and get prioritized over it.

- The background re-render caused by `useDeferredValue` does not fire Effects until it's committed to the screen. If the background re-render suspends, its Effects will run after the data loads and the UI updates.

---

## Usage {/*usage*/}

### Showing stale content while fresh content is loading {/*showing-stale-content-while-fresh-content-is-loading*/}

Call `useDeferredValue` at the top level of your component to defer updating some part of your UI.

```js [[1, 5, "query"], [2, 5, "deferredQuery"]]
import { useState, useDeferredValue } from 'react';

function SearchPage() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 // ...
}
```

During the initial render, the <CodeStep step={2}>deferred value</CodeStep> will be the same as the <CodeStep step={1}>value</CodeStep> you provided.

During updates, the <CodeStep step={2}>deferred value</CodeStep> will "lag behind" the latest <CodeStep step={1}>value</CodeStep>. In particular, React will first re-render *without* updating the deferred value, and then try to re-render with the newly received value in the background.

**Let's walk through an example to see when this is useful.**

<Note>

This example assumes you use a data source that [activates a Suspense boundary](/reference/react/Suspense#what-activates-a-suspense-boundary), such as a Promise you read with [`use`](/reference/react/use).

[Learn more about Suspense.](/reference/react/Suspense)

</Note>

In this example, the `SearchResults` component [suspends](/reference/react/Suspense#displaying-a-fallback-while-content-is-loading) while fetching the search results. Try typing `"a"`, waiting for the results, and then editing it to `"ab"`. The results for `"a"` get replaced by the loading fallback.

<Sandpack>

```js src/App.js
import { Suspense, useState } from 'react';
import SearchResults from './SearchResults.js';

export default function App() {
 const [query, setQuery] = useState('');
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <SearchResults query={query} />
 </Suspense>
 </>
 );
}
```

```js src/SearchResults.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function SearchResults({ query }) {
 if (query === '') {
 return null;
 }
 const albums = use(fetchData(`/search?q=${query}`));
 if (albums.length === 0) {
 return <p>No matches for <i>"{query}"</i></p>;
 }
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/search?q=')) {
 return await getSearchResults(url.slice('/search?q='.length));
 } else {
 throw Error('Not implemented');
 }
}

async function getSearchResults(query) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 const allAlbums = [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];

 const lowerQuery = query.trim().toLowerCase();
 return allAlbums.filter(album => {
 const lowerTitle = album.title.toLowerCase();
 return (
 lowerTitle.startsWith(lowerQuery) ||
 lowerTitle.indexOf(' ' + lowerQuery) !== -1
 )
 });
}
```

```css
input { margin: 10px; }
```

</Sandpack>

A common alternative UI pattern is to *defer* updating the list of results and to keep showing the previous results until the new results are ready. Call `useDeferredValue` to pass a deferred version of the query down:

```js {3,11}
export default function App() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <SearchResults query={deferredQuery} />
 </Suspense>
 </>
 );
}
```

The `query` will update immediately, so the input will display the new value. However, the `deferredQuery` will keep its previous value until the data has loaded, so `SearchResults` will show the stale results for a bit.

Enter `"a"` in the example below, wait for the results to load, and then edit the input to `"ab"`. Notice how instead of the Suspense fallback, you now see the stale result list until the new results have loaded:

<Sandpack>

```js src/App.js
import { Suspense, useState, useDeferredValue } from 'react';
import SearchResults from './SearchResults.js';

export default function App() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <SearchResults query={deferredQuery} />
 </Suspense>
 </>
 );
}
```

```js src/SearchResults.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function SearchResults({ query }) {
 if (query === '') {
 return null;
 }
 const albums = use(fetchData(`/search?q=${query}`));
 if (albums.length === 0) {
 return <p>No matches for <i>"{query}"</i></p>;
 }
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/search?q=')) {
 return await getSearchResults(url.slice('/search?q='.length));
 } else {
 throw Error('Not implemented');
 }
}

async function getSearchResults(query) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 const allAlbums = [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];

 const lowerQuery = query.trim().toLowerCase();
 return allAlbums.filter(album => {
 const lowerTitle = album.title.toLowerCase();
 return (
 lowerTitle.startsWith(lowerQuery) ||
 lowerTitle.indexOf(' ' + lowerQuery) !== -1
 )
 });
}
```

```css
input { margin: 10px; }
```

</Sandpack>

<DeepDive>

#### How does deferring a value work under the hood? {/*how-does-deferring-a-value-work-under-the-hood*/}

You can think of it as happening in two steps:

1. **First, React re-renders with the new `query` (`"ab"`) but with the old `deferredQuery` (still `"a"`).** The `deferredQuery` value, which you pass to the result list, is *deferred:* it "lags behind" the `query` value.

2. **In the background, React tries to re-render with *both* `query` and `deferredQuery` updated to `"ab"`.** If this re-render completes, React will show it on the screen. However, if it suspends (the results for `"ab"` have not loaded yet), React will abandon this rendering attempt, and retry this re-render again after the data has loaded. The user will keep seeing the stale deferred value until the data is ready.

The deferred "background" rendering is interruptible. For example, if you type into the input again, React will abandon it and restart with the new value. React will always use the latest provided value.

Note that there is still a network request per each keystroke. What's being deferred here is displaying results (until they're ready), not the network requests themselves. Even if the user continues typing, responses for each keystroke get cached, so pressing Backspace is instant and doesn't fetch again.

</DeepDive>

---

### Indicating that the content is stale {/*indicating-that-the-content-is-stale*/}

In the example above, there is no indication that the result list for the latest query is still loading. This can be confusing to the user if the new results take a while to load. To make it more obvious to the user that the result list does not match the latest query, you can add a visual indication when the stale result list is displayed:

```js {2}
<div style={{
 opacity: query !== deferredQuery ? 0.5 : 1,
}}>
 <SearchResults query={deferredQuery} />
</div>
```

With this change, as soon as you start typing, the stale result list gets slightly dimmed until the new result list loads. You can also add a CSS transition to delay dimming so that it feels gradual, like in the example below:

<Sandpack>

```js src/App.js
import { Suspense, useState, useDeferredValue } from 'react';
import SearchResults from './SearchResults.js';

export default function App() {
 const [query, setQuery] = useState('');
 const deferredQuery = useDeferredValue(query);
 const isStale = query !== deferredQuery;
 return (
 <>
 <label>
 Search albums:
 <input value={query} onChange={e => setQuery(e.target.value)} />
 </label>
 <Suspense fallback={<h2>Loading...</h2>}>
 <div style={{
 opacity: isStale ? 0.5 : 1,
 transition: isStale ? 'opacity 0.2s 0.2s linear' : 'opacity 0s 0s linear'
 }}>
 <SearchResults query={deferredQuery} />
 </div>
 </Suspense>
 </>
 );
}
```

```js src/SearchResults.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function SearchResults({ query }) {
 if (query === '') {
 return null;
 }
 const albums = use(fetchData(`/search?q=${query}`));
 if (albums.length === 0) {
 return <p>No matches for <i>"{query}"</i></p>;
 }
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/search?q=')) {
 return await getSearchResults(url.slice('/search?q='.length));
 } else {
 throw Error('Not implemented');
 }
}

async function getSearchResults(query) {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });

 const allAlbums = [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];

 const lowerQuery = query.trim().toLowerCase();
 return allAlbums.filter(album => {
 const lowerTitle = album.title.toLowerCase();
 return (
 lowerTitle.startsWith(lowerQuery) ||
 lowerTitle.indexOf(' ' + lowerQuery) !== -1
 )
 });
}
```

```css
input { margin: 10px; }
```

</Sandpack>

---

### Deferring re-rendering for a part of the UI {/*deferring-re-rendering-for-a-part-of-the-ui*/}

You can also apply `useDeferredValue` as a performance optimization. It is useful when a part of your UI is slow to re-render, there's no easy way to optimize it, and you want to prevent it from blocking the rest of the UI.

Imagine you have a text field and a component (like a chart or a long list) that re-renders on every keystroke:

```js
function App() {
 const [text, setText] = useState('');
 return (
 <>
 <input value={text} onChange={e => setText(e.target.value)} />
 <SlowList text={text} />
 </>
 );
}
```

First, optimize `SlowList` to skip re-rendering when its props are the same. To do this, [wrap it in `memo`:](/reference/react/memo#skipping-re-rendering-when-props-are-unchanged)

```js {1,3}
const SlowList = memo(function SlowList({ text }) {
 // ...
});
```

However, this only helps if the `SlowList` props are *the same* as during the previous render. The problem you're facing now is that it's slow when they're *different,* and when you actually need to show different visual output.

Concretely, the main performance problem is that whenever you type into the input, the `SlowList` receives new props, and re-rendering its entire tree makes the typing feel janky. In this case, `useDeferredValue` lets you prioritize updating the input (which must be fast) over updating the result list (which is allowed to be slower):

```js {3,7}
function App() {
 const [text, setText] = useState('');
 const deferredText = useDeferredValue(text);
 return (
 <>
 <input value={text} onChange={e => setText(e.target.value)} />
 <SlowList text={deferredText} />
 </>
 );
}
```

This does not make re-rendering of the `SlowList` faster. However, it tells React that re-rendering the list can be deprioritized so that it doesn't block the keystrokes. The list will "lag behind" the input and then "catch up". Like before, React will attempt to update the list as soon as possible, but will not block the user from typing.

<Recipes titleText="The difference between useDeferredValue and unoptimized re-rendering" titleId="examples">

#### Deferred re-rendering of the list {/*deferred-re-rendering-of-the-list*/}

In this example, each item in the `SlowList` component is **artificially slowed down** so that you can see how `useDeferredValue` lets you keep the input responsive. Type into the input and notice that typing feels snappy while the list "lags behind" it.

<Sandpack>

```js
import { useState, useDeferredValue } from 'react';
import SlowList from './SlowList.js';

export default function App() {
 const [text, setText] = useState('');
 const deferredText = useDeferredValue(text);
 return (
 <>
 <input value={text} onChange={e => setText(e.target.value)} />
 <SlowList text={deferredText} />
 </>
 );
}
```

```js {expectedErrors: {'react-compiler': [19, 20]}} src/SlowList.js
import { memo } from 'react';

const SlowList = memo(function SlowList({ text }) {
 // Log once. The actual slowdown is inside SlowItem.
 console.log('[ARTIFICIALLY SLOW] Rendering 250 <SlowItem />');

 let items = [];
 for (let i = 0; i < 250; i++) {
 items.push(<SlowItem key={i} text={text} />);
 }
 return (
 <ul className="items">
 {items}
 </ul>
 );
});

function SlowItem({ text }) {
 let startTime = performance.now();
 while (performance.now() - startTime < 1) {
 // Do nothing for 1 ms per item to emulate extremely slow code
 }

 return (
 <li className="item">
 Text: {text}
 </li>
 )
}

export default SlowList;
```

```css
.items {
 padding: 0;
 max-height: 300px;
 overflow: auto;
}

.item {
 list-style: none;
 display: block;
 height: 40px;
 padding: 5px;
 margin-top: 10px;
 border-radius: 4px;
 border: 1px solid #aaa;
}
```

</Sandpack>

<Solution />

#### Unoptimized re-rendering of the list {/*unoptimized-re-rendering-of-the-list*/}

In this example, each item in the `SlowList` component is **artificially slowed down**, but there is no `useDeferredValue`.

Notice how typing into the input feels very janky. This is because without `useDeferredValue`, each keystroke forces the entire list to re-render immediately in a non-interruptible way.

<Sandpack>

```js
import { useState } from 'react';
import SlowList from './SlowList.js';

export default function App() {
 const [text, setText] = useState('');
 return (
 <>
 <input value={text} onChange={e => setText(e.target.value)} />
 <SlowList text={text} />
 </>
 );
}
```

```js {expectedErrors: {'react-compiler': [19, 20]}} src/SlowList.js
import { memo } from 'react';

const SlowList = memo(function SlowList({ text }) {
 // Log once. The actual slowdown is inside SlowItem.
 console.log('[ARTIFICIALLY SLOW] Rendering 250 <SlowItem />');

 let items = [];
 for (let i = 0; i < 250; i++) {
 items.push(<SlowItem key={i} text={text} />);
 }
 return (
 <ul className="items">
 {items}
 </ul>
 );
});

function SlowItem({ text }) {
 let startTime = performance.now();
 while (performance.now() - startTime < 1) {
 // Do nothing for 1 ms per item to emulate extremely slow code
 }

 return (
 <li className="item">
 Text: {text}
 </li>
 )
}

export default SlowList;
```

```css
.items {
 padding: 0;
 max-height: 300px;
 overflow: auto;
}

.item {
 list-style: none;
 display: block;
 height: 40px;
 padding: 5px;
 margin-top: 10px;
 border-radius: 4px;
 border: 1px solid #aaa;
}
```

</Sandpack>

<Solution />

</Recipes>

<Pitfall>

This optimization requires `SlowList` to be wrapped in [`memo`.](/reference/react/memo) This is because whenever the `text` changes, React needs to be able to re-render the parent component quickly. During that re-render, `deferredText` still has its previous value, so `SlowList` is able to skip re-rendering (its props have not changed). Without [`memo`,](/reference/react/memo) it would have to re-render anyway, defeating the point of the optimization.

</Pitfall>

<DeepDive>

#### How is deferring a value different from debouncing and throttling? {/*how-is-deferring-a-value-different-from-debouncing-and-throttling*/}

There are two common optimization techniques you might have used before in this scenario:

- *Debouncing* means you'd wait for the user to stop typing (e.g. for a second) before updating the list.
- *Throttling* means you'd update the list every once in a while (e.g. at most once a second).

While these techniques are helpful in some cases, `useDeferredValue` is better suited to optimizing rendering because it is deeply integrated with React itself and adapts to the user's device.

Unlike debouncing or throttling, it doesn't require choosing any fixed delay. If the user's device is fast (e.g. powerful laptop), the deferred re-render would happen almost immediately and wouldn't be noticeable. If the user's device is slow, the list would "lag behind" the input proportionally to how slow the device is.

Also, unlike with debouncing or throttling, deferred re-renders done by `useDeferredValue` are interruptible by default. This means that if React is in the middle of re-rendering a large list, but the user makes another keystroke, React will abandon that re-render, handle the keystroke, and then start rendering in the background again. By contrast, debouncing and throttling still produce a janky experience because they're *blocking:* they merely postpone the moment when rendering blocks the keystroke.

If the work you're optimizing doesn't happen during rendering, debouncing and throttling are still useful. For example, they can let you fire fewer network requests. You can also use these techniques together.

</DeepDive>

---
title: useEffect
---

<Intro>

`useEffect` is a React Hook that lets you [synchronize a component with an external system.](/learn/synchronizing-with-effects)

```js
useEffect(setup, dependencies?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useEffect(setup, dependencies?)` {/*useeffect*/}

Call `useEffect` at the top level of your component to declare an Effect:

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [serverUrl, roomId]);
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `setup`: The function with your Effect's logic. Your setup function may also optionally return a *cleanup* function. When your [component commits](/learn/render-and-commit#step-3-react-commits-changes-to-the-dom), React will run your setup function. After every commit with changed dependencies, React will first run the cleanup function (if you provided it) with the old values, and then run your setup function with the new values. After your component is removed from the DOM, React will run your cleanup function.

* **optional** `dependencies`: The list of all reactive values referenced inside of the `setup` code. Reactive values include props, state, and all the variables and functions declared directly inside your component body. If your linter is [configured for React](/learn/editor-setup#linting), it will verify that every reactive value is correctly specified as a dependency. The list of dependencies must have a constant number of items and be written inline like `[dep1, dep2, dep3]`. React will compare each dependency with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. If you omit this argument, your Effect will re-run after every commit of the component. [See the difference between passing an array of dependencies, an empty array, and no dependencies at all.](#examples-dependencies)

#### Returns {/*returns*/}

`useEffect` returns `undefined`.

#### Caveats {/*caveats*/}

* `useEffect` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.

* If you're **not trying to synchronize with some external system,** [you probably don't need an Effect.](/learn/you-might-not-need-an-effect)

* When Strict Mode is on, React will **run one extra development-only setup+cleanup cycle** before the first real setup. This is a stress-test that ensures that your cleanup logic "mirrors" your setup logic and that it stops or undoes whatever the setup is doing. If this causes a problem, [implement the cleanup function.](/learn/synchronizing-with-effects#how-to-handle-the-effect-firing-twice-in-development)

* If some of your dependencies are objects or functions defined inside the component, there is a risk that they will **cause the Effect to re-run more often than needed.** To fix this, remove unnecessary [object](#removing-unnecessary-object-dependencies) and [function](#removing-unnecessary-function-dependencies) dependencies. You can also [extract state updates](#updating-state-based-on-previous-state-from-an-effect) and [non-reactive logic](#reading-the-latest-props-and-state-from-an-effect) outside of your Effect.

* If your Effect wasn't caused by an interaction (like a click), React will generally let the browser **paint the updated screen first before running your Effect.** If your Effect is doing something visual (for example, positioning a tooltip), and the delay is noticeable (for example, it flickers), replace `useEffect` with [`useLayoutEffect`.](/reference/react/useLayoutEffect)

* If your Effect is caused by an interaction (like a click), **React may run your Effect before the browser paints the updated screen**. This ensures that the result of the Effect can be observed by the event system. Usually, this works as expected. However, if you must defer the work until after paint, such as an `alert()`, you can use `setTimeout`. See [reactwg/react-18/128](https://github.com/reactwg/react-18/discussions/128) for more information.

* Even if your Effect was caused by an interaction (like a click), **React may allow the browser to repaint the screen before processing the state updates inside your Effect.** Usually, this works as expected. However, if you must block the browser from repainting the screen, you need to replace `useEffect` with [`useLayoutEffect`.](/reference/react/useLayoutEffect)

* Effects **only run on the client.** They don't run during server rendering.

---

## Usage {/*usage*/}

### Connecting to an external system {/*connecting-to-an-external-system*/}

Some components need to stay connected to the network, some browser API, or a third-party library, while they are displayed on the page. These systems aren't controlled by React, so they are called *external.*

To [connect your component to some external system,](/learn/synchronizing-with-effects) call `useEffect` at the top level of your component:

```js [[1, 8, "const connection = createConnection(serverUrl, roomId);"], [1, 9, "connection.connect();"], [2, 11, "connection.disconnect();"], [3, 13, "[serverUrl, roomId]"]]
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [serverUrl, roomId]);
 // ...
}
```

You need to pass two arguments to `useEffect`:

1. A *setup function* with <CodeStep step={1}>setup code</CodeStep> that connects to that system.
 - It should return a *cleanup function* with <CodeStep step={2}>cleanup code</CodeStep> that disconnects from that system.
2. A <CodeStep step={3}>list of dependencies</CodeStep> including every value from your component used inside of those functions.

**React calls your setup and cleanup functions whenever it's necessary, which may happen multiple times:**

1. Your <CodeStep step={1}>setup code</CodeStep> runs when your component is added to the page *(mounts)*.
2. After every commit of your component where the <CodeStep step={3}>dependencies</CodeStep> have changed:
 - First, your <CodeStep step={2}>cleanup code</CodeStep> runs with the old props and state.
 - Then, your <CodeStep step={1}>setup code</CodeStep> runs with the new props and state.
3. Your <CodeStep step={2}>cleanup code</CodeStep> runs one final time after your component is removed from the page *(unmounts).*

**Let's illustrate this sequence for the example above.**

When the `ChatRoom` component above gets added to the page, it will connect to the chat room with the initial `serverUrl` and `roomId`. If either `serverUrl` or `roomId` change as a result of a commit (say, if the user picks a different chat room in a dropdown), your Effect will *disconnect from the previous room, and connect to the next one.* When the `ChatRoom` component is removed from the page, your Effect will disconnect one last time.

**To [help you find bugs,](/learn/synchronizing-with-effects#step-3-add-cleanup-if-needed) in development React runs <CodeStep step={1}>setup</CodeStep> and <CodeStep step={2}>cleanup</CodeStep> one extra time before the <CodeStep step={1}>setup</CodeStep>.** This is a stress-test that verifies your Effect's logic is implemented correctly. If this causes visible issues, your cleanup function is missing some logic. The cleanup function should stop or undo whatever the setup function was doing. The rule of thumb is that the user shouldn't be able to distinguish between the setup being called once (as in production) and a *setup* → *cleanup* → *setup* sequence (as in development). [See common solutions.](/learn/synchronizing-with-effects#how-to-handle-the-effect-firing-twice-in-development)

**Try to [write every Effect as an independent process](/learn/lifecycle-of-reactive-effects#each-effect-represents-a-separate-synchronization-process) and [think about a single setup/cleanup cycle at a time.](/learn/lifecycle-of-reactive-effects#thinking-from-the-effects-perspective)** It shouldn't matter whether your component is mounting, updating, or unmounting. When your cleanup logic correctly "mirrors" the setup logic, your Effect is resilient to running setup and cleanup as often as needed.

<Note>

An Effect lets you [keep your component synchronized](/learn/synchronizing-with-effects) with some external system (like a chat service). Here, *external system* means any piece of code that's not controlled by React, such as:

* A timer managed with <CodeStep step={1}>[`setInterval()`](https://developer.mozilla.org/en-US/docs/Web/API/setInterval)</CodeStep> and <CodeStep step={2}>[`clearInterval()`](https://developer.mozilla.org/en-US/docs/Web/API/clearInterval)</CodeStep>.
* An event subscription using <CodeStep step={1}>[`window.addEventListener()`](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/addEventListener)</CodeStep> and <CodeStep step={2}>[`window.removeEventListener()`](https://developer.mozilla.org/en-US/docs/Web/API/EventTarget/removeEventListener)</CodeStep>.
* A third-party animation library with an API like <CodeStep step={1}>`animation.start()`</CodeStep> and <CodeStep step={2}>`animation.reset()`</CodeStep>.

**If you're not connecting to any external system, [you probably don't need an Effect.](/learn/you-might-not-need-an-effect)**

</Note>

<Recipes titleText="Examples of connecting to an external system" titleId="examples-connecting">

#### Connecting to a chat server {/*connecting-to-a-chat-server*/}

In this example, the `ChatRoom` component uses an Effect to stay connected to an external system defined in `chat.js`. Press "Open chat" to make the `ChatRoom` component appear. This sandbox runs in development mode, so there is an extra connect-and-disconnect cycle, as [explained here.](/learn/synchronizing-with-effects#step-3-add-cleanup-if-needed) Try changing the `roomId` and `serverUrl` using the dropdown and the input, and see how the Effect re-connects to the chat. Press "Close chat" to see the Effect disconnect one last time.

<Sandpack>

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [roomId, serverUrl]);

 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={serverUrl}
 onChange={e => setServerUrl(e.target.value)}
 />
 </label>
 <h1>Welcome to the {roomId} room!</h1>
 </>
 );
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

<Solution />

#### Listening to a global browser event {/*listening-to-a-global-browser-event*/}

In this example, the external system is the browser DOM itself. Normally, you'd specify event listeners with JSX, but you can't listen to the global [`window`](https://developer.mozilla.org/en-US/docs/Web/API/Window) object this way. An Effect lets you connect to the `window` object and listen to its events. Listening to the `pointermove` event lets you track the cursor (or finger) position and update the red dot to move with it.

<Sandpack>

```js
import { useState, useEffect } from 'react';

export default function App() {
 const [position, setPosition] = useState({ x: 0, y: 0 });

 useEffect(() => {
 function handleMove(e) {
 setPosition({ x: e.clientX, y: e.clientY });
 }
 window.addEventListener('pointermove', handleMove);
 return () => {
 window.removeEventListener('pointermove', handleMove);
 };
 }, []);

 return (
 <div style={{
 position: 'absolute',
 backgroundColor: 'pink',
 borderRadius: '50%',
 opacity: 0.6,
 transform: `translate(${position.x}px, ${position.y}px)`,
 pointerEvents: 'none',
 left: -20,
 top: -20,
 width: 40,
 height: 40,
 }} />
 );
}
```

```css
body {
 min-height: 300px;
}
```

</Sandpack>

<Solution />

#### Triggering an animation {/*triggering-an-animation*/}

In this example, the external system is the animation library in `animation.js`. It provides a JavaScript class called `FadeInAnimation` that takes a DOM node as an argument and exposes `start()` and `stop()` methods to control the animation. This component [uses a ref](/learn/manipulating-the-dom-with-refs) to access the underlying DOM node. The Effect reads the DOM node from the ref and automatically starts the animation for that node when the component appears.

<Sandpack>

```js
import { useState, useEffect, useRef } from 'react';
import { FadeInAnimation } from './animation.js';

function Welcome() {
 const ref = useRef(null);

 useEffect(() => {
 const animation = new FadeInAnimation(ref.current);
 animation.start(1000);
 return () => {
 animation.stop();
 };
 }, []);

 return (
 <h1
 ref={ref}
 style={{
 opacity: 0,
 color: 'white',
 padding: 50,
 textAlign: 'center',
 fontSize: 50,
 backgroundImage: 'radial-gradient(circle, rgba(63,94,251,1) 0%, rgba(252,70,107,1) 100%)'
 }}
 >
 Welcome
 </h1>
 );
}

export default function App() {
 const [show, setShow] = useState(false);
 return (
 <>
 <button onClick={() => setShow(!show)}>
 {show ? 'Remove' : 'Show'}
 </button>
 <hr />
 {show && <Welcome />}
 </>
 );
}
```

```js src/animation.js
export class FadeInAnimation {
 constructor(node) {
 this.node = node;
 }
 start(duration) {
 this.duration = duration;
 if (this.duration === 0) {
 // Jump to end immediately
 this.onProgress(1);
 } else {
 this.onProgress(0);
 // Start animating
 this.startTime = performance.now();
 this.frameId = requestAnimationFrame(() => this.onFrame());
 }
 }
 onFrame() {
 const timePassed = performance.now() - this.startTime;
 const progress = Math.min(timePassed / this.duration, 1);
 this.onProgress(progress);
 if (progress < 1) {
 // We still have more frames to paint
 this.frameId = requestAnimationFrame(() => this.onFrame());
 }
 }
 onProgress(progress) {
 this.node.style.opacity = progress;
 }
 stop() {
 cancelAnimationFrame(this.frameId);
 this.startTime = null;
 this.frameId = null;
 this.duration = 0;
 }
}
```

```css
label, button { display: block; margin-bottom: 20px; }
html, body { min-height: 300px; }
```

</Sandpack>

<Solution />

#### Controlling a modal dialog {/*controlling-a-modal-dialog*/}

In this example, the external system is the browser DOM. The `ModalDialog` component renders a [`<dialog>`](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/dialog) element. It uses an Effect to synchronize the `isOpen` prop to the [`showModal()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLDialogElement/showModal) and [`close()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLDialogElement/close) method calls.

<Sandpack>

```js
import { useState } from 'react';
import ModalDialog from './ModalDialog.js';

export default function App() {
 const [show, setShow] = useState(false);
 return (
 <>
 <button onClick={() => setShow(true)}>
 Open dialog
 </button>
 <ModalDialog isOpen={show}>
 Hello there!
 <br />
 <button onClick={() => {
 setShow(false);
 }}>Close</button>
 </ModalDialog>
 </>
 );
}
```

```js src/ModalDialog.js active
import { useEffect, useRef } from 'react';

export default function ModalDialog({ isOpen, children }) {
 const ref = useRef();

 useEffect(() => {
 if (!isOpen) {
 return;
 }
 const dialog = ref.current;
 dialog.showModal();
 return () => {
 dialog.close();
 };
 }, [isOpen]);

 return <dialog ref={ref}>{children}</dialog>;
}
```

```css
body {
 min-height: 300px;
}
```

</Sandpack>

<Solution />

#### Tracking element visibility {/*tracking-element-visibility*/}

In this example, the external system is again the browser DOM. The `App` component displays a long list, then a `Box` component, and then another long list. Scroll the list down. Notice that when all of the `Box` component is fully visible in the viewport, the background color changes to black. To implement this, the `Box` component uses an Effect to manage an [`IntersectionObserver`](https://developer.mozilla.org/en-US/docs/Web/API/Intersection_Observer_API). This browser API notifies you when the DOM element is visible in the viewport.

<Sandpack>

```js
import Box from './Box.js';

export default function App() {
 return (
 <>
 <LongSection />
 <Box />
 <LongSection />
 <Box />
 <LongSection />
 </>
 );
}

function LongSection() {
 const items = [];
 for (let i = 0; i < 50; i++) {
 items.push(<li key={i}>Item #{i} (keep scrolling)</li>);
 }
 return <ul>{items}</ul>
}
```

```js src/Box.js active
import { useRef, useEffect } from 'react';

export default function Box() {
 const ref = useRef(null);

 useEffect(() => {
 const div = ref.current;
 const observer = new IntersectionObserver(entries => {
 const entry = entries[0];
 if (entry.isIntersecting) {
 document.body.style.backgroundColor = 'black';
 document.body.style.color = 'white';
 } else {
 document.body.style.backgroundColor = 'white';
 document.body.style.color = 'black';
 }
 }, {
 threshold: 1.0
 });
 observer.observe(div);
 return () => {
 observer.disconnect();
 }
 }, []);

 return (
 <div ref={ref} style={{
 margin: 20,
 height: 100,
 width: 100,
 border: '2px solid black',
 backgroundColor: 'blue'
 }} />
 );
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Wrapping Effects in custom Hooks {/*wrapping-effects-in-custom-hooks*/}

Effects are an ["escape hatch":](/learn/escape-hatches) you use them when you need to "step outside React" and when there is no better built-in solution for your use case. If you find yourself often needing to manually write Effects, it's usually a sign that you need to extract some [custom Hooks](/learn/reusing-logic-with-custom-hooks) for common behaviors your components rely on.

For example, this `useChatRoom` custom Hook "hides" the logic of your Effect behind a more declarative API:

```js {1,11}
function useChatRoom({ serverUrl, roomId }) {
 useEffect(() => {
 const options = {
 serverUrl: serverUrl,
 roomId: roomId
 };
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId, serverUrl]);
}
```

Then you can use it from any component like this:

```js {4-7}
function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useChatRoom({
 roomId: roomId,
 serverUrl: serverUrl
 });
 // ...
```

There are also many excellent custom Hooks for every purpose available in the React ecosystem.

[Learn more about wrapping Effects in custom Hooks.](/learn/reusing-logic-with-custom-hooks)

<Recipes titleText="Examples of wrapping Effects in custom Hooks" titleId="examples-custom-hooks">

#### Custom `useChatRoom` Hook {/*custom-usechatroom-hook*/}

This example is identical to one of the [earlier examples,](#examples-connecting) but the logic is extracted to a custom Hook.

<Sandpack>

```js
import { useState } from 'react';
import { useChatRoom } from './useChatRoom.js';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useChatRoom({
 roomId: roomId,
 serverUrl: serverUrl
 });

 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={serverUrl}
 onChange={e => setServerUrl(e.target.value)}
 />
 </label>
 <h1>Welcome to the {roomId} room!</h1>
 </>
 );
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [show, setShow] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId} />}
 </>
 );
}
```

```js src/useChatRoom.js
import { useEffect } from 'react';
import { createConnection } from './chat.js';

export function useChatRoom({ serverUrl, roomId }) {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [roomId, serverUrl]);
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

<Solution />

#### Custom `useWindowListener` Hook {/*custom-usewindowlistener-hook*/}

This example is identical to one of the [earlier examples,](#examples-connecting) but the logic is extracted to a custom Hook.

<Sandpack>

```js
import { useState } from 'react';
import { useWindowListener } from './useWindowListener.js';

export default function App() {
 const [position, setPosition] = useState({ x: 0, y: 0 });

 useWindowListener('pointermove', (e) => {
 setPosition({ x: e.clientX, y: e.clientY });
 });

 return (
 <div style={{
 position: 'absolute',
 backgroundColor: 'pink',
 borderRadius: '50%',
 opacity: 0.6,
 transform: `translate(${position.x}px, ${position.y}px)`,
 pointerEvents: 'none',
 left: -20,
 top: -20,
 width: 40,
 height: 40,
 }} />
 );
}
```

```js src/useWindowListener.js
import { useState, useEffect } from 'react';

export function useWindowListener(eventType, listener) {
 useEffect(() => {
 window.addEventListener(eventType, listener);
 return () => {
 window.removeEventListener(eventType, listener);
 };
 }, [eventType, listener]);
}
```

```css
body {
 min-height: 300px;
}
```

</Sandpack>

<Solution />

#### Custom `useIntersectionObserver` Hook {/*custom-useintersectionobserver-hook*/}

This example is identical to one of the [earlier examples,](#examples-connecting) but the logic is partially extracted to a custom Hook.

<Sandpack>

```js
import Box from './Box.js';

export default function App() {
 return (
 <>
 <LongSection />
 <Box />
 <LongSection />
 <Box />
 <LongSection />
 </>
 );
}

function LongSection() {
 const items = [];
 for (let i = 0; i < 50; i++) {
 items.push(<li key={i}>Item #{i} (keep scrolling)</li>);
 }
 return <ul>{items}</ul>
}
```

```js src/Box.js active
import { useRef, useEffect } from 'react';
import { useIntersectionObserver } from './useIntersectionObserver.js';

export default function Box() {
 const ref = useRef(null);
 const isIntersecting = useIntersectionObserver(ref);

 useEffect(() => {
 if (isIntersecting) {
 document.body.style.backgroundColor = 'black';
 document.body.style.color = 'white';
 } else {
 document.body.style.backgroundColor = 'white';
 document.body.style.color = 'black';
 }
 }, [isIntersecting]);

 return (
 <div ref={ref} style={{
 margin: 20,
 height: 100,
 width: 100,
 border: '2px solid black',
 backgroundColor: 'blue'
 }} />
 );
}
```

```js src/useIntersectionObserver.js
import { useState, useEffect } from 'react';

export function useIntersectionObserver(ref) {
 const [isIntersecting, setIsIntersecting] = useState(false);

 useEffect(() => {
 const div = ref.current;
 const observer = new IntersectionObserver(entries => {
 const entry = entries[0];
 setIsIntersecting(entry.isIntersecting);
 }, {
 threshold: 1.0
 });
 observer.observe(div);
 return () => {
 observer.disconnect();
 }
 }, [ref]);

 return isIntersecting;
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Controlling a non-React widget {/*controlling-a-non-react-widget*/}

Sometimes, you want to keep an external system synchronized to some prop or state of your component.

For example, if you have a third-party map widget or a video player component written without React, you can use an Effect to call methods on it that make its state match the current state of your React component. This Effect creates an instance of a `MapWidget` class defined in `map-widget.js`. When you change the `zoomLevel` prop of the `Map` component, the Effect calls the `setZoom()` on the class instance to keep it synchronized:

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "leaflet": "1.9.1",
 "react": "latest",
 "react-dom": "latest",
 "react-scripts": "latest",
 "remarkable": "2.0.1"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js src/App.js
import { useState } from 'react';
import Map from './Map.js';

export default function App() {
 const [zoomLevel, setZoomLevel] = useState(0);
 return (
 <>
 Zoom level: {zoomLevel}x
 <button onClick={() => setZoomLevel(zoomLevel + 1)}>+</button>
 <button onClick={() => setZoomLevel(zoomLevel - 1)}>-</button>
 <hr />
 <Map zoomLevel={zoomLevel} />
 </>
 );
}
```

```js src/Map.js active
import { useRef, useEffect } from 'react';
import { MapWidget } from './map-widget.js';

export default function Map({ zoomLevel }) {
 const containerRef = useRef(null);
 const mapRef = useRef(null);

 useEffect(() => {
 if (mapRef.current === null) {
 mapRef.current = new MapWidget(containerRef.current);
 }

 const map = mapRef.current;
 map.setZoom(zoomLevel);
 }, [zoomLevel]);

 return (
 <div
 style={{ width: 200, height: 200 }}
 ref={containerRef}
 />
 );
}
```

```js src/map-widget.js
import 'leaflet/dist/leaflet.css';
import * as L from 'leaflet';

export class MapWidget {
 constructor(domNode) {
 this.map = L.map(domNode, {
 zoomControl: false,
 doubleClickZoom: false,
 boxZoom: false,
 keyboard: false,
 scrollWheelZoom: false,
 zoomAnimation: false,
 touchZoom: false,
 zoomSnap: 0.1
 });
 L.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
 maxZoom: 19,
 attribution: '© OpenStreetMap'
 }).addTo(this.map);
 this.map.setView([0, 0], 0);
 }
 setZoom(level) {
 this.map.setZoom(level);
 }
}
```

```css
button { margin: 5px; }
```

</Sandpack>

In this example, a cleanup function is not needed because the `MapWidget` class manages only the DOM node that was passed to it. After the `Map` React component is removed from the tree, both the DOM node and the `MapWidget` class instance will be automatically garbage-collected by the browser JavaScript engine.

---

### Fetching data with Effects {/*fetching-data-with-effects*/}

You can use an Effect to fetch data for your component. Note that [if you use a framework,](/learn/creating-a-react-app#full-stack-frameworks) using your framework's data fetching mechanism will be a lot more efficient than writing Effects manually.

If you want to fetch data from an Effect manually, your code might look like this:

```js
import { useState, useEffect } from 'react';
import { fetchBio } from './api.js';

export default function Page() {
 const [person, setPerson] = useState('Alice');
 const [bio, setBio] = useState(null);

 useEffect(() => {
 let ignore = false;
 setBio(null);
 fetchBio(person).then(result => {
 if (!ignore) {
 setBio(result);
 }
 });
 return () => {
 ignore = true;
 };
 }, [person]);

 // ...
```

Note the `ignore` variable which is initialized to `false`, and is set to `true` during cleanup. This ensures [your code doesn't suffer from "race conditions":](https://maxrozen.com/race-conditions-fetching-data-react-with-useeffect) network responses may arrive in a different order than you sent them.

<Sandpack>

{/* TODO(@poteto) - investigate potential false positives in react compiler validation */}
```js {expectedErrors: {'react-compiler': [9]}} src/App.js
import { useState, useEffect } from 'react';
import { fetchBio } from './api.js';

export default function Page() {
 const [person, setPerson] = useState('Alice');
 const [bio, setBio] = useState(null);
 useEffect(() => {
 let ignore = false;
 setBio(null);
 fetchBio(person).then(result => {
 if (!ignore) {
 setBio(result);
 }
 });
 return () => {
 ignore = true;
 }
 }, [person]);

 return (
 <>
 <select value={person} onChange={e => {
 setPerson(e.target.value);
 }}>
 <option value="Alice">Alice</option>
 <option value="Bob">Bob</option>
 <option value="Taylor">Taylor</option>
 </select>
 <hr />
 <p><i>{bio ?? 'Loading...'}</i></p>
 </>
 );
}
```

```js src/api.js hidden
export async function fetchBio(person) {
 const delay = person === 'Bob' ? 2000 : 200;
 return new Promise(resolve => {
 setTimeout(() => {
 resolve('This is ' + person + '’s bio.');
 }, delay);
 })
}
```

</Sandpack>

You can also rewrite using the [`async` / `await`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function) syntax, but you still need to provide a cleanup function:

<Sandpack>

```js src/App.js
import { useState, useEffect } from 'react';
import { fetchBio } from './api.js';

export default function Page() {
 const [person, setPerson] = useState('Alice');
 const [bio, setBio] = useState(null);
 useEffect(() => {
 async function startFetching() {
 setBio(null);
 const result = await fetchBio(person);
 if (!ignore) {
 setBio(result);
 }
 }

 let ignore = false;
 startFetching();
 return () => {
 ignore = true;
 }
 }, [person]);

 return (
 <>
 <select value={person} onChange={e => {
 setPerson(e.target.value);
 }}>
 <option value="Alice">Alice</option>
 <option value="Bob">Bob</option>
 <option value="Taylor">Taylor</option>
 </select>
 <hr />
 <p><i>{bio ?? 'Loading...'}</i></p>
 </>
 );
}
```

```js src/api.js hidden
export async function fetchBio(person) {
 const delay = person === 'Bob' ? 2000 : 200;
 return new Promise(resolve => {
 setTimeout(() => {
 resolve('This is ' + person + '’s bio.');
 }, delay);
 })
}
```

</Sandpack>

Writing data fetching directly in Effects gets repetitive and makes it difficult to add optimizations like caching and server rendering later. [It's easier to use a custom Hook--either your own or maintained by the community.](/learn/reusing-logic-with-custom-hooks#when-to-use-custom-hooks)

<DeepDive>

#### What are good alternatives to data fetching in Effects? {/*what-are-good-alternatives-to-data-fetching-in-effects*/}

Writing `fetch` calls inside Effects is a [popular way to fetch data](https://www.robinwieruch.de/react-hooks-fetch-data/), especially in fully client-side apps. This is, however, a very manual approach and it has significant downsides:

- **Effects don't run on the server.** This means that the initial server-rendered HTML will only include a loading state with no data. The client computer will have to download all JavaScript and render your app only to discover that now it needs to load the data. This is not very efficient.
- **Fetching directly in Effects makes it easy to create "network waterfalls".** You render the parent component, it fetches some data, renders the child components, and then they start fetching their data. If the network is not very fast, this is significantly slower than fetching all data in parallel.
- **Fetching directly in Effects usually means you don't preload or cache data.** For example, if the component unmounts and then mounts again, it would have to fetch the data again.
- **It's not very ergonomic.** There's quite a bit of boilerplate code involved when writing `fetch` calls in a way that doesn't suffer from bugs like [race conditions.](https://maxrozen.com/race-conditions-fetching-data-react-with-useeffect)

This list of downsides is not specific to React. It applies to fetching data on mount with any library. Like with routing, data fetching is not trivial to do well, so we recommend the following approaches:

- **If you use a [framework](/learn/creating-a-react-app#full-stack-frameworks), use its built-in data fetching mechanism.** Modern React frameworks have integrated data fetching mechanisms that are efficient and don't suffer from the above pitfalls.
- **Otherwise, consider using or building a client-side cache.** Popular open source solutions include [TanStack Query](https://tanstack.com/query/latest/), [useSWR](https://swr.vercel.app/), and [React Router 6.4+.](https://beta.reactrouter.com/en/main/start/overview) You can build your own solution too, in which case you would use Effects under the hood but also add logic for deduplicating requests, caching responses, and avoiding network waterfalls (by preloading data or hoisting data requirements to routes).

You can continue fetching data directly in Effects if neither of these approaches suit you.

</DeepDive>

---

### Specifying reactive dependencies {/*specifying-reactive-dependencies*/}

**Notice that you can't "choose" the dependencies of your Effect.** Every <CodeStep step={2}>reactive value</CodeStep> used by your Effect's code must be declared as a dependency. Your Effect's dependency list is determined by the surrounding code:

```js [[2, 1, "roomId"], [2, 2, "serverUrl"], [2, 5, "serverUrl"], [2, 5, "roomId"], [2, 8, "serverUrl"], [2, 8, "roomId"]]
function ChatRoom({ roomId }) { // This is a reactive value
 const [serverUrl, setServerUrl] = useState('https://localhost:1234'); // This is a reactive value too

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId); // This Effect reads these reactive values
 connection.connect();
 return () => connection.disconnect();
 }, [serverUrl, roomId]); // ✅ So you must specify them as dependencies of your Effect
 // ...
}
```

If either `serverUrl` or `roomId` change, your Effect will reconnect to the chat using the new values.

**[Reactive values](/learn/lifecycle-of-reactive-effects#effects-react-to-reactive-values) include props and all variables and functions declared directly inside of your component.** Since `roomId` and `serverUrl` are reactive values, you can't remove them from the dependencies. If you try to omit them and [your linter is correctly configured for React,](/learn/editor-setup#linting) the linter will flag this as a mistake you need to fix:

```js {8}
function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => connection.disconnect();
 }, []); // 🔴 React Hook useEffect has missing dependencies: 'roomId' and 'serverUrl'
 // ...
}
```

**To remove a dependency, you need to ["prove" to the linter that it *doesn't need* to be a dependency.](/learn/removing-effect-dependencies#removing-unnecessary-dependencies)** For example, you can move `serverUrl` out of your component to prove that it's not reactive and won't change on re-renders:

```js {1,8}
const serverUrl = 'https://localhost:1234'; // Not a reactive value anymore

function ChatRoom({ roomId }) {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]); // ✅ All dependencies declared
 // ...
}
```

Now that `serverUrl` is not a reactive value (and can't change on a re-render), it doesn't need to be a dependency. **If your Effect's code doesn't use any reactive values, its dependency list should be empty (`[]`):**

```js {1,2,9}
const serverUrl = 'https://localhost:1234'; // Not a reactive value anymore
const roomId = 'music'; // Not a reactive value anymore

function ChatRoom() {
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => connection.disconnect();
 }, []); // ✅ All dependencies declared
 // ...
}
```

[An Effect with empty dependencies](/learn/lifecycle-of-reactive-effects#what-an-effect-with-empty-dependencies-means) doesn't re-run when any of your component's props or state change.

<Pitfall>

If you have an existing codebase, you might have some Effects that suppress the linter like this:

```js {3-4}
useEffect(() => {
 // ...
 // 🔴 Avoid suppressing the linter like this:
 // eslint-ignore-next-line react-hooks/exhaustive-deps
}, []);
```

**When dependencies don't match the code, there is a high risk of introducing bugs.** By suppressing the linter, you "lie" to React about the values your Effect depends on. [Instead, prove they're unnecessary.](/learn/removing-effect-dependencies#removing-unnecessary-dependencies)

</Pitfall>

<Recipes titleText="Examples of passing reactive dependencies" titleId="examples-dependencies">

#### Passing a dependency array {/*passing-a-dependency-array*/}

If you specify the dependencies, your Effect runs **after the initial commit _and_ after commits with changed dependencies.**

```js {3}
useEffect(() => {
 // ...
}, [a, b]); // Runs again if a or b are different
```

In the below example, `serverUrl` and `roomId` are [reactive values,](/learn/lifecycle-of-reactive-effects#effects-react-to-reactive-values) so they both must be specified as dependencies. As a result, selecting a different room in the dropdown or editing the server URL input causes the chat to re-connect. However, since `message` isn't used in the Effect (and so it isn't a dependency), editing the message doesn't re-connect to the chat.

<Sandpack>

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');
 const [message, setMessage] = useState('');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [serverUrl, roomId]);

 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={serverUrl}
 onChange={e => setServerUrl(e.target.value)}
 />
 </label>
 <h1>Welcome to the {roomId} room!</h1>
 <label>
 Your message:{' '}
 <input value={message} onChange={e => setMessage(e.target.value)} />
 </label>
 </>
 );
}

export default function App() {
 const [show, setShow] = useState(false);
 const [roomId, setRoomId] = useState('general');
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 </label>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId}/>}
 </>
 );
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { margin-bottom: 10px; }
button { margin-left: 5px; }
```

</Sandpack>

<Solution />

#### Passing an empty dependency array {/*passing-an-empty-dependency-array*/}

If your Effect truly doesn't use any reactive values, it will only run **after the initial commit.**

```js {3}
useEffect(() => {
 // ...
}, []); // Does not run again (except once in development)
```

**Even with empty dependencies, setup and cleanup will [run one extra time in development](/learn/synchronizing-with-effects#how-to-handle-the-effect-firing-twice-in-development) to help you find bugs.**

In this example, both `serverUrl` and `roomId` are hardcoded. Since they're declared outside the component, they are not reactive values, and so they aren't dependencies. The dependency list is empty, so the Effect doesn't re-run on re-renders.

<Sandpack>

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';
const roomId = 'music';

function ChatRoom() {
 const [message, setMessage] = useState('');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => connection.disconnect();
 }, []);

 return (
 <>
 <h1>Welcome to the {roomId} room!</h1>
 <label>
 Your message:{' '}
 <input value={message} onChange={e => setMessage(e.target.value)} />
 </label>
 </>
 );
}

export default function App() {
 const [show, setShow] = useState(false);
 return (
 <>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 {show && <hr />}
 {show && <ChatRoom />}
 </>
 );
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

</Sandpack>

<Solution />

#### Passing no dependency array at all {/*passing-no-dependency-array-at-all*/}

If you pass no dependency array at all, your Effect runs **after every single commit** of your component.

```js {3}
useEffect(() => {
 // ...
}); // Always runs again
```

In this example, the Effect re-runs when you change `serverUrl` and `roomId`, which is sensible. However, it *also* re-runs when you change the `message`, which is probably undesirable. This is why usually you'll specify the dependency array.

<Sandpack>

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

function ChatRoom({ roomId }) {
 const [serverUrl, setServerUrl] = useState('https://localhost:1234');
 const [message, setMessage] = useState('');

 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }); // No dependency array at all

 return (
 <>
 <label>
 Server URL:{' '}
 <input
 value={serverUrl}
 onChange={e => setServerUrl(e.target.value)}
 />
 </label>
 <h1>Welcome to the {roomId} room!</h1>
 <label>
 Your message:{' '}
 <input value={message} onChange={e => setMessage(e.target.value)} />
 </label>
 </>
 );
}

export default function App() {
 const [show, setShow] = useState(false);
 const [roomId, setRoomId] = useState('general');
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 <button onClick={() => setShow(!show)}>
 {show ? 'Close chat' : 'Open chat'}
 </button>
 </label>
 {show && <hr />}
 {show && <ChatRoom roomId={roomId}/>}
 </>
 );
}
```

```js src/chat.js
export function createConnection(serverUrl, roomId) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { margin-bottom: 10px; }
button { margin-left: 5px; }
```

</Sandpack>

<Solution />

</Recipes>

---

### Updating state based on previous state from an Effect {/*updating-state-based-on-previous-state-from-an-effect*/}

When you want to update state based on previous state from an Effect, you might run into a problem:

```js {6,9}
function Counter() {
 const [count, setCount] = useState(0);

 useEffect(() => {
 const intervalId = setInterval(() => {
 setCount(count + 1); // You want to increment the counter every second...
 }, 1000)
 return () => clearInterval(intervalId);
 }, [count]); // 🚩 ... but specifying `count` as a dependency always resets the interval.
 // ...
}
```

Since `count` is a reactive value, it must be specified in the list of dependencies. However, that causes the Effect to cleanup and setup again every time the `count` changes. This is not ideal.

To fix this, [pass the `c => c + 1` state updater](/reference/react/useState#updating-state-based-on-the-previous-state) to `setCount`:

<Sandpack>

```js
import { useState, useEffect } from 'react';

export default function Counter() {
 const [count, setCount] = useState(0);

 useEffect(() => {
 const intervalId = setInterval(() => {
 setCount(c => c + 1); // ✅ Pass a state updater
 }, 1000);
 return () => clearInterval(intervalId);
 }, []); // ✅ Now count is not a dependency

 return <h1>{count}</h1>;
}
```

```css
label {
 display: block;
 margin-top: 20px;
 margin-bottom: 20px;
}

body {
 min-height: 150px;
}
```

</Sandpack>

Now that you're passing `c => c + 1` instead of `count + 1`, [your Effect no longer needs to depend on `count`.](/learn/removing-effect-dependencies#are-you-reading-some-state-to-calculate-the-next-state) As a result of this fix, it won't need to cleanup and setup the interval again every time the `count` changes.

---

### Removing unnecessary object dependencies {/*removing-unnecessary-object-dependencies*/}

If your Effect depends on an object or a function created during rendering, it might run too often. For example, this Effect re-connects after every commit because the `options` object is [different for every render:](/learn/removing-effect-dependencies#does-some-reactive-value-change-unintentionally)

```js {6-9,12,15}
const serverUrl = 'https://localhost:1234';

function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 const options = { // 🚩 This object is created from scratch on every re-render
 serverUrl: serverUrl,
 roomId: roomId
 };

 useEffect(() => {
 const connection = createConnection(options); // It's used inside the Effect
 connection.connect();
 return () => connection.disconnect();
 }, [options]); // 🚩 As a result, these dependencies are always different on a commit
 // ...
```

Avoid using an object created during rendering as a dependency. Instead, create the object inside the Effect:

<Sandpack>

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';

function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 useEffect(() => {
 const options = {
 serverUrl: serverUrl,
 roomId: roomId
 };
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]);

 return (
 <>
 <h1>Welcome to the {roomId} room!</h1>
 <input value={message} onChange={e => setMessage(e.target.value)} />
 </>
 );
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <hr />
 <ChatRoom roomId={roomId} />
 </>
 );
}
```

```js src/chat.js
export function createConnection({ serverUrl, roomId }) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

Now that you create the `options` object inside the Effect, the Effect itself only depends on the `roomId` string.

With this fix, typing into the input doesn't reconnect the chat. Unlike an object which gets re-created, a string like `roomId` doesn't change unless you set it to another value. [Read more about removing dependencies.](/learn/removing-effect-dependencies)

---

### Removing unnecessary function dependencies {/*removing-unnecessary-function-dependencies*/}

If your Effect depends on an object or a function created during rendering, it might run too often. For example, this Effect re-connects after every commit because the `createOptions` function is [different for every render:](/learn/removing-effect-dependencies#does-some-reactive-value-change-unintentionally)

```js {4-9,12,16}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 function createOptions() { // 🚩 This function is created from scratch on every re-render
 return {
 serverUrl: serverUrl,
 roomId: roomId
 };
 }

 useEffect(() => {
 const options = createOptions(); // It's used inside the Effect
 const connection = createConnection();
 connection.connect();
 return () => connection.disconnect();
 }, [createOptions]); // 🚩 As a result, these dependencies are always different on a commit
 // ...
```

By itself, creating a function from scratch on every re-render is not a problem. You don't need to optimize that. However, if you use it as a dependency of your Effect, it will cause your Effect to re-run after every commit.

Avoid using a function created during rendering as a dependency. Instead, declare it inside the Effect:

<Sandpack>

```js
import { useState, useEffect } from 'react';
import { createConnection } from './chat.js';

const serverUrl = 'https://localhost:1234';

function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 useEffect(() => {
 function createOptions() {
 return {
 serverUrl: serverUrl,
 roomId: roomId
 };
 }

 const options = createOptions();
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]);

 return (
 <>
 <h1>Welcome to the {roomId} room!</h1>
 <input value={message} onChange={e => setMessage(e.target.value)} />
 </>
 );
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <hr />
 <ChatRoom roomId={roomId} />
 </>
 );
}
```

```js src/chat.js
export function createConnection({ serverUrl, roomId }) {
 // A real implementation would actually connect to the server
 return {
 connect() {
 console.log('✅ Connecting to "' + roomId + '" room at ' + serverUrl + '...');
 },
 disconnect() {
 console.log('❌ Disconnected from "' + roomId + '" room at ' + serverUrl);
 }
 };
}
```

```css
input { display: block; margin-bottom: 20px; }
button { margin-left: 10px; }
```

</Sandpack>

Now that you define the `createOptions` function inside the Effect, the Effect itself only depends on the `roomId` string. With this fix, typing into the input doesn't reconnect the chat. Unlike a function which gets re-created, a string like `roomId` doesn't change unless you set it to another value. [Read more about removing dependencies.](/learn/removing-effect-dependencies)

---

### Reading the latest props and state from an Effect {/*reading-the-latest-props-and-state-from-an-effect*/}

By default, when you read a reactive value from an Effect, you have to add it as a dependency. This ensures that your Effect "reacts" to every change of that value. For most dependencies, that's the behavior you want.

**However, sometimes you'll want to read the *latest* props and state from an Effect without "reacting" to them.** For example, imagine you want to log the number of the items in the shopping cart for every page visit:

```js {3}
function Page({ url, shoppingCart }) {
 useEffect(() => {
 logVisit(url, shoppingCart.length);
 }, [url, shoppingCart]); // ✅ All dependencies declared
 // ...
}
```

**What if you want to log a new page visit after every `url` change, but *not* if only the `shoppingCart` changes?** You can't exclude `shoppingCart` from dependencies without breaking the [reactivity rules.](#specifying-reactive-dependencies) However, you can express that you *don't want* a piece of code to "react" to changes even though it is called from inside an Effect. [Declare an *Effect Event*](/learn/separating-events-from-effects#declaring-an-effect-event) with the [`useEffectEvent`](/reference/react/useEffectEvent) Hook, and move the code reading `shoppingCart` inside of it:

```js {2-4,7,8}
function Page({ url, shoppingCart }) {
 const onVisit = useEffectEvent(visitedUrl => {
 logVisit(visitedUrl, shoppingCart.length)
 });

 useEffect(() => {
 onVisit(url);
 }, [url]); // ✅ All dependencies declared
 // ...
}
```

**Effect Events are not reactive and must always be omitted from dependencies of your Effect.** This is what lets you put non-reactive code (where you can read the latest value of some props and state) inside of them. By reading `shoppingCart` inside of `onVisit`, you ensure that `shoppingCart` won't re-run your Effect.

[Read more about how Effect Events let you separate reactive and non-reactive code.](/learn/separating-events-from-effects#reading-latest-props-and-state-with-effect-events)

---

### Displaying different content on the server and the client {/*displaying-different-content-on-the-server-and-the-client*/}

If your app uses server rendering (either [directly](/reference/react-dom/server) or via a [framework](/learn/creating-a-react-app#full-stack-frameworks)), your component will render in two different environments. On the server, it will render to produce the initial HTML. On the client, React will run the rendering code again so that it can attach your event handlers to that HTML. This is why, for [hydration](/reference/react-dom/client/hydrateRoot#hydrating-server-rendered-html) to work, your initial render output must be identical on the client and the server.

In rare cases, you might need to display different content on the client. For example, if your app reads some data from [`localStorage`](https://developer.mozilla.org/en-US/docs/Web/API/Window/localStorage), it can't possibly do that on the server. Here is how you could implement this:

{/* TODO(@poteto) - investigate potential false positives in react compiler validation */}
```js {expectedErrors: {'react-compiler': [5]}}
function MyComponent() {
 const [didMount, setDidMount] = useState(false);

 useEffect(() => {
 setDidMount(true);
 }, []);

 if (didMount) {
 // ... return client-only JSX ...
 } else {
 // ... return initial JSX ...
 }
}
```

While the app is loading, the user will see the initial render output. Then, when it's loaded and hydrated, your Effect will run and set `didMount` to `true`, triggering a re-render. This will switch to the client-only render output. Effects don't run on the server, so this is why `didMount` was `false` during the initial server render.

Use this pattern sparingly. Keep in mind that users with a slow connection will see the initial content for quite a bit of time--potentially, many seconds--so you don't want to make jarring changes to your component's appearance. In many cases, you can avoid the need for this by conditionally showing different things with CSS.

---

## Troubleshooting {/*troubleshooting*/}

### My Effect runs twice when the component mounts {/*my-effect-runs-twice-when-the-component-mounts*/}

When Strict Mode is on, in development, React runs setup and cleanup one extra time before the actual setup.

This is a stress-test that verifies your Effect’s logic is implemented correctly. If this causes visible issues, your cleanup function is missing some logic. The cleanup function should stop or undo whatever the setup function was doing. The rule of thumb is that the user shouldn’t be able to distinguish between the setup being called once (as in production) and a setup → cleanup → setup sequence (as in development).

Read more about [how this helps find bugs](/learn/synchronizing-with-effects#step-3-add-cleanup-if-needed) and [how to fix your logic.](/learn/synchronizing-with-effects#how-to-handle-the-effect-firing-twice-in-development)

---

### My Effect runs after every re-render {/*my-effect-runs-after-every-re-render*/}

First, check that you haven't forgotten to specify the dependency array:

```js {3}
useEffect(() => {
 // ...
}); // 🚩 No dependency array: re-runs after every commit!
```

If you've specified the dependency array but your Effect still re-runs in a loop, it's because one of your dependencies is different on every re-render.

You can debug this problem by manually logging your dependencies to the console:

```js {5}
 useEffect(() => {
 // ..
 }, [serverUrl, roomId]);

 console.log([serverUrl, roomId]);
```

You can then right-click on the arrays from different re-renders in the console and select "Store as a global variable" for both of them. Assuming the first one got saved as `temp1` and the second one got saved as `temp2`, you can then use the browser console to check whether each dependency in both arrays is the same:

```js
Object.is(temp1[0], temp2[0]); // Is the first dependency the same between the arrays?
Object.is(temp1[1], temp2[1]); // Is the second dependency the same between the arrays?
Object.is(temp1[2], temp2[2]); // ... and so on for every dependency ...
```

When you find the dependency that is different on every re-render, you can usually fix it in one of these ways:

- [Updating state based on previous state from an Effect](#updating-state-based-on-previous-state-from-an-effect)
- [Removing unnecessary object dependencies](#removing-unnecessary-object-dependencies)
- [Removing unnecessary function dependencies](#removing-unnecessary-function-dependencies)
- [Reading the latest props and state from an Effect](#reading-the-latest-props-and-state-from-an-effect)

As a last resort (if these methods didn't help), wrap its creation with [`useMemo`](/reference/react/useMemo#memoizing-a-dependency-of-another-hook) or [`useCallback`](/reference/react/useCallback#preventing-an-effect-from-firing-too-often) (for functions).

---

### My Effect keeps re-running in an infinite cycle {/*my-effect-keeps-re-running-in-an-infinite-cycle*/}

If your Effect runs in an infinite cycle, these two things must be true:

- Your Effect is updating some state.
- That state leads to a re-render, which causes the Effect's dependencies to change.

Before you start fixing the problem, ask yourself whether your Effect is connecting to some external system (like DOM, network, a third-party widget, and so on). Why does your Effect need to set state? Does it synchronize with that external system? Or are you trying to manage your application's data flow with it?

If there is no external system, consider whether [removing the Effect altogether](/learn/you-might-not-need-an-effect) would simplify your logic.

If you're genuinely synchronizing with some external system, think about why and under what conditions your Effect should update the state. Has something changed that affects your component's visual output? If you need to keep track of some data that isn't used by rendering, a [ref](/reference/react/useRef#referencing-a-value-with-a-ref) (which doesn't trigger re-renders) might be more appropriate. Verify your Effect doesn't update the state (and trigger re-renders) more than needed.

Finally, if your Effect is updating the state at the right time, but there is still a loop, it's because that state update leads to one of the Effect's dependencies changing. [Read how to debug dependency changes.](/reference/react/useEffect#my-effect-runs-after-every-re-render)

---

### My cleanup logic runs even though my component didn't unmount {/*my-cleanup-logic-runs-even-though-my-component-didnt-unmount*/}

The cleanup function runs not only during unmount, but before every re-render with changed dependencies. Additionally, in development, React [runs setup+cleanup one extra time immediately after component mounts.](#my-effect-runs-twice-when-the-component-mounts)

If you have cleanup code without corresponding setup code, it's usually a code smell:

```js {2-5}
useEffect(() => {
 // 🔴 Avoid: Cleanup logic without corresponding setup logic
 return () => {
 doSomething();
 };
}, []);
```

Your cleanup logic should be "symmetrical" to the setup logic, and should stop or undo whatever setup did:

```js {2-3,5}
 useEffect(() => {
 const connection = createConnection(serverUrl, roomId);
 connection.connect();
 return () => {
 connection.disconnect();
 };
 }, [serverUrl, roomId]);
```

[Learn how the Effect lifecycle is different from the component's lifecycle.](/learn/lifecycle-of-reactive-effects#the-lifecycle-of-an-effect)

---

### My Effect does something visual, and I see a flicker before it runs {/*my-effect-does-something-visual-and-i-see-a-flicker-before-it-runs*/}

If your Effect must block the browser from [painting the screen,](/learn/render-and-commit#epilogue-browser-paint) replace `useEffect` with [`useLayoutEffect`](/reference/react/useLayoutEffect). Note that **this shouldn't be needed for the vast majority of Effects.** You'll only need this if it's crucial to run your Effect before the browser paint: for example, to measure and position a tooltip before the user sees it.

---
title: useEffectEvent
---

<Intro>

`useEffectEvent` is a React Hook that lets you separate events from Effects.

```js
const onEvent = useEffectEvent(callback)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useEffectEvent(callback)` {/*useeffectevent*/}

Call `useEffectEvent` at the top level of your component to create an Effect Event.

```js {4,6}
import { useEffectEvent, useEffect } from 'react';

function ChatRoom({ roomId, theme }) {
 const onConnected = useEffectEvent(() => {
 showNotification('Connected!', theme);
 });
}
```

Effect Events are a part of your Effect logic, but they behave more like an event handler. They always “see” the latest values from render (like props and state) without re-synchronizing your Effect, so they're excluded from Effect dependencies. See [Separating Events from Effects](/learn/separating-events-from-effects#extracting-non-reactive-logic-out-of-effects) to learn more.

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `callback`: A function containing the logic for your Effect Event. The function can accept any number of arguments and return any value. When you call the returned Effect Event function, the `callback` always accesses the latest committed values from render at the time of the call.

#### Returns {/*returns*/}

`useEffectEvent` returns an Effect Event function with the same type signature as your `callback`.

You can call this function inside `useEffect`, `useLayoutEffect`, `useInsertionEffect`, or from within other Effect Events in the same component.

#### Caveats {/*caveats*/}

* `useEffectEvent` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the Effect Event into it.
* Effect Events can only be called from inside Effects or other Effect Events. Do not call them during rendering or pass them to other components or Hooks. The [`eslint-plugin-react-hooks`](/reference/eslint-plugin-react-hooks) linter enforces this restriction.
* Do not use `useEffectEvent` to avoid specifying dependencies in your Effect's dependency array. This hides bugs and makes your code harder to understand. Only use it for logic that is genuinely an event fired from Effects.
* Effect Event functions do not have a stable identity. Their identity intentionally changes on every render.

<DeepDive>

#### Why are Effect Events not stable? {/*why-are-effect-events-not-stable*/}

Unlike `set` functions from `useState` or refs, Effect Event functions do not have a stable identity. Their identity intentionally changes on every render:

```js
// 🔴 Wrong: including Effect Event in dependencies
useEffect(() => {
 onSomething();
}, [onSomething]); // ESLint will warn about this
```

This is a deliberate design choice. Effect Events are meant to be called only from within Effects in the same component. Since you can only call them locally and cannot pass them to other components or include them in dependency arrays, a stable identity would serve no purpose, and would actually mask bugs.

The non-stable identity acts as a runtime assertion: if your code incorrectly depends on the function identity, you'll see the Effect re-running on every render, making the bug obvious.

This design reinforces that Effect Events conceptually belong to a particular effect, and are not a general purpose API to opt-out of reactivity.

</DeepDive>

---

## Usage {/*usage*/}

### Using an event in an Effect {/*using-an-event-in-an-effect*/}

Call `useEffectEvent` at the top level of your component to create an *Effect Event*:

```js [[1, 1, "onConnected"]]
const onConnected = useEffectEvent(() => {
 if (!muted) {
 showNotification('Connected!');
 }
});
```

`useEffectEvent` accepts an `event callback` and returns an <CodeStep step={1}>Effect Event</CodeStep>. The Effect Event is a function that can be called inside of Effects without re-connecting the Effect:

```js [[1, 3, "onConnected"]]
useEffect(() => {
 const connection = createConnection(roomId);
 connection.on('connected', onConnected);
 connection.connect();
 return () => {
 connection.disconnect();
 }
}, [roomId]);
```

Since `onConnected` is an <CodeStep step={1}>Effect Event</CodeStep>, `muted` and `onConnect` are not in the Effect dependencies.

<Pitfall>

##### Don't use Effect Events to skip dependencies {/*pitfall-skip-dependencies*/}

It might be tempting to use `useEffectEvent` to avoid listing dependencies that you think are "unnecessary." However, this hides bugs and makes your code harder to understand:

```js
// 🔴 Wrong: Using Effect Events to hide dependencies
const logVisit = useEffectEvent(() => {
 log(pageUrl);
});

useEffect(() => {
 logVisit()
}, []); // Missing pageUrl means you miss logs
```

If a value should cause your Effect to re-run, keep it as a dependency. Only use Effect Events for logic that genuinely should not re-trigger your Effect.

See [Separating Events from Effects](/learn/separating-events-from-effects) to learn more.

</Pitfall>

---

### Using a timer with latest values {/*using-a-timer-with-latest-values*/}

When you use `setInterval` or `setTimeout` in an Effect, you often want to read the latest values from render without restarting the timer whenever those values change.

This counter increments `count` by the current `increment` value every second. The `onTick` Effect Event reads the latest `count` and `increment` without causing the interval to restart:

<Sandpack>

```js
import { useState, useEffect, useEffectEvent } from 'react';

export default function Timer() {
 const [count, setCount] = useState(0);
 const [increment, setIncrement] = useState(1);

 const onTick = useEffectEvent(() => {
 setCount(count + increment);
 });

 useEffect(() => {
 const id = setInterval(() => {
 onTick();
 }, 1000);
 return () => {
 clearInterval(id);
 };
 }, []);

 return (
 <>
 <h1>
 Counter: {count}
 <button onClick={() => setCount(0)}>Reset</button>
 </h1>
 <hr />
 <p>
 Every second, increment by:
 <button disabled={increment === 0} onClick={() => {
 setIncrement(i => i - 1);
 }}>–</button>
 <b>{increment}</b>
 <button onClick={() => {
 setIncrement(i => i + 1);
 }}>+</button>
 </p>
 </>
 );
}
```

```css
button { margin: 10px; }
```

</Sandpack>

Try changing the increment value while the timer is running. The counter immediately uses the new increment value, but the timer keeps ticking smoothly without restarting.

---

### Using an event listener with latest values {/*using-an-event-listener-with-latest-values*/}

When you set up an event listener in an Effect, you often need to read the latest values from render in the callback. Without `useEffectEvent`, you would need to include the values in your dependencies, causing the listener to be removed and re-added on every change.

This example shows a dot that follows the cursor, but only when "Can move" is checked. The `onMove` Effect Event always reads the latest `canMove` value without re-running the Effect:

<Sandpack>

```js
import { useState, useEffect, useEffectEvent } from 'react';

export default function App() {
 const [position, setPosition] = useState({ x: 0, y: 0 });
 const [canMove, setCanMove] = useState(true);

 const onMove = useEffectEvent(e => {
 if (canMove) {
 setPosition({ x: e.clientX, y: e.clientY });
 }
 });

 useEffect(() => {
 window.addEventListener('pointermove', onMove);
 return () => window.removeEventListener('pointermove', onMove);
 }, []);

 return (
 <>
 <label>
 <input
 type="checkbox"
 checked={canMove}
 onChange={e => setCanMove(e.target.checked)}
 />
 The dot is allowed to move
 </label>
 <hr />
 <div style={{
 position: 'absolute',
 backgroundColor: 'pink',
 borderRadius: '50%',
 opacity: 0.6,
 transform: `translate(${position.x}px, ${position.y}px)`,
 pointerEvents: 'none',
 left: -20,
 top: -20,
 width: 40,
 height: 40,
 }} />
 </>
 );
}
```

```css
body {
 height: 200px;
}
```

</Sandpack>

Toggle the checkbox and move your cursor. The dot responds immediately to the checkbox state, but the event listener is only set up once when the component mounts.

---

### Avoid reconnecting to external systems {/*showing-a-notification-without-reconnecting*/}

A common use case for `useEffectEvent` is when you want to do something in response to an Effect, but that "something" depends on a value you don't want to react to.

In this example, a chat component connects to a room and shows a notification when connected. The user can mute notifications with a checkbox. However, you don't want to reconnect to the chat room every time the user changes the settings:

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "react": "latest",
 "react-dom": "latest",
 "react-scripts": "latest",
 "toastify-js": "1.12.0"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js
import { useState, useEffect, useEffectEvent } from 'react';
import { createConnection } from './chat.js';
import { showNotification } from './notifications.js';

function ChatRoom({ roomId, muted }) {
 const onConnected = useEffectEvent((roomId) => {
 console.log('✅ Connected to ' + roomId + ' (muted: ' + muted + ')');
 if (!muted) {
 showNotification('Connected to ' + roomId);
 }
 });

 useEffect(() => {
 const connection = createConnection(roomId);
 console.log('⏳ Connecting to ' + roomId + '...');
 connection.on('connected', () => {
 onConnected(roomId);
 });
 connection.connect();
 return () => {
 console.log('❌ Disconnected from ' + roomId);
 connection.disconnect();
 }
 }, [roomId]);

 return <h1>Welcome to the {roomId} room!</h1>;
}

export default function App() {
 const [roomId, setRoomId] = useState('general');
 const [muted, setMuted] = useState(false);
 return (
 <>
 <label>
 Choose the chat room:{' '}
 <select
 value={roomId}
 onChange={e => setRoomId(e.target.value)}
 >
 <option value="general">general</option>
 <option value="travel">travel</option>
 <option value="music">music</option>
 </select>
 </label>
 <label>
 <input
 type="checkbox"
 checked={muted}
 onChange={e => setMuted(e.target.checked)}
 />
 Mute notifications
 </label>
 <hr />
 <ChatRoom
 roomId={roomId}
 muted={muted}
 />
 </>
 );
}
```

```js src/chat.js
const serverUrl = 'https://localhost:1234';

export function createConnection(roomId) {
 // A real implementation would actually connect to the server
 let connectedCallback;
 let timeout;
 return {
 connect() {
 timeout = setTimeout(() => {
 if (connectedCallback) {
 connectedCallback();
 }
 }, 100);
 },
 on(event, callback) {
 if (connectedCallback) {
 throw Error('Cannot add the handler twice.');
 }
 if (event !== 'connected') {
 throw Error('Only "connected" event is supported.');
 }
 connectedCallback = callback;
 },
 disconnect() {
 clearTimeout(timeout);
 }
 };
}
```

```js src/notifications.js
import Toastify from 'toastify-js';
import 'toastify-js/src/toastify.css';

export function showNotification(message, theme) {
 Toastify({
 text: message,
 duration: 2000,
 gravity: 'top',
 position: 'right',
 style: {
 background: theme === 'dark' ? 'black' : 'white',
 color: theme === 'dark' ? 'white' : 'black',
 },
 }).showToast();
}
```

```css
label { display: block; margin-top: 10px; }
```

</Sandpack>

Try switching rooms. The chat reconnects and shows a notification. Now mute the notifications. Since `muted` is read inside the Effect Event rather than the Effect, the chat stays connected.

---

### Using Effect Events in custom Hooks {/*using-effect-events-in-custom-hooks*/}

You can use `useEffectEvent` inside your own custom Hooks. This lets you create reusable Hooks that encapsulate Effects while keeping some values non-reactive:

<Sandpack>

```js
import { useState, useEffect, useEffectEvent } from 'react';

function useInterval(callback, delay) {
 const onTick = useEffectEvent(callback);

 useEffect(() => {
 if (delay === null) {
 return;
 }
 const id = setInterval(() => {
 onTick();
 }, delay);
 return () => clearInterval(id);
 }, [delay]);
}

function Counter({ incrementBy }) {
 const [count, setCount] = useState(0);

 useInterval(() => {
 setCount(c => c + incrementBy);
 }, 1000);

 return (
 <div>
 <h2>Count: {count}</h2>
 <p>Incrementing by {incrementBy} every second</p>
 </div>
 );
}

export default function App() {
 const [incrementBy, setIncrementBy] = useState(1);

 return (
 <>
 <label>
 Increment by:{' '}
 <select
 value={incrementBy}
 onChange={(e) => setIncrementBy(Number(e.target.value))}
 >
 <option value={1}>1</option>
 <option value={5}>5</option>
 <option value={10}>10</option>
 </select>
 </label>
 <hr />
 <Counter incrementBy={incrementBy} />
 </>
 );
}
```

```css
label { display: block; margin-bottom: 8px; }
```

</Sandpack>

In this example, `useInterval` is a custom Hook that sets up an interval. The `callback` passed to it is wrapped in an Effect Event, so the interval does not reset even if a new `callback` is passed in every render.

---

## Troubleshooting {/*troubleshooting*/}

### I'm getting an error: "A function wrapped in useEffectEvent can't be called during rendering" {/*cant-call-during-rendering*/}

This error means you're calling an Effect Event function during the render phase of your component. Effect Events can only be called from inside Effects or other Effect Events.

```js
function MyComponent({ data }) {
 const onLog = useEffectEvent(() => {
 console.log(data);
 });

 // 🔴 Wrong: calling during render
 onLog();

 // ✅ Correct: call from an Effect
 useEffect(() => {
 onLog();
 }, []);

 return <div>{data}</div>;
}
```

If you need to run logic during render, don't wrap it in `useEffectEvent`. Call the logic directly or move it into an Effect.

---

### I'm getting a lint error: "Functions returned from useEffectEvent must not be included in the dependency array" {/*effect-event-in-deps*/}

If you see a warning like "Functions returned from `useEffectEvent` must not be included in the dependency array", remove the Effect Event from your dependencies:

```js
const onSomething = useEffectEvent(() => {
 // ...
});

// 🔴 Wrong: Effect Event in dependencies
useEffect(() => {
 onSomething();
}, [onSomething]);

// ✅ Correct: no Effect Event in dependencies
useEffect(() => {
 onSomething();
}, []);
```

Effect Events are designed to be called from Effects without being listed as dependencies. The linter enforces this because the function identity is [intentionally not stable](#why-are-effect-events-not-stable). Including it would cause your Effect to re-run on every render.

---

### I'm getting a lint error: "... is a function created with useEffectEvent, and can only be called from Effects" {/*effect-event-called-outside-effect*/}

If you see a warning like "... is a function created with React Hook `useEffectEvent`, and can only be called from Effects and Effect Events", you're calling the function from the wrong place:

```js
const onSomething = useEffectEvent(() => {
 console.log(value);
});

// 🔴 Wrong: calling from event handler
function handleClick() {
 onSomething();
}

// 🔴 Wrong: passing to child component
return <Child onSomething={onSomething} />;

// ✅ Correct: calling from Effect
useEffect(() => {
 onSomething();
}, []);
```

Effect Events are specifically designed to be used in Effects local to the component they're defined in. If you need a callback for event handlers or to pass to children, use a regular function or `useCallback` instead.

---
title: useId
---

<Intro>

`useId` is a React Hook for generating unique IDs that can be passed to accessibility attributes.

```js
const id = useId()
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useId()` {/*useid*/}

Call `useId` at the top level of your component to generate a unique ID:

```js
import { useId } from 'react';

function PasswordField() {
 const passwordHintId = useId();
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

`useId` does not take any parameters.

#### Returns {/*returns*/}

`useId` returns a unique ID string associated with this particular `useId` call in this particular component.

#### Caveats {/*caveats*/}

* `useId` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.

* `useId` **should not be used to generate cache keys** for [use()](/reference/react/use). The ID is stable when a component is mounted but may change during rendering. Cache keys should be generated from your data.

* `useId` **should not be used to generate keys** in a list. [Keys should be generated from your data.](/learn/rendering-lists#where-to-get-your-key)

* `useId` currently cannot be used in [async Server Components](/reference/rsc/server-components#async-components-with-server-components).

---

## Usage {/*usage*/}

<Pitfall>

**Do not call `useId` to generate keys in a list.** [Keys should be generated from your data.](/learn/rendering-lists#where-to-get-your-key)

</Pitfall>

### Generating unique IDs for accessibility attributes {/*generating-unique-ids-for-accessibility-attributes*/}

Call `useId` at the top level of your component to generate a unique ID:

```js [[1, 4, "passwordHintId"]]
import { useId } from 'react';

function PasswordField() {
 const passwordHintId = useId();
 // ...
```

You can then pass the <CodeStep step={1}>generated ID</CodeStep> to different attributes:

```js [[1, 2, "passwordHintId"], [1, 3, "passwordHintId"]]
<>
 <input type="password" aria-describedby={passwordHintId} />
 <p id={passwordHintId}>
</>
```

**Let's walk through an example to see when this is useful.**

[HTML accessibility attributes](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA) like [`aria-describedby`](https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Attributes/aria-describedby) let you specify that two tags are related to each other. For example, you can specify that an element (like an input) is described by another element (like a paragraph).

In regular HTML, you would write it like this:

```html {5,8}
<label>
 Password:
 <input
 type="password"
 aria-describedby="password-hint"
 />
</label>
<p id="password-hint">
 The password should contain at least 18 characters
</p>
```

However, hardcoding IDs like this is not a good practice in React. A component may be rendered more than once on the page--but IDs have to be unique! Instead of hardcoding an ID, generate a unique ID with `useId`:

```js {4,11,14}
import { useId } from 'react';

function PasswordField() {
 const passwordHintId = useId();
 return (
 <>
 <label>
 Password:
 <input
 type="password"
 aria-describedby={passwordHintId}
 />
 </label>
 <p id={passwordHintId}>
 The password should contain at least 18 characters
 </p>
 </>
 );
}
```

Now, even if `PasswordField` appears multiple times on the screen, the generated IDs won't clash.

<Sandpack>

```js
import { useId } from 'react';

function PasswordField() {
 const passwordHintId = useId();
 return (
 <>
 <label>
 Password:
 <input
 type="password"
 aria-describedby={passwordHintId}
 />
 </label>
 <p id={passwordHintId}>
 The password should contain at least 18 characters
 </p>
 </>
 );
}

export default function App() {
 return (
 <>
 <h2>Choose password</h2>
 <PasswordField />
 <h2>Confirm password</h2>
 <PasswordField />
 </>
 );
}
```

```css
input { margin: 5px; }
```

</Sandpack>

[Watch this video](https://www.youtube.com/watch?v=0dNzNcuEuOo) to see the difference in the user experience with assistive technologies.

<Pitfall>

With [server rendering](/reference/react-dom/server), **`useId` requires an identical component tree on the server and the client**. If the trees you render on the server and the client don't match exactly, the generated IDs won't match.

</Pitfall>

<DeepDive>

#### Why is useId better than an incrementing counter? {/*why-is-useid-better-than-an-incrementing-counter*/}

You might be wondering why `useId` is better than incrementing a global variable like `nextId++`.

The primary benefit of `useId` is that React ensures that it works with [server rendering.](/reference/react-dom/server) During server rendering, your components generate HTML output. Later, on the client, [hydration](/reference/react-dom/client/hydrateRoot) attaches your event handlers to the generated HTML. For hydration to work, the client output must match the server HTML.

This is very difficult to guarantee with an incrementing counter because the order in which the Client Components are hydrated may not match the order in which the server HTML was emitted. By calling `useId`, you ensure that hydration will work, and the output will match between the server and the client.

Inside React, `useId` is generated from the "parent path" of the calling component. This is why, if the client and the server tree are the same, the "parent path" will match up regardless of rendering order.

</DeepDive>

---

### Generating IDs for several related elements {/*generating-ids-for-several-related-elements*/}

If you need to give IDs to multiple related elements, you can call `useId` to generate a shared prefix for them:

<Sandpack>

```js
import { useId } from 'react';

export default function Form() {
 const id = useId();
 return (
 <form>
 <label htmlFor={id + '-firstName'}>First Name:</label>
 <input id={id + '-firstName'} type="text" />
 <hr />
 <label htmlFor={id + '-lastName'}>Last Name:</label>
 <input id={id + '-lastName'} type="text" />
 </form>
 );
}
```

```css
input { margin: 5px; }
```

</Sandpack>

This lets you avoid calling `useId` for every single element that needs a unique ID.

---

### Specifying a shared prefix for all generated IDs {/*specifying-a-shared-prefix-for-all-generated-ids*/}

If you render multiple independent React applications on a single page, pass `identifierPrefix` as an option to your [`createRoot`](/reference/react-dom/client/createRoot#parameters) or [`hydrateRoot`](/reference/react-dom/client/hydrateRoot) calls. This ensures that the IDs generated by the two different apps never clash because every identifier generated with `useId` will start with the distinct prefix you've specified.

<Sandpack>

```html public/index.html
<!DOCTYPE html>
<html>
 <head><title>My app</title></head>
 <body>
 <div id="root1"></div>
 <div id="root2"></div>
 </body>
</html>
```

```js
import { useId } from 'react';

function PasswordField() {
 const passwordHintId = useId();
 console.log('Generated identifier:', passwordHintId)
 return (
 <>
 <label>
 Password:
 <input
 type="password"
 aria-describedby={passwordHintId}
 />
 </label>
 <p id={passwordHintId}>
 The password should contain at least 18 characters
 </p>
 </>
 );
}

export default function App() {
 return (
 <>
 <h2>Choose password</h2>
 <PasswordField />
 </>
 );
}
```

```js src/index.js active
import { createRoot } from 'react-dom/client';
import App from './App.js';
import './styles.css';

const root1 = createRoot(document.getElementById('root1'), {
 identifierPrefix: 'my-first-app-'
});
root1.render(<App />);

const root2 = createRoot(document.getElementById('root2'), {
 identifierPrefix: 'my-second-app-'
});
root2.render(<App />);
```

```css
#root1 {
 border: 5px solid blue;
 padding: 10px;
 margin: 5px;
}

#root2 {
 border: 5px solid green;
 padding: 10px;
 margin: 5px;
}

input { margin: 5px; }
```

</Sandpack>

---

### Using the same ID prefix on the client and the server {/*using-the-same-id-prefix-on-the-client-and-the-server*/}

If you [render multiple independent React apps on the same page](#specifying-a-shared-prefix-for-all-generated-ids), and some of these apps are server-rendered, make sure that the `identifierPrefix` you pass to the [`hydrateRoot`](/reference/react-dom/client/hydrateRoot) call on the client side is the same as the `identifierPrefix` you pass to the [server APIs](/reference/react-dom/server) such as [`renderToPipeableStream`.](/reference/react-dom/server/renderToPipeableStream)

```js
// Server
import { renderToPipeableStream } from 'react-dom/server';

const { pipe } = renderToPipeableStream(
 <App />,
 { identifierPrefix: 'react-app1' }
);
```

```js
// Client
import { hydrateRoot } from 'react-dom/client';

const domNode = document.getElementById('root');
const root = hydrateRoot(
 domNode,
 reactNode,
 { identifierPrefix: 'react-app1' }
);
```

You do not need to pass `identifierPrefix` if you only have one React app on the page.

---
title: useImperativeHandle
---

<Intro>

`useImperativeHandle` is a React Hook that lets you customize the handle exposed as a [ref.](/learn/manipulating-the-dom-with-refs)

```js
useImperativeHandle(ref, createHandle, dependencies?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useImperativeHandle(ref, createHandle, dependencies?)` {/*useimperativehandle*/}

Call `useImperativeHandle` at the top level of your component to customize the ref handle it exposes:

```js
import { useImperativeHandle } from 'react';

function MyInput({ ref }) {
 useImperativeHandle(ref, () => {
 return {
 // ... your methods ...
 };
 }, []);
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `ref`: The `ref` you received as a prop to the `MyInput` component.

* `createHandle`: A function that takes no arguments and returns the ref handle you want to expose. That ref handle can have any type. Usually, you will return an object with the methods you want to expose.

* **optional** `dependencies`: The list of all reactive values referenced inside of the `createHandle` code. Reactive values include props, state, and all the variables and functions declared directly inside your component body. If your linter is [configured for React](/learn/editor-setup#linting), it will verify that every reactive value is correctly specified as a dependency. The list of dependencies must have a constant number of items and be written inline like `[dep1, dep2, dep3]`. React will compare each dependency with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. If a re-render resulted in a change to some dependency, or if you omitted this argument, your `createHandle` function will re-execute, and the newly created handle will be assigned to the ref.

<Note>

Starting with React 19, [`ref` is available as a prop.](/blog/2024/12/05/react-19#ref-as-a-prop) In React 18 and earlier, it was necessary to get the `ref` from [`forwardRef`.](/reference/react/forwardRef)

</Note>

#### Returns {/*returns*/}

`useImperativeHandle` returns `undefined`.

---

## Usage {/*usage*/}

### Exposing a custom ref handle to the parent component {/*exposing-a-custom-ref-handle-to-the-parent-component*/}

To expose a DOM node to the parent element, pass in the `ref` prop to the node.

```js {2}
function MyInput({ ref }) {
 return <input ref={ref} />;
};
```

With the code above, [a ref to `MyInput` will receive the `<input>` DOM node.](/learn/manipulating-the-dom-with-refs) However, you can expose a custom value instead. To customize the exposed handle, call `useImperativeHandle` at the top level of your component:

```js {4-8}
import { useImperativeHandle } from 'react';

function MyInput({ ref }) {
 useImperativeHandle(ref, () => {
 return {
 // ... your methods ...
 };
 }, []);

 return <input />;
};
```

Note that in the code above, the `ref` is no longer passed to the `<input>`.

For example, suppose you don't want to expose the entire `<input>` DOM node, but you want to expose two of its methods: `focus` and `scrollIntoView`. To do this, keep the real browser DOM in a separate ref. Then use `useImperativeHandle` to expose a handle with only the methods that you want the parent component to call:

```js {7-14}
import { useRef, useImperativeHandle } from 'react';

function MyInput({ ref }) {
 const inputRef = useRef(null);

 useImperativeHandle(ref, () => {
 return {
 focus() {
 inputRef.current.focus();
 },
 scrollIntoView() {
 inputRef.current.scrollIntoView();
 },
 };
 }, []);

 return <input ref={inputRef} />;
};
```

Now, if the parent component gets a ref to `MyInput`, it will be able to call the `focus` and `scrollIntoView` methods on it. However, it will not have full access to the underlying `<input>` DOM node.

<Sandpack>

```js
import { useRef } from 'react';
import MyInput from './MyInput.js';

export default function Form() {
 const ref = useRef(null);

 function handleClick() {
 ref.current.focus();
 // This won't work because the DOM node isn't exposed:
 // ref.current.style.opacity = 0.5;
 }

 return (
 <form>
 <MyInput placeholder="Enter your name" ref={ref} />
 <button type="button" onClick={handleClick}>
 Edit
 </button>
 </form>
 );
}
```

```js src/MyInput.js
import { useRef, useImperativeHandle } from 'react';

function MyInput({ ref, ...props }) {
 const inputRef = useRef(null);

 useImperativeHandle(ref, () => {
 return {
 focus() {
 inputRef.current.focus();
 },
 scrollIntoView() {
 inputRef.current.scrollIntoView();
 },
 };
 }, []);

 return <input {...props} ref={inputRef} />;
};

export default MyInput;
```

```css
input {
 margin: 5px;
}
```

</Sandpack>

---

### Exposing your own imperative methods {/*exposing-your-own-imperative-methods*/}

The methods you expose via an imperative handle don't have to match the DOM methods exactly. For example, this `Post` component exposes a `scrollAndFocusAddComment` method via an imperative handle. This lets the parent `Page` scroll the list of comments *and* focus the input field when you click the button:

<Sandpack>

```js
import { useRef } from 'react';
import Post from './Post.js';

export default function Page() {
 const postRef = useRef(null);

 function handleClick() {
 postRef.current.scrollAndFocusAddComment();
 }

 return (
 <>
 <button onClick={handleClick}>
 Write a comment
 </button>
 <Post ref={postRef} />
 </>
 );
}
```

```js src/Post.js
import { useRef, useImperativeHandle } from 'react';
import CommentList from './CommentList.js';
import AddComment from './AddComment.js';

function Post({ ref }) {
 const commentsRef = useRef(null);
 const addCommentRef = useRef(null);

 useImperativeHandle(ref, () => {
 return {
 scrollAndFocusAddComment() {
 commentsRef.current.scrollToBottom();
 addCommentRef.current.focus();
 }
 };
 }, []);

 return (
 <>
 <article>
 <p>Welcome to my blog!</p>
 </article>
 <CommentList ref={commentsRef} />
 <AddComment ref={addCommentRef} />
 </>
 );
};

export default Post;
```

```js src/CommentList.js
import { useRef, useImperativeHandle } from 'react';

function CommentList({ ref }) {
 const divRef = useRef(null);

 useImperativeHandle(ref, () => {
 return {
 scrollToBottom() {
 const node = divRef.current;
 node.scrollTop = node.scrollHeight;
 }
 };
 }, []);

 let comments = [];
 for (let i = 0; i < 50; i++) {
 comments.push(<p key={i}>Comment #{i}</p>);
 }

 return (
 <div className="CommentList" ref={divRef}>
 {comments}
 </div>
 );
}

export default CommentList;
```

```js src/AddComment.js
import { useRef, useImperativeHandle } from 'react';

function AddComment({ ref }) {
 return <input placeholder="Add comment..." ref={ref} />;
}

export default AddComment;
```

```css
.CommentList {
 height: 100px;
 overflow: scroll;
 border: 1px solid black;
 margin-top: 20px;
 margin-bottom: 20px;
}
```

</Sandpack>

<Pitfall>

**Do not overuse refs.** You should only use refs for *imperative* behaviors that you can't express as props: for example, scrolling to a node, focusing a node, triggering an animation, selecting text, and so on.

**If you can express something as a prop, you should not use a ref.** For example, instead of exposing an imperative handle like `{ open, close }` from a `Modal` component, it is better to take `isOpen` as a prop like `<Modal isOpen={isOpen} />`. [Effects](/learn/synchronizing-with-effects) can help you expose imperative behaviors via props.

</Pitfall>

---
title: useInsertionEffect
---

<Pitfall>

`useInsertionEffect` is for CSS-in-JS library authors. Unless you are working on a CSS-in-JS library and need a place to inject the styles, you probably want [`useEffect`](/reference/react/useEffect) or [`useLayoutEffect`](/reference/react/useLayoutEffect) instead.

</Pitfall>

<Intro>

`useInsertionEffect` allows inserting elements into the DOM before any layout Effects fire.

```js
useInsertionEffect(setup, dependencies?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useInsertionEffect(setup, dependencies?)` {/*useinsertioneffect*/}

Call `useInsertionEffect` to insert styles before any Effects fire that may need to read layout:

```js
import { useInsertionEffect } from 'react';

// Inside your CSS-in-JS library
function useCSS(rule) {
 useInsertionEffect(() => {
 // ... inject <style> tags here ...
 });
 return rule;
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `setup`: The function with your Effect's logic. Your setup function may also optionally return a *cleanup* function. When your component is added to the DOM, but before any layout Effects fire, React will run your setup function. After every re-render with changed dependencies, React will first run the cleanup function (if you provided it) with the old values, and then run your setup function with the new values. When your component is removed from the DOM, React will run your cleanup function.

* **optional** `dependencies`: The list of all reactive values referenced inside of the `setup` code. Reactive values include props, state, and all the variables and functions declared directly inside your component body. If your linter is [configured for React](/learn/editor-setup#linting), it will verify that every reactive value is correctly specified as a dependency. The list of dependencies must have a constant number of items and be written inline like `[dep1, dep2, dep3]`. React will compare each dependency with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison algorithm. If you don't specify the dependencies at all, your Effect will re-run after every re-render of the component.

#### Returns {/*returns*/}

`useInsertionEffect` returns `undefined`.

#### Caveats {/*caveats*/}

* Effects only run on the client. They don't run during server rendering.
* You can't update state from inside `useInsertionEffect`.
* By the time `useInsertionEffect` runs, refs are not attached yet.
* `useInsertionEffect` may run either before or after the DOM has been updated. You shouldn't rely on the DOM being updated at any particular time.
* Unlike other types of Effects, which fire cleanup for every Effect and then setup for every Effect, `useInsertionEffect` will fire both cleanup and setup one component at a time. This results in an "interleaving" of the cleanup and setup functions.
---

## Usage {/*usage*/}

### Injecting dynamic styles from CSS-in-JS libraries {/*injecting-dynamic-styles-from-css-in-js-libraries*/}

Traditionally, you would style React components using plain CSS.

```js
// In your JS file:
<button className="success" />

// In your CSS file:
.success { color: green; }
```

Some teams prefer to author styles directly in JavaScript code instead of writing CSS files. This usually requires using a CSS-in-JS library or a tool. There are three common approaches to CSS-in-JS:

1. Static extraction to CSS files with a compiler
2. Inline styles, e.g. `<div style={{ opacity: 1 }}>`
3. Runtime injection of `<style>` tags

If you use CSS-in-JS, we recommend a combination of the first two approaches (CSS files for static styles, inline styles for dynamic styles). **We don't recommend runtime `<style>` tag injection for two reasons:**

1. Runtime injection forces the browser to recalculate the styles a lot more often.
2. Runtime injection can be very slow if it happens at the wrong time in the React lifecycle.

The first problem is not solvable, but `useInsertionEffect` helps you solve the second problem.

Call `useInsertionEffect` to insert the styles before any layout Effects fire:

```js {4-11}
// Inside your CSS-in-JS library
let isInserted = new Set();
function useCSS(rule) {
 useInsertionEffect(() => {
 // As explained earlier, we don't recommend runtime injection of <style> tags.
 // But if you have to do it, then it's important to do in useInsertionEffect.
 if (!isInserted.has(rule)) {
 isInserted.add(rule);
 document.head.appendChild(getStyleForRule(rule));
 }
 });
 return rule;
}

function Button() {
 const className = useCSS('...');
 return <div className={className} />;
}
```

Similarly to `useEffect`, `useInsertionEffect` does not run on the server. If you need to collect which CSS rules have been used on the server, you can do it during rendering:

```js {1,4-6}
let collectedRulesSet = new Set();

function useCSS(rule) {
 if (typeof window === 'undefined') {
 collectedRulesSet.add(rule);
 }
 useInsertionEffect(() => {
 // ...
 });
 return rule;
}
```

[Read more about upgrading CSS-in-JS libraries with runtime injection to `useInsertionEffect`.](https://github.com/reactwg/react-18/discussions/110)

<DeepDive>

#### How is this better than injecting styles during rendering or useLayoutEffect? {/*how-is-this-better-than-injecting-styles-during-rendering-or-uselayouteffect*/}

If you insert styles during rendering and React is processing a [non-blocking update,](/reference/react/useTransition#perform-non-blocking-updates-with-actions) the browser will recalculate the styles every single frame while rendering a component tree, which can be **extremely slow.**

`useInsertionEffect` is better than inserting styles during [`useLayoutEffect`](/reference/react/useLayoutEffect) or [`useEffect`](/reference/react/useEffect) because it ensures that by the time other Effects run in your components, the `<style>` tags have already been inserted. Otherwise, layout calculations in regular Effects would be wrong due to outdated styles.

</DeepDive>

---
title: useLayoutEffect
---

<Pitfall>

`useLayoutEffect` can hurt performance. Prefer [`useEffect`](/reference/react/useEffect) when possible.

</Pitfall>

<Intro>

`useLayoutEffect` is a version of [`useEffect`](/reference/react/useEffect) that fires before the browser repaints the screen.

```js
useLayoutEffect(setup, dependencies?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useLayoutEffect(setup, dependencies?)` {/*useinsertioneffect*/}

Call `useLayoutEffect` to perform the layout measurements before the browser repaints the screen:

```js
import { useState, useRef, useLayoutEffect } from 'react';

function Tooltip() {
 const ref = useRef(null);
 const [tooltipHeight, setTooltipHeight] = useState(0);

 useLayoutEffect(() => {
 const { height } = ref.current.getBoundingClientRect();
 setTooltipHeight(height);
 }, []);
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `setup`: The function with your Effect's logic. Your setup function may also optionally return a *cleanup* function. After your [component commits](/learn/render-and-commit#step-3-react-commits-changes-to-the-dom) to the DOM and before the browser repaints the screen, React will run your setup function. After every commit with changed dependencies, React will first run the cleanup function (if you provided it) with the old values, and then run your setup function with the new values. Before your component is removed from the DOM, React will run your cleanup function.

* **optional** `dependencies`: The list of all reactive values referenced inside of the `setup` code. Reactive values include props, state, and all the variables and functions declared directly inside your component body. If your linter is [configured for React](/learn/editor-setup#linting), it will verify that every reactive value is correctly specified as a dependency. The list of dependencies must have a constant number of items and be written inline like `[dep1, dep2, dep3]`. React will compare each dependency with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. If you omit this argument, your Effect will re-run after every commit of the component.

#### Returns {/*returns*/}

`useLayoutEffect` returns `undefined`.

#### Caveats {/*caveats*/}

* `useLayoutEffect` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a component and move the Effect there.

* When Strict Mode is on, React will **run one extra development-only setup+cleanup cycle** before the first real setup. This is a stress-test that ensures that your cleanup logic "mirrors" your setup logic and that it stops or undoes whatever the setup is doing. If this causes a problem, [implement the cleanup function.](/learn/synchronizing-with-effects#how-to-handle-the-effect-firing-twice-in-development)

* If some of your dependencies are objects or functions defined inside the component, there is a risk that they will **cause the Effect to re-run more often than needed.** To fix this, remove unnecessary [object](/reference/react/useEffect#removing-unnecessary-object-dependencies) and [function](/reference/react/useEffect#removing-unnecessary-function-dependencies) dependencies. You can also [extract state updates](/reference/react/useEffect#updating-state-based-on-previous-state-from-an-effect) and [non-reactive logic](/reference/react/useEffect#reading-the-latest-props-and-state-from-an-effect) outside of your Effect.

* Effects **only run on the client.** They don't run during server rendering.

* The code inside `useLayoutEffect` and all state updates scheduled from it **block the browser from repainting the screen.** When used excessively, this makes your app slow. When possible, prefer [`useEffect`.](/reference/react/useEffect)

* If you trigger a state update inside `useLayoutEffect`, React will execute all remaining Effects immediately including `useEffect`.

---

## Usage {/*usage*/}

### Measuring layout before the browser repaints the screen {/*measuring-layout-before-the-browser-repaints-the-screen*/}

Most components don't need to know their position and size on the screen to decide what to render. They only return some JSX. Then the browser calculates their *layout* (position and size) and repaints the screen.

Sometimes, that's not enough. Imagine a tooltip that appears next to some element on hover. If there's enough space, the tooltip should appear above the element, but if it doesn't fit, it should appear below. In order to render the tooltip at the right final position, you need to know its height (i.e. whether it fits at the top).

To do this, you need to render in two passes:

1. Render the tooltip anywhere (even with a wrong position).
2. Measure its height and decide where to place the tooltip.
3. Render the tooltip *again* in the correct place.

**All of this needs to happen before the browser repaints the screen.** You don't want the user to see the tooltip moving. Call `useLayoutEffect` to perform the layout measurements before the browser repaints the screen:

```js {5-8}
function Tooltip() {
 const ref = useRef(null);
 const [tooltipHeight, setTooltipHeight] = useState(0); // You don't know real height yet

 useLayoutEffect(() => {
 const { height } = ref.current.getBoundingClientRect();
 setTooltipHeight(height); // Re-render now that you know the real height
 }, []);

 // ...use tooltipHeight in the rendering logic below...
}
```

Here's how this works step by step:

1. `Tooltip` renders with the initial `tooltipHeight = 0` (so the tooltip may be wrongly positioned).
2. React places it in the DOM and runs the code in `useLayoutEffect`.
3. Your `useLayoutEffect` [measures the height](https://developer.mozilla.org/en-US/docs/Web/API/Element/getBoundingClientRect) of the tooltip content and triggers an immediate re-render.
4. `Tooltip` renders again with the real `tooltipHeight` (so the tooltip is correctly positioned).
5. React updates it in the DOM, and the browser finally displays the tooltip.

Hover over the buttons below and see how the tooltip adjusts its position depending on whether it fits:

<Sandpack>

```js
import ButtonWithTooltip from './ButtonWithTooltip.js';

export default function App() {
 return (
 <div>
 <ButtonWithTooltip
 tooltipContent={
 <div>
 This tooltip does not fit above the button.
 <br />
 This is why it's displayed below instead!
 </div>
 }
 >
 Hover over me (tooltip above)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 </div>
 );
}
```

```js src/ButtonWithTooltip.js
import { useState, useRef } from 'react';
import Tooltip from './Tooltip.js';

export default function ButtonWithTooltip({ tooltipContent, ...rest }) {
 const [targetRect, setTargetRect] = useState(null);
 const buttonRef = useRef(null);
 return (
 <>
 <button
 {...rest}
 ref={buttonRef}
 onPointerEnter={() => {
 const rect = buttonRef.current.getBoundingClientRect();
 setTargetRect({
 left: rect.left,
 top: rect.top,
 right: rect.right,
 bottom: rect.bottom,
 });
 }}
 onPointerLeave={() => {
 setTargetRect(null);
 }}
 />
 {targetRect !== null && (
 <Tooltip targetRect={targetRect}>
 {tooltipContent}
 </Tooltip>
 )
 }
 </>
 );
}
```

```js src/Tooltip.js active
import { useRef, useLayoutEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import TooltipContainer from './TooltipContainer.js';

export default function Tooltip({ children, targetRect }) {
 const ref = useRef(null);
 const [tooltipHeight, setTooltipHeight] = useState(0);

 useLayoutEffect(() => {
 const { height } = ref.current.getBoundingClientRect();
 setTooltipHeight(height);
 console.log('Measured tooltip height: ' + height);
 }, []);

 let tooltipX = 0;
 let tooltipY = 0;
 if (targetRect !== null) {
 tooltipX = targetRect.left;
 tooltipY = targetRect.top - tooltipHeight;
 if (tooltipY < 0) {
 // It doesn't fit above, so place below.
 tooltipY = targetRect.bottom;
 }
 }

 return createPortal(
 <TooltipContainer x={tooltipX} y={tooltipY} contentRef={ref}>
 {children}
 </TooltipContainer>,
 document.body
 );
}
```

```js src/TooltipContainer.js
export default function TooltipContainer({ children, x, y, contentRef }) {
 return (
 <div
 style={{
 position: 'absolute',
 pointerEvents: 'none',
 left: 0,
 top: 0,
 transform: `translate3d(${x}px, ${y}px, 0)`
 }}
 >
 <div ref={contentRef} className="tooltip">
 {children}
 </div>
 </div>
 );
}
```

```css
.tooltip {
 color: white;
 background: #222;
 border-radius: 4px;
 padding: 4px;
}
```

</Sandpack>

Notice that even though the `Tooltip` component has to render in two passes (first, with `tooltipHeight` initialized to `0` and then with the real measured height), you only see the final result. This is why you need `useLayoutEffect` instead of [`useEffect`](/reference/react/useEffect) for this example. Let's look at the difference in detail below.

<Recipes titleText="useLayoutEffect vs useEffect" titleId="examples">

#### `useLayoutEffect` blocks the browser from repainting {/*uselayouteffect-blocks-the-browser-from-repainting*/}

React guarantees that the code inside `useLayoutEffect` and any state updates scheduled inside it will be processed **before the browser repaints the screen.** This lets you render the tooltip, measure it, and re-render the tooltip again without the user noticing the first extra render. In other words, `useLayoutEffect` blocks the browser from painting.

<Sandpack>

```js
import ButtonWithTooltip from './ButtonWithTooltip.js';

export default function App() {
 return (
 <div>
 <ButtonWithTooltip
 tooltipContent={
 <div>
 This tooltip does not fit above the button.
 <br />
 This is why it's displayed below instead!
 </div>
 }
 >
 Hover over me (tooltip above)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 </div>
 );
}
```

```js src/ButtonWithTooltip.js
import { useState, useRef } from 'react';
import Tooltip from './Tooltip.js';

export default function ButtonWithTooltip({ tooltipContent, ...rest }) {
 const [targetRect, setTargetRect] = useState(null);
 const buttonRef = useRef(null);
 return (
 <>
 <button
 {...rest}
 ref={buttonRef}
 onPointerEnter={() => {
 const rect = buttonRef.current.getBoundingClientRect();
 setTargetRect({
 left: rect.left,
 top: rect.top,
 right: rect.right,
 bottom: rect.bottom,
 });
 }}
 onPointerLeave={() => {
 setTargetRect(null);
 }}
 />
 {targetRect !== null && (
 <Tooltip targetRect={targetRect}>
 {tooltipContent}
 </Tooltip>
 )
 }
 </>
 );
}
```

```js src/Tooltip.js active
import { useRef, useLayoutEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import TooltipContainer from './TooltipContainer.js';

export default function Tooltip({ children, targetRect }) {
 const ref = useRef(null);
 const [tooltipHeight, setTooltipHeight] = useState(0);

 useLayoutEffect(() => {
 const { height } = ref.current.getBoundingClientRect();
 setTooltipHeight(height);
 }, []);

 let tooltipX = 0;
 let tooltipY = 0;
 if (targetRect !== null) {
 tooltipX = targetRect.left;
 tooltipY = targetRect.top - tooltipHeight;
 if (tooltipY < 0) {
 // It doesn't fit above, so place below.
 tooltipY = targetRect.bottom;
 }
 }

 return createPortal(
 <TooltipContainer x={tooltipX} y={tooltipY} contentRef={ref}>
 {children}
 </TooltipContainer>,
 document.body
 );
}
```

```js src/TooltipContainer.js
export default function TooltipContainer({ children, x, y, contentRef }) {
 return (
 <div
 style={{
 position: 'absolute',
 pointerEvents: 'none',
 left: 0,
 top: 0,
 transform: `translate3d(${x}px, ${y}px, 0)`
 }}
 >
 <div ref={contentRef} className="tooltip">
 {children}
 </div>
 </div>
 );
}
```

```css
.tooltip {
 color: white;
 background: #222;
 border-radius: 4px;
 padding: 4px;
}
```

</Sandpack>

<Solution />

#### `useEffect` does not block the browser {/*useeffect-does-not-block-the-browser*/}

Here is the same example, but with [`useEffect`](/reference/react/useEffect) instead of `useLayoutEffect`. If you're on a slower device, you might notice that sometimes the tooltip "flickers" and you briefly see its initial position before the corrected position.

<Sandpack>

```js
import ButtonWithTooltip from './ButtonWithTooltip.js';

export default function App() {
 return (
 <div>
 <ButtonWithTooltip
 tooltipContent={
 <div>
 This tooltip does not fit above the button.
 <br />
 This is why it's displayed below instead!
 </div>
 }
 >
 Hover over me (tooltip above)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 </div>
 );
}
```

```js src/ButtonWithTooltip.js
import { useState, useRef } from 'react';
import Tooltip from './Tooltip.js';

export default function ButtonWithTooltip({ tooltipContent, ...rest }) {
 const [targetRect, setTargetRect] = useState(null);
 const buttonRef = useRef(null);
 return (
 <>
 <button
 {...rest}
 ref={buttonRef}
 onPointerEnter={() => {
 const rect = buttonRef.current.getBoundingClientRect();
 setTargetRect({
 left: rect.left,
 top: rect.top,
 right: rect.right,
 bottom: rect.bottom,
 });
 }}
 onPointerLeave={() => {
 setTargetRect(null);
 }}
 />
 {targetRect !== null && (
 <Tooltip targetRect={targetRect}>
 {tooltipContent}
 </Tooltip>
 )
 }
 </>
 );
}
```

```js src/Tooltip.js active
import { useRef, useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import TooltipContainer from './TooltipContainer.js';

export default function Tooltip({ children, targetRect }) {
 const ref = useRef(null);
 const [tooltipHeight, setTooltipHeight] = useState(0);

 useEffect(() => {
 const { height } = ref.current.getBoundingClientRect();
 setTooltipHeight(height);
 }, []);

 let tooltipX = 0;
 let tooltipY = 0;
 if (targetRect !== null) {
 tooltipX = targetRect.left;
 tooltipY = targetRect.top - tooltipHeight;
 if (tooltipY < 0) {
 // It doesn't fit above, so place below.
 tooltipY = targetRect.bottom;
 }
 }

 return createPortal(
 <TooltipContainer x={tooltipX} y={tooltipY} contentRef={ref}>
 {children}
 </TooltipContainer>,
 document.body
 );
}
```

```js src/TooltipContainer.js
export default function TooltipContainer({ children, x, y, contentRef }) {
 return (
 <div
 style={{
 position: 'absolute',
 pointerEvents: 'none',
 left: 0,
 top: 0,
 transform: `translate3d(${x}px, ${y}px, 0)`
 }}
 >
 <div ref={contentRef} className="tooltip">
 {children}
 </div>
 </div>
 );
}
```

```css
.tooltip {
 color: white;
 background: #222;
 border-radius: 4px;
 padding: 4px;
}
```

</Sandpack>

To make the bug easier to reproduce, this version adds an artificial delay during rendering. React will let the browser paint the screen before it processes the state update inside `useEffect`. As a result, the tooltip flickers:

<Sandpack>

```js
import ButtonWithTooltip from './ButtonWithTooltip.js';

export default function App() {
 return (
 <div>
 <ButtonWithTooltip
 tooltipContent={
 <div>
 This tooltip does not fit above the button.
 <br />
 This is why it's displayed below instead!
 </div>
 }
 >
 Hover over me (tooltip above)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 <div style={{ height: 50 }} />
 <ButtonWithTooltip
 tooltipContent={
 <div>This tooltip fits above the button</div>
 }
 >
 Hover over me (tooltip below)
 </ButtonWithTooltip>
 </div>
 );
}
```

```js src/ButtonWithTooltip.js
import { useState, useRef } from 'react';
import Tooltip from './Tooltip.js';

export default function ButtonWithTooltip({ tooltipContent, ...rest }) {
 const [targetRect, setTargetRect] = useState(null);
 const buttonRef = useRef(null);
 return (
 <>
 <button
 {...rest}
 ref={buttonRef}
 onPointerEnter={() => {
 const rect = buttonRef.current.getBoundingClientRect();
 setTargetRect({
 left: rect.left,
 top: rect.top,
 right: rect.right,
 bottom: rect.bottom,
 });
 }}
 onPointerLeave={() => {
 setTargetRect(null);
 }}
 />
 {targetRect !== null && (
 <Tooltip targetRect={targetRect}>
 {tooltipContent}
 </Tooltip>
 )
 }
 </>
 );
}
```

```js {expectedErrors: {'react-compiler': [10, 11]}} src/Tooltip.js active
import { useRef, useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import TooltipContainer from './TooltipContainer.js';

export default function Tooltip({ children, targetRect }) {
 const ref = useRef(null);
 const [tooltipHeight, setTooltipHeight] = useState(0);

 // This artificially slows down rendering
 let now = performance.now();
 while (performance.now() - now < 100) {
 // Do nothing for a bit...
 }

 useEffect(() => {
 const { height } = ref.current.getBoundingClientRect();
 setTooltipHeight(height);
 }, []);

 let tooltipX = 0;
 let tooltipY = 0;
 if (targetRect !== null) {
 tooltipX = targetRect.left;
 tooltipY = targetRect.top - tooltipHeight;
 if (tooltipY < 0) {
 // It doesn't fit above, so place below.
 tooltipY = targetRect.bottom;
 }
 }

 return createPortal(
 <TooltipContainer x={tooltipX} y={tooltipY} contentRef={ref}>
 {children}
 </TooltipContainer>,
 document.body
 );
}
```

```js src/TooltipContainer.js
export default function TooltipContainer({ children, x, y, contentRef }) {
 return (
 <div
 style={{
 position: 'absolute',
 pointerEvents: 'none',
 left: 0,
 top: 0,
 transform: `translate3d(${x}px, ${y}px, 0)`
 }}
 >
 <div ref={contentRef} className="tooltip">
 {children}
 </div>
 </div>
 );
}
```

```css
.tooltip {
 color: white;
 background: #222;
 border-radius: 4px;
 padding: 4px;
}
```

</Sandpack>

Edit this example to `useLayoutEffect` and observe that it blocks the paint even if rendering is slowed down.

<Solution />

</Recipes>

<Note>

Rendering in two passes and blocking the browser hurts performance. Try to avoid this when you can.

</Note>

---

## Troubleshooting {/*troubleshooting*/}

### I'm getting an error: "`useLayoutEffect` does nothing on the server" {/*im-getting-an-error-uselayouteffect-does-nothing-on-the-server*/}

The purpose of `useLayoutEffect` is to let your component [use layout information for rendering:](#measuring-layout-before-the-browser-repaints-the-screen)

1. Render the initial content.
2. Measure the layout *before the browser repaints the screen.*
3. Render the final content using the layout information you've read.

When you or your framework uses [server rendering](/reference/react-dom/server), your React app renders to HTML on the server for the initial render. This lets you show the initial HTML before the JavaScript code loads.

The problem is that on the server, there is no layout information.

In the [earlier example](#measuring-layout-before-the-browser-repaints-the-screen), the `useLayoutEffect` call in the `Tooltip` component lets it position itself correctly (either above or below content) depending on the content height. If you tried to render `Tooltip` as a part of the initial server HTML, this would be impossible to determine. On the server, there is no layout yet! So, even if you rendered it on the server, its position would "jump" on the client after the JavaScript loads and runs.

Usually, components that rely on layout information don't need to render on the server anyway. For example, it probably doesn't make sense to show a `Tooltip` during the initial render. It is triggered by a client interaction.

However, if you're running into this problem, you have a few different options:

- Replace `useLayoutEffect` with [`useEffect`.](/reference/react/useEffect) This tells React that it's okay to display the initial render result without blocking the paint (because the original HTML will become visible before your Effect runs).

- Alternatively, [mark your component as client-only.](/reference/react/Suspense#providing-a-fallback-for-server-errors-and-client-only-content) This tells React to replace its content up to the closest [`<Suspense>`](/reference/react/Suspense) boundary with a loading fallback (for example, a spinner or a glimmer) during server rendering.

- Alternatively, you can render a component with `useLayoutEffect` only after hydration. Keep a boolean `isMounted` state that's initialized to `false`, and set it to `true` inside a `useEffect` call. Your rendering logic can then be like `return isMounted ? <RealContent /> : <FallbackContent />`. On the server and during the hydration, the user will see `FallbackContent` which should not call `useLayoutEffect`. Then React will replace it with `RealContent` which runs on the client only and can include `useLayoutEffect` calls.

- If you synchronize your component with an external data store and rely on `useLayoutEffect` for different reasons than measuring layout, consider [`useSyncExternalStore`](/reference/react/useSyncExternalStore) instead which [supports server rendering.](/reference/react/useSyncExternalStore#adding-support-for-server-rendering)

---
title: useMemo
---

<Intro>

`useMemo` is a React Hook that lets you cache the result of a calculation between re-renders.

```js
const cachedValue = useMemo(calculateValue, dependencies)
```

</Intro>

<Note>

[React Compiler](/learn/react-compiler) automatically memoizes values and functions, reducing the need for manual `useMemo` calls. You can use the compiler to handle memoization automatically.

</Note>

<InlineToc />

---

## Reference {/*reference*/}

### `useMemo(calculateValue, dependencies)` {/*usememo*/}

Call `useMemo` at the top level of your component to cache a calculation between re-renders:

```js
import { useMemo } from 'react';

function TodoList({ todos, tab }) {
 const visibleTodos = useMemo(
 () => filterTodos(todos, tab),
 [todos, tab]
 );
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `calculateValue`: The function calculating the value that you want to cache. It should be pure, should take no arguments, and should return a value of any type. React will call your function during the initial render. On next renders, React will return the same value again if the `dependencies` have not changed since the last render. Otherwise, it will call `calculateValue`, return its result, and store it so it can be reused later.

* `dependencies`: The list of all reactive values referenced inside of the `calculateValue` code. Reactive values include props, state, and all the variables and functions declared directly inside your component body. If your linter is [configured for React](/learn/editor-setup#linting), it will verify that every reactive value is correctly specified as a dependency. The list of dependencies must have a constant number of items and be written inline like `[dep1, dep2, dep3]`. React will compare each dependency with its previous value using the [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison.

#### Returns {/*returns*/}

On the initial render, `useMemo` returns the result of calling `calculateValue` with no arguments.

During next renders, it will either return an already stored value from the last render (if the dependencies haven't changed), or call `calculateValue` again, and return the result that `calculateValue` has returned.

#### Caveats {/*caveats*/}

* `useMemo` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.
* In Strict Mode, React will **call your calculation function twice** in order to [help you find accidental impurities.](#my-calculation-runs-twice-on-every-re-render) This is development-only behavior and does not affect production. If your calculation function is pure (as it should be), this should not affect your logic. The result from one of the calls will be ignored.
* React **will not throw away the cached value unless there is a specific reason to do that.** For example, in development, React throws away the cache when you edit the file of your component. Both in development and in production, React will throw away the cache if your component suspends during the initial mount. In the future, React may add more features that take advantage of throwing away the cache--for example, if React adds built-in support for virtualized lists in the future, it would make sense to throw away the cache for items that scroll out of the virtualized table viewport. This should be fine if you rely on `useMemo` solely as a performance optimization. Otherwise, a [state variable](/reference/react/useState#avoiding-recreating-the-initial-state) or a [ref](/reference/react/useRef#avoiding-recreating-the-ref-contents) may be more appropriate.

<Note>

Caching return values like this is also known as [*memoization*,](https://en.wikipedia.org/wiki/Memoization) which is why this Hook is called `useMemo`.

</Note>

---

## Usage {/*usage*/}

### Skipping expensive recalculations {/*skipping-expensive-recalculations*/}

To cache a calculation between re-renders, wrap it in a `useMemo` call at the top level of your component:

```js [[3, 4, "visibleTodos"], [1, 4, "() => filterTodos(todos, tab)"], [2, 4, "[todos, tab]"]]
import { useMemo } from 'react';

function TodoList({ todos, tab, theme }) {
 const visibleTodos = useMemo(() => filterTodos(todos, tab), [todos, tab]);
 // ...
}
```

You need to pass two things to `useMemo`:

1. A <CodeStep step={1}>calculation function</CodeStep> that takes no arguments, like `() =>`, and returns what you wanted to calculate.
2. A <CodeStep step={2}>list of dependencies</CodeStep> including every value within your component that's used inside your calculation.

On the initial render, the <CodeStep step={3}>value</CodeStep> you'll get from `useMemo` will be the result of calling your <CodeStep step={1}>calculation</CodeStep>.

On every subsequent render, React will compare the <CodeStep step={2}>dependencies</CodeStep> with the dependencies you passed during the last render. If none of the dependencies have changed (compared with [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is)), `useMemo` will return the value you already calculated before. Otherwise, React will re-run your calculation and return the new value.

In other words, `useMemo` caches a calculation result between re-renders until its dependencies change.

**Let's walk through an example to see when this is useful.**

By default, React will re-run the entire body of your component every time that it re-renders. For example, if this `TodoList` updates its state or receives new props from its parent, the `filterTodos` function will re-run:

```js {2}
function TodoList({ todos, tab, theme }) {
 const visibleTodos = filterTodos(todos, tab);
 // ...
}
```

Usually, this isn't a problem because most calculations are very fast. However, if you're filtering or transforming a large array, or doing some expensive computation, you might want to skip doing it again if data hasn't changed. If both `todos` and `tab` are the same as they were during the last render, wrapping the calculation in `useMemo` like earlier lets you reuse `visibleTodos` you've already calculated before.

This type of caching is called *[memoization.](https://en.wikipedia.org/wiki/Memoization)*

<Note>

**You should only rely on `useMemo` as a performance optimization.** If your code doesn't work without it, find the underlying problem and fix it first. Then you may add `useMemo` to improve performance.

</Note>

<DeepDive>

#### How to tell if a calculation is expensive? {/*how-to-tell-if-a-calculation-is-expensive*/}

In general, unless you're creating or looping over thousands of objects, it's probably not expensive. If you want to get more confidence, you can add a console log to measure the time spent in a piece of code:

```js {1,3}
console.time('filter array');
const visibleTodos = filterTodos(todos, tab);
console.timeEnd('filter array');
```

Perform the interaction you're measuring (for example, typing into the input). You will then see logs like `filter array: 0.15ms` in your console. If the overall logged time adds up to a significant amount (say, `1ms` or more), it might make sense to memoize that calculation. As an experiment, you can then wrap the calculation in `useMemo` to verify whether the total logged time has decreased for that interaction or not:

```js
console.time('filter array');
const visibleTodos = useMemo(() => {
 return filterTodos(todos, tab); // Skipped if todos and tab haven't changed
}, [todos, tab]);
console.timeEnd('filter array');
```

`useMemo` won't make the *first* render faster. It only helps you skip unnecessary work on updates.

Keep in mind that your machine is probably faster than your users' so it's a good idea to test the performance with an artificial slowdown. For example, Chrome offers a [CPU Throttling](https://developer.chrome.com/blog/new-in-devtools-61/#throttling) option for this.

Also note that measuring performance in development will not give you the most accurate results. (For example, when [Strict Mode](/reference/react/StrictMode) is on, you will see each component render twice rather than once.) To get the most accurate timings, build your app for production and test it on a device like your users have.

</DeepDive>

<DeepDive>

#### Should you add useMemo everywhere? {/*should-you-add-usememo-everywhere*/}

If your app is like this site, and most interactions are coarse (like replacing a page or an entire section), memoization is usually unnecessary. On the other hand, if your app is more like a drawing editor, and most interactions are granular (like moving shapes), then you might find memoization very helpful.

Optimizing with `useMemo` is only valuable in a few cases:

- The calculation you're putting in `useMemo` is noticeably slow, and its dependencies rarely change.
- You pass it as a prop to a component wrapped in [`memo`.](/reference/react/memo) You want to skip re-rendering if the value hasn't changed. Memoization lets your component re-render only when dependencies aren't the same.
- The value you're passing is later used as a dependency of some Hook. For example, maybe another `useMemo` calculation value depends on it. Or maybe you are depending on this value from [`useEffect.`](/reference/react/useEffect)

There is no benefit to wrapping a calculation in `useMemo` in other cases. There is no significant harm to doing that either, so some teams choose to not think about individual cases, and memoize as much as possible. The downside of this approach is that code becomes less readable. Also, not all memoization is effective: a single value that's "always new" is enough to break memoization for an entire component.

**In practice, you can make a lot of memoization unnecessary by following a few principles:**

1. When a component visually wraps other components, let it [accept JSX as children.](/learn/passing-props-to-a-component#passing-jsx-as-children) This way, when the wrapper component updates its own state, React knows that its children don't need to re-render.
1. Prefer local state and don't [lift state up](/learn/sharing-state-between-components) any further than necessary. For example, don't keep transient state like forms and whether an item is hovered at the top of your tree or in a global state library.
1. Keep your [rendering logic pure.](/learn/keeping-components-pure) If re-rendering a component causes a problem or produces some noticeable visual artifact, it's a bug in your component! Fix the bug instead of adding memoization.
1. Avoid [unnecessary Effects that update state.](/learn/you-might-not-need-an-effect) Most performance problems in React apps are caused by chains of updates originating from Effects that cause your components to render over and over.
1. Try to [remove unnecessary dependencies from your Effects.](/learn/removing-effect-dependencies) For example, instead of memoization, it's often simpler to move some object or a function inside an Effect or outside the component.

If a specific interaction still feels laggy, [use the React Developer Tools profiler](https://legacy.reactjs.org/blog/2018/09/10/introducing-the-react-profiler.html) to see which components would benefit the most from memoization, and add memoization where needed. These principles make your components easier to debug and understand, so it's good to follow them in any case. In the long term, we're researching [doing granular memoization automatically](https://www.youtube.com/watch?v=lGEMwh32soc) to solve this once and for all.

</DeepDive>

<Recipes titleText="The difference between useMemo and calculating a value directly" titleId="examples-recalculation">

#### Skipping recalculation with `useMemo` {/*skipping-recalculation-with-usememo*/}

In this example, the `filterTodos` implementation is **artificially slowed down** so that you can see what happens when some JavaScript function you're calling during rendering is genuinely slow. Try switching the tabs and toggling the theme.

Switching the tabs feels slow because it forces the slowed down `filterTodos` to re-execute. That's expected because the `tab` has changed, and so the entire calculation *needs* to re-run. (If you're curious why it runs twice, it's explained [here.](#my-calculation-runs-twice-on-every-re-render))

Toggle the theme. **Thanks to `useMemo`, it's fast despite the artificial slowdown!** The slow `filterTodos` call was skipped because both `todos` and `tab` (which you pass as dependencies to `useMemo`) haven't changed since the last render.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import { createTodos } from './utils.js';
import TodoList from './TodoList.js';

const todos = createTodos();

export default function App() {
 const [tab, setTab] = useState('all');
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <button onClick={() => setTab('all')}>
 All
 </button>
 <button onClick={() => setTab('active')}>
 Active
 </button>
 <button onClick={() => setTab('completed')}>
 Completed
 </button>
 <br />
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <TodoList
 todos={todos}
 tab={tab}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}

```

```js src/TodoList.js active
import { useMemo } from 'react';
import { filterTodos } from './utils.js'

export default function TodoList({ todos, theme, tab }) {
 const visibleTodos = useMemo(
 () => filterTodos(todos, tab),
 [todos, tab]
 );
 return (
 <div className={theme}>
 <p><b>Note: <code>filterTodos</code> is artificially slowed down!</b></p>
 <ul>
 {visibleTodos.map(todo => (
 <li key={todo.id}>
 {todo.completed ?
 <s>{todo.text}</s> :
 todo.text
 }
 </li>
 ))}
 </ul>
 </div>
 );
}
```

```js src/utils.js
export function createTodos() {
 const todos = [];
 for (let i = 0; i < 50; i++) {
 todos.push({
 id: i,
 text: "Todo " + (i + 1),
 completed: Math.random() > 0.5
 });
 }
 return todos;
}

export function filterTodos(todos, tab) {
 console.log('[ARTIFICIALLY SLOW] Filtering ' + todos.length + ' todos for "' + tab + '" tab.');
 let startTime = performance.now();
 while (performance.now() - startTime < 500) {
 // Do nothing for 500 ms to emulate extremely slow code
 }

 return todos.filter(todo => {
 if (tab === 'all') {
 return true;
 } else if (tab === 'active') {
 return !todo.completed;
 } else if (tab === 'completed') {
 return todo.completed;
 }
 });
}
```

```css
label {
 display: block;
 margin-top: 10px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

<Solution />

#### Always recalculating a value {/*always-recalculating-a-value*/}

In this example, the `filterTodos` implementation is also **artificially slowed down** so that you can see what happens when some JavaScript function you're calling during rendering is genuinely slow. Try switching the tabs and toggling the theme.

Unlike in the previous example, toggling the theme is also slow now! This is because **there is no `useMemo` call in this version,** so the artificially slowed down `filterTodos` gets called on every re-render. It is called even if only `theme` has changed.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import { createTodos } from './utils.js';
import TodoList from './TodoList.js';

const todos = createTodos();

export default function App() {
 const [tab, setTab] = useState('all');
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <button onClick={() => setTab('all')}>
 All
 </button>
 <button onClick={() => setTab('active')}>
 Active
 </button>
 <button onClick={() => setTab('completed')}>
 Completed
 </button>
 <br />
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <TodoList
 todos={todos}
 tab={tab}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}

```

```js src/TodoList.js active
import { filterTodos } from './utils.js'

export default function TodoList({ todos, theme, tab }) {
 const visibleTodos = filterTodos(todos, tab);
 return (
 <div className={theme}>
 <ul>
 <p><b>Note: <code>filterTodos</code> is artificially slowed down!</b></p>
 {visibleTodos.map(todo => (
 <li key={todo.id}>
 {todo.completed ?
 <s>{todo.text}</s> :
 todo.text
 }
 </li>
 ))}
 </ul>
 </div>
 );
}
```

```js src/utils.js
export function createTodos() {
 const todos = [];
 for (let i = 0; i < 50; i++) {
 todos.push({
 id: i,
 text: "Todo " + (i + 1),
 completed: Math.random() > 0.5
 });
 }
 return todos;
}

export function filterTodos(todos, tab) {
 console.log('[ARTIFICIALLY SLOW] Filtering ' + todos.length + ' todos for "' + tab + '" tab.');
 let startTime = performance.now();
 while (performance.now() - startTime < 500) {
 // Do nothing for 500 ms to emulate extremely slow code
 }

 return todos.filter(todo => {
 if (tab === 'all') {
 return true;
 } else if (tab === 'active') {
 return !todo.completed;
 } else if (tab === 'completed') {
 return todo.completed;
 }
 });
}
```

```css
label {
 display: block;
 margin-top: 10px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

However, here is the same code **with the artificial slowdown removed.** Does the lack of `useMemo` feel noticeable or not?

<Sandpack>

```js src/App.js
import { useState } from 'react';
import { createTodos } from './utils.js';
import TodoList from './TodoList.js';

const todos = createTodos();

export default function App() {
 const [tab, setTab] = useState('all');
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <button onClick={() => setTab('all')}>
 All
 </button>
 <button onClick={() => setTab('active')}>
 Active
 </button>
 <button onClick={() => setTab('completed')}>
 Completed
 </button>
 <br />
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <TodoList
 todos={todos}
 tab={tab}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}

```

```js src/TodoList.js active
import { filterTodos } from './utils.js'

export default function TodoList({ todos, theme, tab }) {
 const visibleTodos = filterTodos(todos, tab);
 return (
 <div className={theme}>
 <ul>
 {visibleTodos.map(todo => (
 <li key={todo.id}>
 {todo.completed ?
 <s>{todo.text}</s> :
 todo.text
 }
 </li>
 ))}
 </ul>
 </div>
 );
}
```

```js src/utils.js
export function createTodos() {
 const todos = [];
 for (let i = 0; i < 50; i++) {
 todos.push({
 id: i,
 text: "Todo " + (i + 1),
 completed: Math.random() > 0.5
 });
 }
 return todos;
}

export function filterTodos(todos, tab) {
 console.log('Filtering ' + todos.length + ' todos for "' + tab + '" tab.');

 return todos.filter(todo => {
 if (tab === 'all') {
 return true;
 } else if (tab === 'active') {
 return !todo.completed;
 } else if (tab === 'completed') {
 return todo.completed;
 }
 });
}
```

```css
label {
 display: block;
 margin-top: 10px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

Quite often, code without memoization works fine. If your interactions are fast enough, you might not need memoization.

You can try increasing the number of todo items in `utils.js` and see how the behavior changes. This particular calculation wasn't very expensive to begin with, but if the number of todos grows significantly, most of the overhead will be in re-rendering rather than in the filtering. Keep reading below to see how you can optimize re-rendering with `useMemo`.

<Solution />

</Recipes>

---

### Skipping re-rendering of components {/*skipping-re-rendering-of-components*/}

In some cases, `useMemo` can also help you optimize performance of re-rendering child components. To illustrate this, let's say this `TodoList` component passes the `visibleTodos` as a prop to the child `List` component:

```js {5}
export default function TodoList({ todos, tab, theme }) {
 // ...
 return (
 <div className={theme}>
 <List items={visibleTodos} />
 </div>
 );
}
```

You've noticed that toggling the `theme` prop freezes the app for a moment, but if you remove `<List />` from your JSX, it feels fast. This tells you that it's worth trying to optimize the `List` component.

**By default, when a component re-renders, React re-renders all of its children recursively.** This is why, when `TodoList` re-renders with a different `theme`, the `List` component *also* re-renders. This is fine for components that don't require much calculation to re-render. But if you've verified that a re-render is slow, you can tell `List` to skip re-rendering when its props are the same as on last render by wrapping it in [`memo`:](/reference/react/memo)

```js {3,5}
import { memo } from 'react';

const List = memo(function List({ items }) {
 // ...
});
```

**With this change, `List` will skip re-rendering if all of its props are the *same* as on the last render.** This is where caching the calculation becomes important! Imagine that you calculated `visibleTodos` without `useMemo`:

```js {2-3,6-7}
export default function TodoList({ todos, tab, theme }) {
 // Every time the theme changes, this will be a different array...
 const visibleTodos = filterTodos(todos, tab);
 return (
 <div className={theme}>
 {/* ... so List's props will never be the same, and it will re-render every time */}
 <List items={visibleTodos} />
 </div>
 );
}
```

**In the above example, the `filterTodos` function always creates a *different* array,** similar to how the `{}` object literal always creates a new object. Normally, this wouldn't be a problem, but it means that `List` props will never be the same, and your [`memo`](/reference/react/memo) optimization won't work. This is where `useMemo` comes in handy:

```js {2-3,5,9-10}
export default function TodoList({ todos, tab, theme }) {
 // Tell React to cache your calculation between re-renders...
 const visibleTodos = useMemo(
 () => filterTodos(todos, tab),
 [todos, tab] // ...so as long as these dependencies don't change...
 );
 return (
 <div className={theme}>
 {/* ...List will receive the same props and can skip re-rendering */}
 <List items={visibleTodos} />
 </div>
 );
}
```

**By wrapping the `visibleTodos` calculation in `useMemo`, you ensure that it has the *same* value between the re-renders** (until dependencies change). You don't *have to* wrap a calculation in `useMemo` unless you do it for some specific reason. In this example, the reason is that you pass it to a component wrapped in [`memo`,](/reference/react/memo) and this lets it skip re-rendering. There are a few other reasons to add `useMemo` which are described further on this page.

<DeepDive>

#### Memoizing individual JSX nodes {/*memoizing-individual-jsx-nodes*/}

Instead of wrapping `List` in [`memo`](/reference/react/memo), you could wrap the `<List />` JSX node itself in `useMemo`:

```js {3,6}
export default function TodoList({ todos, tab, theme }) {
 const visibleTodos = useMemo(() => filterTodos(todos, tab), [todos, tab]);
 const children = useMemo(() => <List items={visibleTodos} />, [visibleTodos]);
 return (
 <div className={theme}>
 {children}
 </div>
 );
}
```

The behavior would be the same. If the `visibleTodos` haven't changed, `List` won't be re-rendered.

A JSX node like `<List items={visibleTodos} />` is an object like `{ type: List, props: { items: visibleTodos } }`. Creating this object is very cheap, but React doesn't know whether its contents is the same as last time or not. This is why by default, React will re-render the `List` component.

However, if React sees the same exact JSX as during the previous render, it won't try to re-render your component. This is because JSX nodes are [immutable.](https://en.wikipedia.org/wiki/Immutable_object) A JSX node object could not have changed over time, so React knows it's safe to skip a re-render. However, for this to work, the node has to *actually be the same object*, not merely look the same in code. This is what `useMemo` does in this example.

Manually wrapping JSX nodes into `useMemo` is not convenient. For example, you can't do this conditionally. This is usually why you would wrap components with [`memo`](/reference/react/memo) instead of wrapping JSX nodes.

</DeepDive>

<Recipes titleText="The difference between skipping re-renders and always re-rendering" titleId="examples-rerendering">

#### Skipping re-rendering with `useMemo` and `memo` {/*skipping-re-rendering-with-usememo-and-memo*/}

In this example, the `List` component is **artificially slowed down** so that you can see what happens when a React component you're rendering is genuinely slow. Try switching the tabs and toggling the theme.

Switching the tabs feels slow because it forces the slowed down `List` to re-render. That's expected because the `tab` has changed, and so you need to reflect the user's new choice on the screen.

Next, try toggling the theme. **Thanks to `useMemo` together with [`memo`](/reference/react/memo), it’s fast despite the artificial slowdown!** The `List` skipped re-rendering because the `visibleTodos` array has not changed since the last render. The `visibleTodos` array has not changed because both `todos` and `tab` (which you pass as dependencies to `useMemo`) haven't changed since the last render.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import { createTodos } from './utils.js';
import TodoList from './TodoList.js';

const todos = createTodos();

export default function App() {
 const [tab, setTab] = useState('all');
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <button onClick={() => setTab('all')}>
 All
 </button>
 <button onClick={() => setTab('active')}>
 Active
 </button>
 <button onClick={() => setTab('completed')}>
 Completed
 </button>
 <br />
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <TodoList
 todos={todos}
 tab={tab}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}
```

```js src/TodoList.js active
import { useMemo } from 'react';
import List from './List.js';
import { filterTodos } from './utils.js'

export default function TodoList({ todos, theme, tab }) {
 const visibleTodos = useMemo(
 () => filterTodos(todos, tab),
 [todos, tab]
 );
 return (
 <div className={theme}>
 <p><b>Note: <code>List</code> is artificially slowed down!</b></p>
 <List items={visibleTodos} />
 </div>
 );
}
```

```js {expectedErrors: {'react-compiler': [5, 6]}} src/List.js
import { memo } from 'react';

const List = memo(function List({ items }) {
 console.log('[ARTIFICIALLY SLOW] Rendering <List /> with ' + items.length + ' items');
 let startTime = performance.now();
 while (performance.now() - startTime < 500) {
 // Do nothing for 500 ms to emulate extremely slow code
 }

 return (
 <ul>
 {items.map(item => (
 <li key={item.id}>
 {item.completed ?
 <s>{item.text}</s> :
 item.text
 }
 </li>
 ))}
 </ul>
 );
});

export default List;
```

```js src/utils.js
export function createTodos() {
 const todos = [];
 for (let i = 0; i < 50; i++) {
 todos.push({
 id: i,
 text: "Todo " + (i + 1),
 completed: Math.random() > 0.5
 });
 }
 return todos;
}

export function filterTodos(todos, tab) {
 return todos.filter(todo => {
 if (tab === 'all') {
 return true;
 } else if (tab === 'active') {
 return !todo.completed;
 } else if (tab === 'completed') {
 return todo.completed;
 }
 });
}
```

```css
label {
 display: block;
 margin-top: 10px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

<Solution />

#### Always re-rendering a component {/*always-re-rendering-a-component*/}

In this example, the `List` implementation is also **artificially slowed down** so that you can see what happens when some React component you're rendering is genuinely slow. Try switching the tabs and toggling the theme.

Unlike in the previous example, toggling the theme is also slow now! This is because **there is no `useMemo` call in this version,** so the `visibleTodos` is always a different array, and the slowed down `List` component can't skip re-rendering.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import { createTodos } from './utils.js';
import TodoList from './TodoList.js';

const todos = createTodos();

export default function App() {
 const [tab, setTab] = useState('all');
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <button onClick={() => setTab('all')}>
 All
 </button>
 <button onClick={() => setTab('active')}>
 Active
 </button>
 <button onClick={() => setTab('completed')}>
 Completed
 </button>
 <br />
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <TodoList
 todos={todos}
 tab={tab}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}
```

```js src/TodoList.js active
import List from './List.js';
import { filterTodos } from './utils.js'

export default function TodoList({ todos, theme, tab }) {
 const visibleTodos = filterTodos(todos, tab);
 return (
 <div className={theme}>
 <p><b>Note: <code>List</code> is artificially slowed down!</b></p>
 <List items={visibleTodos} />
 </div>
 );
}
```

```js {expectedErrors: {'react-compiler': [5, 6]}} src/List.js
import { memo } from 'react';

const List = memo(function List({ items }) {
 console.log('[ARTIFICIALLY SLOW] Rendering <List /> with ' + items.length + ' items');
 let startTime = performance.now();
 while (performance.now() - startTime < 500) {
 // Do nothing for 500 ms to emulate extremely slow code
 }

 return (
 <ul>
 {items.map(item => (
 <li key={item.id}>
 {item.completed ?
 <s>{item.text}</s> :
 item.text
 }
 </li>
 ))}
 </ul>
 );
});

export default List;
```

```js src/utils.js
export function createTodos() {
 const todos = [];
 for (let i = 0; i < 50; i++) {
 todos.push({
 id: i,
 text: "Todo " + (i + 1),
 completed: Math.random() > 0.5
 });
 }
 return todos;
}

export function filterTodos(todos, tab) {
 return todos.filter(todo => {
 if (tab === 'all') {
 return true;
 } else if (tab === 'active') {
 return !todo.completed;
 } else if (tab === 'completed') {
 return todo.completed;
 }
 });
}
```

```css
label {
 display: block;
 margin-top: 10px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

However, here is the same code **with the artificial slowdown removed.** Does the lack of `useMemo` feel noticeable or not?

<Sandpack>

```js src/App.js
import { useState } from 'react';
import { createTodos } from './utils.js';
import TodoList from './TodoList.js';

const todos = createTodos();

export default function App() {
 const [tab, setTab] = useState('all');
 const [isDark, setIsDark] = useState(false);
 return (
 <>
 <button onClick={() => setTab('all')}>
 All
 </button>
 <button onClick={() => setTab('active')}>
 Active
 </button>
 <button onClick={() => setTab('completed')}>
 Completed
 </button>
 <br />
 <label>
 <input
 type="checkbox"
 checked={isDark}
 onChange={e => setIsDark(e.target.checked)}
 />
 Dark mode
 </label>
 <hr />
 <TodoList
 todos={todos}
 tab={tab}
 theme={isDark ? 'dark' : 'light'}
 />
 </>
 );
}
```

```js src/TodoList.js active
import List from './List.js';
import { filterTodos } from './utils.js'

export default function TodoList({ todos, theme, tab }) {
 const visibleTodos = filterTodos(todos, tab);
 return (
 <div className={theme}>
 <List items={visibleTodos} />
 </div>
 );
}
```

```js src/List.js
import { memo } from 'react';

function List({ items }) {
 return (
 <ul>
 {items.map(item => (
 <li key={item.id}>
 {item.completed ?
 <s>{item.text}</s> :
 item.text
 }
 </li>
 ))}
 </ul>
 );
}

export default memo(List);
```

```js src/utils.js
export function createTodos() {
 const todos = [];
 for (let i = 0; i < 50; i++) {
 todos.push({
 id: i,
 text: "Todo " + (i + 1),
 completed: Math.random() > 0.5
 });
 }
 return todos;
}

export function filterTodos(todos, tab) {
 return todos.filter(todo => {
 if (tab === 'all') {
 return true;
 } else if (tab === 'active') {
 return !todo.completed;
 } else if (tab === 'completed') {
 return todo.completed;
 }
 });
}
```

```css
label {
 display: block;
 margin-top: 10px;
}

.dark {
 background-color: black;
 color: white;
}

.light {
 background-color: white;
 color: black;
}
```

</Sandpack>

Quite often, code without memoization works fine. If your interactions are fast enough, you don't need memoization.

Keep in mind that you need to run React in production mode, disable [React Developer Tools](/learn/react-developer-tools), and use devices similar to the ones your app's users have in order to get a realistic sense of what's actually slowing down your app.

<Solution />

</Recipes>

---

### Preventing an Effect from firing too often {/*preventing-an-effect-from-firing-too-often*/}

Sometimes, you might want to use a value inside an [Effect:](/learn/synchronizing-with-effects)

```js {4-7,10}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 const options = {
 serverUrl: 'https://localhost:1234',
 roomId: roomId
 }

 useEffect(() => {
 const connection = createConnection(options);
 connection.connect();
 // ...
```

This creates a problem. [Every reactive value must be declared as a dependency of your Effect.](/learn/lifecycle-of-reactive-effects#react-verifies-that-you-specified-every-reactive-value-as-a-dependency) However, if you declare `options` as a dependency, it will cause your Effect to constantly reconnect to the chat room:

```js {5}
 useEffect(() => {
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [options]); // 🔴 Problem: This dependency changes on every render
 // ...
```

To solve this, you can wrap the object you need to call from an Effect in `useMemo`:

```js {4-9,16}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 const options = useMemo(() => {
 return {
 serverUrl: 'https://localhost:1234',
 roomId: roomId
 };
 }, [roomId]); // ✅ Only changes when roomId changes

 useEffect(() => {
 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [options]); // ✅ Only changes when options changes
 // ...
```

This ensures that the `options` object is the same between re-renders if `useMemo` returns the cached object.

However, since `useMemo` is performance optimization, not a semantic guarantee, React may throw away the cached value if [there is a specific reason to do that](#caveats). This will also cause the effect to re-fire, **so it's even better to remove the need for a function dependency** by moving your object *inside* the Effect:

```js {5-8,13}
function ChatRoom({ roomId }) {
 const [message, setMessage] = useState('');

 useEffect(() => {
 const options = { // ✅ No need for useMemo or object dependencies!
 serverUrl: 'https://localhost:1234',
 roomId: roomId
 }

 const connection = createConnection(options);
 connection.connect();
 return () => connection.disconnect();
 }, [roomId]); // ✅ Only changes when roomId changes
 // ...
```

Now your code is simpler and doesn't need `useMemo`. [Learn more about removing Effect dependencies.](/learn/removing-effect-dependencies#move-dynamic-objects-and-functions-inside-your-effect)

### Memoizing a dependency of another Hook {/*memoizing-a-dependency-of-another-hook*/}

Suppose you have a calculation that depends on an object created directly in the component body:

```js {2}
function Dropdown({ allItems, text }) {
 const searchOptions = { matchMode: 'whole-word', text };

 const visibleItems = useMemo(() => {
 return searchItems(allItems, searchOptions);
 }, [allItems, searchOptions]); // 🚩 Caution: Dependency on an object created in the component body
 // ...
```

Depending on an object like this defeats the point of memoization. When a component re-renders, all of the code directly inside the component body runs again. **The lines of code creating the `searchOptions` object will also run on every re-render.** Since `searchOptions` is a dependency of your `useMemo` call, and it's different every time, React knows the dependencies are different, and recalculate `searchItems` every time.

To fix this, you could memoize the `searchOptions` object *itself* before passing it as a dependency:

```js {2-4}
function Dropdown({ allItems, text }) {
 const searchOptions = useMemo(() => {
 return { matchMode: 'whole-word', text };
 }, [text]); // ✅ Only changes when text changes

 const visibleItems = useMemo(() => {
 return searchItems(allItems, searchOptions);
 }, [allItems, searchOptions]); // ✅ Only changes when allItems or searchOptions changes
 // ...
```

In the example above, if the `text` did not change, the `searchOptions` object also won't change. However, an even better fix is to move the `searchOptions` object declaration *inside* of the `useMemo` calculation function:

```js {3}
function Dropdown({ allItems, text }) {
 const visibleItems = useMemo(() => {
 const searchOptions = { matchMode: 'whole-word', text };
 return searchItems(allItems, searchOptions);
 }, [allItems, text]); // ✅ Only changes when allItems or text changes
 // ...
```

Now your calculation depends on `text` directly (which is a string and can't "accidentally" become different).

---

### Memoizing a function {/*memoizing-a-function*/}

Suppose the `Form` component is wrapped in [`memo`.](/reference/react/memo) You want to pass a function to it as a prop:

```js {2-7}
export default function ProductPage({ productId, referrer }) {
 function handleSubmit(orderDetails) {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails
 });
 }

 return <Form onSubmit={handleSubmit} />;
}
```

Just as `{}` creates a different object, function declarations like `function() {}` and expressions like `() => {}` produce a *different* function on every re-render. By itself, creating a new function is not a problem. This is not something to avoid! However, if the `Form` component is memoized, presumably you want to skip re-rendering it when no props have changed. A prop that is *always* different would defeat the point of memoization.

To memoize a function with `useMemo`, your calculation function would have to return another function:

```js {2-3,8-9}
export default function Page({ productId, referrer }) {
 const handleSubmit = useMemo(() => {
 return (orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails
 });
 };
 }, [productId, referrer]);

 return <Form onSubmit={handleSubmit} />;
}
```

This looks clunky! **Memoizing functions is common enough that React has a built-in Hook specifically for that. Wrap your functions into [`useCallback`](/reference/react/useCallback) instead of `useMemo`** to avoid having to write an extra nested function:

```js {2,7}
export default function Page({ productId, referrer }) {
 const handleSubmit = useCallback((orderDetails) => {
 post('/product/' + productId + '/buy', {
 referrer,
 orderDetails
 });
 }, [productId, referrer]);

 return <Form onSubmit={handleSubmit} />;
}
```

The two examples above are completely equivalent. The only benefit to `useCallback` is that it lets you avoid writing an extra nested function inside. It doesn't do anything else. [Read more about `useCallback`.](/reference/react/useCallback)

---

## Troubleshooting {/*troubleshooting*/}

### My calculation runs twice on every re-render {/*my-calculation-runs-twice-on-every-re-render*/}

In [Strict Mode](/reference/react/StrictMode), React will call some of your functions twice instead of once:

```js {2,5,6}
function TodoList({ todos, tab }) {
 // This component function will run twice for every render.

 const visibleTodos = useMemo(() => {
 // This calculation will run twice if any of the dependencies change.
 return filterTodos(todos, tab);
 }, [todos, tab]);

 // ...
```

This is expected and shouldn't break your code.

This **development-only** behavior helps you [keep components pure.](/learn/keeping-components-pure) React uses the result of one of the calls, and ignores the result of the other call. As long as your component and calculation functions are pure, this shouldn't affect your logic. However, if they are accidentally impure, this helps you notice and fix the mistake.

For example, this impure calculation function mutates an array you received as a prop:

```js {2-3}
 const visibleTodos = useMemo(() => {
 // 🚩 Mistake: mutating a prop
 todos.push({ id: 'last', text: 'Go for a walk!' });
 const filtered = filterTodos(todos, tab);
 return filtered;
 }, [todos, tab]);
```

React calls your function twice, so you'd notice the todo is added twice. Your calculation shouldn't change any existing objects, but it's okay to change any *new* objects you created during the calculation. For example, if the `filterTodos` function always returns a *different* array, you can mutate *that* array instead:

```js {3,4}
 const visibleTodos = useMemo(() => {
 const filtered = filterTodos(todos, tab);
 // ✅ Correct: mutating an object you created during the calculation
 filtered.push({ id: 'last', text: 'Go for a walk!' });
 return filtered;
 }, [todos, tab]);
```

Read [keeping components pure](/learn/keeping-components-pure) to learn more about purity.

Also, check out the guides on [updating objects](/learn/updating-objects-in-state) and [updating arrays](/learn/updating-arrays-in-state) without mutation.

---

### My `useMemo` call is supposed to return an object, but returns undefined {/*my-usememo-call-is-supposed-to-return-an-object-but-returns-undefined*/}

This code doesn't work:

```js {1-2,5}
 // 🔴 You can't return an object from an arrow function with () => {
 const searchOptions = useMemo(() => {
 matchMode: 'whole-word',
 text: text
 }, [text]);
```

In JavaScript, `() => {` starts the arrow function body, so the `{` brace is not a part of your object. This is why it doesn't return an object, and leads to mistakes. You could fix it by adding parentheses like `({` and `})`:

```js {1-2,5}
 // This works, but is easy for someone to break again
 const searchOptions = useMemo(() => ({
 matchMode: 'whole-word',
 text: text
 }), [text]);
```

However, this is still confusing and too easy for someone to break by removing the parentheses.

To avoid this mistake, write a `return` statement explicitly:

```js {1-3,6-7}
 // ✅ This works and is explicit
 const searchOptions = useMemo(() => {
 return {
 matchMode: 'whole-word',
 text: text
 };
 }, [text]);
```

---

### Every time my component renders, the calculation in `useMemo` re-runs {/*every-time-my-component-renders-the-calculation-in-usememo-re-runs*/}

Make sure you've specified the dependency array as a second argument!

If you forget the dependency array, `useMemo` will re-run the calculation every time:

```js {2-3}
function TodoList({ todos, tab }) {
 // 🔴 Recalculates every time: no dependency array
 const visibleTodos = useMemo(() => filterTodos(todos, tab));
 // ...
```

This is the corrected version passing the dependency array as a second argument:

```js {2-3}
function TodoList({ todos, tab }) {
 // ✅ Does not recalculate unnecessarily
 const visibleTodos = useMemo(() => filterTodos(todos, tab), [todos, tab]);
 // ...
```

If this doesn't help, then the problem is that at least one of your dependencies is different from the previous render. You can debug this problem by manually logging your dependencies to the console:

```js
 const visibleTodos = useMemo(() => filterTodos(todos, tab), [todos, tab]);
 console.log([todos, tab]);
```

You can then right-click on the arrays from different re-renders in the console and select "Store as a global variable" for both of them. Assuming the first one got saved as `temp1` and the second one got saved as `temp2`, you can then use the browser console to check whether each dependency in both arrays is the same:

```js
Object.is(temp1[0], temp2[0]); // Is the first dependency the same between the arrays?
Object.is(temp1[1], temp2[1]); // Is the second dependency the same between the arrays?
Object.is(temp1[2], temp2[2]); // ... and so on for every dependency ...
```

When you find which dependency breaks memoization, either find a way to remove it, or [memoize it as well.](#memoizing-a-dependency-of-another-hook)

---

### I need to call `useMemo` for each list item in a loop, but it's not allowed {/*i-need-to-call-usememo-for-each-list-item-in-a-loop-but-its-not-allowed*/}

Suppose the `Chart` component is wrapped in [`memo`](/reference/react/memo). You want to skip re-rendering every `Chart` in the list when the `ReportList` component re-renders. However, you can't call `useMemo` in a loop:

```js {expectedErrors: {'react-compiler': [6]}} {5-11}
function ReportList({ items }) {
 return (
 <article>
 {items.map(item => {
 // 🔴 You can't call useMemo in a loop like this:
 const data = useMemo(() => calculateReport(item), [item]);
 return (
 <figure key={item.id}>
 <Chart data={data} />
 </figure>
 );
 })}
 </article>
 );
}
```

Instead, extract a component for each item and memoize data for individual items:

```js {5,12-18}
function ReportList({ items }) {
 return (
 <article>
 {items.map(item =>
 <Report key={item.id} item={item} />
 )}
 </article>
 );
}

function Report({ item }) {
 // ✅ Call useMemo at the top level:
 const data = useMemo(() => calculateReport(item), [item]);
 return (
 <figure>
 <Chart data={data} />
 </figure>
 );
}
```

Alternatively, you could remove `useMemo` and instead wrap `Report` itself in [`memo`.](/reference/react/memo) If the `item` prop does not change, `Report` will skip re-rendering, so `Chart` will skip re-rendering too:

```js {5,6,12}
function ReportList({ items }) {
 // ...
}

const Report = memo(function Report({ item }) {
 const data = calculateReport(item);
 return (
 <figure>
 <Chart data={data} />
 </figure>
 );
});
```

---
title: useOptimistic
---

<Intro>

`useOptimistic` is a React Hook that lets you optimistically update the UI.

```js
const [optimisticState, setOptimistic] = useOptimistic(value, reducer?);
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useOptimistic(value, reducer?)` {/*useoptimistic*/}

Call `useOptimistic` at the top level of your component to create optimistic state for a value.

```js
import { useOptimistic } from 'react';

function MyComponent({name, todos}) {
 const [optimisticAge, setOptimisticAge] = useOptimistic(28);
 const [optimisticName, setOptimisticName] = useOptimistic(name);
 const [optimisticTodos, setOptimisticTodos] = useOptimistic(todos, todoReducer);
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `value`: The value returned when there are no pending Actions.
* **optional** `reducer(currentState, action)`: The reducer function that specifies how the optimistic state gets updated. It must be pure, should take the current state and reducer action arguments, and should return the next optimistic state.

#### Returns {/*returns*/}

`useOptimistic` returns an array with exactly two values:

1. `optimisticState`: The current optimistic state. It is equal to `value` unless an Action is pending, in which case it is equal to the state returned by `reducer` (or the value passed to the set function if no `reducer` was provided).
2. The [`set` function](#setoptimistic) that lets you update the optimistic state to a different value inside an Action.

---

### `set` functions, like `setOptimistic(optimisticState)` {/*setoptimistic*/}

The `set` function returned by `useOptimistic` lets you update the state for the duration of an [Action](reference/react/useTransition#functions-called-in-starttransition-are-called-actions). You can pass the next state directly, or a function that calculates it from the previous state:

```js
const [optimisticLike, setOptimisticLike] = useOptimistic(false);
const [optimisticSubs, setOptimisticSubs] = useOptimistic(subs);

function handleClick() {
 startTransition(async () => {
 setOptimisticLike(true);
 setOptimisticSubs(a => a + 1);
 await saveChanges();
 });
}
```

#### Parameters {/*setoptimistic-parameters*/}

* `optimisticState`: The value that you want the optimistic state to be during an [Action](reference/react/useTransition#functions-called-in-starttransition-are-called-actions). If you provided a `reducer` to `useOptimistic`, this value will be passed as the second argument to your reducer. It can be a value of any type.
 * If you pass a function as `optimisticState`, it will be treated as an _updater function_. It must be pure, should take the pending state as its only argument, and should return the next optimistic state. React will put your updater function in a queue and re-render your component. During the next render, React will calculate the next state by applying the queued updaters to the previous state similar to [`useState` updaters](/reference/react/useState#setstate-parameters).

#### Returns {/*setoptimistic-returns*/}

`set` functions do not have a return value.

#### Caveats {/*setoptimistic-caveats*/}

* The `set` function must be called inside an [Action](reference/react/useTransition#functions-called-in-starttransition-are-called-actions). If you call the setter outside an Action, [React will show a warning](#an-optimistic-state-update-occurred-outside-a-transition-or-action) and the optimistic state will briefly render.

<DeepDive>

#### How optimistic state works {/*how-optimistic-state-works*/}

`useOptimistic` lets you show a temporary value while an Action is in progress:

```js
const [value, setValue] = useState('a');
const [optimistic, setOptimistic] = useOptimistic(value);

startTransition(async () => {
 setOptimistic('b');
 const newValue = await saveChanges('b');
 setValue(newValue);
});
```

When the setter is called inside an Action, `useOptimistic` will trigger a re-render to show that state while the Action is in progress. Otherwise, the `value` passed to `useOptimistic` is returned.

This state is called the "optimistic" because it is used to immediately present the user with the result of performing an Action, even though the Action actually takes time to complete.

**How the update flows**

1. **Update immediately**: When `setOptimistic('b')` is called, React immediately renders with `'b'`.

2. **(Optional) await in Action**: If you await in the Action, React continues showing `'b'`.

3. **Transition scheduled**: `setValue(newValue)` schedules an update to the real state.

4. **(Optional) wait for Suspense**: If `newValue` suspends, React continues showing `'b'`.

5. **Single render commit**: Finally, the `newValue` commits for `value` and `optimistic`.

There's no extra render to "clear" the optimistic state. The optimistic and real state converge in the same render when the Transition completes.

<Note>

#### Optimistic state is temporary {/*optimistic-state-is-temporary*/}

Optimistic state only renders while an Action is in progress, otherwise `value` is rendered.

If `saveChanges` returned `'c'`, then both `value` and `optimistic` will be `'c'`, not `'b'`.

</Note>

**How the final state is determined**

The `value` argument to `useOptimistic` determines what displays after the Action finishes. How this works depends on the pattern you use:

- **Hardcoded values** like `useOptimistic(false)`: After the Action, `state` is still `false`, so the UI shows `false`. This is useful for pending states where you always start from `false`.

- **Props or state passed in** like `useOptimistic(isLiked)`: If the parent updates `isLiked` during the Action, the new value is used after the Action completes. This is how the UI reflects the result of the Action.

- **Reducer pattern** like `useOptimistic(items, fn)`: If `items` changes while the Action is pending, React re-runs your `reducer` with the new `items` to recalculate the state. This keeps your optimistic additions on top of the latest data.

**What happens when the Action fails**

If the Action throws an error, the Transition still ends, and React renders with whatever `value` currently is. Since the parent typically only updates `value` on success, a failure means `value` hasn't changed, so the UI shows what it showed before the optimistic update. You can catch the error to show a message to the user.

</DeepDive>

---

## Usage {/*usage*/}

### Adding optimistic state to a component {/*adding-optimistic-state-to-a-component*/}

Call `useOptimistic` at the top level of your component to declare one or more optimistic states.

```js [[1, 4, "age"], [1, 5, "name"], [1, 6, "todos"], [2, 4, "optimisticAge"], [2, 5, "optimisticName"], [2, 6, "optimisticTodos"], [3, 4, "setOptimisticAge"], [3, 5, "setOptimisticName"], [3, 6, "setOptimisticTodos"], [4, 6, "reducer"]]
import { useOptimistic } from 'react';

function MyComponent({age, name, todos}) {
 const [optimisticAge, setOptimisticAge] = useOptimistic(age);
 const [optimisticName, setOptimisticName] = useOptimistic(name);
 const [optimisticTodos, setOptimisticTodos] = useOptimistic(todos, reducer);
 // ...
```

`useOptimistic` returns an array with exactly two items:

1. The <CodeStep step={2}>optimistic state</CodeStep>, initially set to the <CodeStep step={1}>value</CodeStep> provided.
2. The <CodeStep step={3}>set function</CodeStep> that lets you temporarily change the state during an [Action](reference/react/useTransition#functions-called-in-starttransition-are-called-actions).
 * If a <CodeStep step={4}>reducer</CodeStep> is provided, it will run before returning the optimistic state.

To use the <CodeStep step={2}>optimistic state</CodeStep>, call the `set` function inside an Action.

Actions are functions called inside `startTransition`:

```js {3}
function onAgeChange(e) {
 startTransition(async () => {
 setOptimisticAge(42);
 const newAge = await postAge(42);
 setAge(newAge);
 });
}
```

React will render the optimistic state `42` first while the `age` remains the current age. The Action waits for POST, and then renders the `newAge` for both `age` and `optimisticAge`.

See [How optimistic state works](#how-optimistic-state-works) for a deep dive.

<Note>

When using [Action props](/reference/react/useTransition#exposing-action-props-from-components), you can call the set function without `startTransition`:

```js [[3, 2, "setOptimisticName"]]
async function submitAction() {
 setOptimisticName('Taylor');
 await updateName('Taylor');
}
```

This works because Action props are already called inside `startTransition`.

For an example, see: [Using optimistic state in Action props](#using-optimistic-state-in-action-props).

</Note>

---

### Using optimistic state in Action props {/*using-optimistic-state-in-action-props*/}

In an [Action prop](/reference/react/useTransition#exposing-action-props-from-components), you can call the optimistic setter directly without `startTransition`.

This example sets optimistic state inside a `<form>` `submitAction` prop:

<Sandpack>

```js src/App.js
import { useState, startTransition } from 'react';
import EditName from './EditName';

export default function App() {
 const [name, setName] = useState('Alice');

 return <EditName name={name} action={setName} />;
}
```

```js src/EditName.js active
import { useOptimistic, startTransition } from 'react';
import { updateName } from './actions.js';

export default function EditName({ name, action }) {
 const [optimisticName, setOptimisticName] = useOptimistic(name);

 async function submitAction(formData) {
 const newName = formData.get('name');
 setOptimisticName(newName);

 const updatedName = await updateName(newName);
 startTransition(() => {
 action(updatedName);
 })
 }

 return (
 <form action={submitAction}>
 <p>Your name is: {optimisticName}</p>
 <p>
 <label>Change it: </label>
 <input
 type="text"
 name="name"
 disabled={name !== optimisticName}
 />
 </p>
 </form>
 );
}
```

```js src/actions.js hidden
export async function updateName(name) {
 await new Promise((res) => setTimeout(res, 1000));
 return name;
}
```

</Sandpack>

In this example, when the user submits the form, the `optimisticName` updates immediately to show the `newName` optimistically while the server request is in progress. When the request completes, `name` and `optimisticName` are rendered with the actual `updatedName` from the response.

<DeepDive>

#### Why doesn't this need `startTransition`? {/*why-doesnt-this-need-starttransition*/}

By convention, props called inside `startTransition` are named with "Action".

Since `submitAction` is named with "Action", you know it's already called inside `startTransition`.

See [Exposing `action` prop from components](/reference/react/useTransition#exposing-action-props-from-components) for the Action prop pattern.

</DeepDive>

---

### Adding optimistic state to Action props {/*adding-optimistic-state-to-action-props*/}

When creating an [Action prop](/reference/react/useTransition#exposing-action-props-from-components), you can add `useOptimistic` to show immediate feedback.

Here's a button that shows "Submitting..." while the `action` is pending:

<Sandpack>

```js src/App.js
import { useState, startTransition } from 'react';
import Button from './Button';
import { submitForm } from './actions.js';

export default function App() {
 const [count, setCount] = useState(0);
 return (
 <div>
 <Button action={async () => {
 await submitForm();
 startTransition(() => {
 setCount(c => c + 1);
 });
 }}>Increment</Button>
 {count > 0 && <p>Submitted {count}!</p>}
 </div>
 );
}
```

```js src/Button.js active
import { useOptimistic, startTransition } from 'react';

export default function Button({ action, children }) {
 const [isPending, setIsPending] = useOptimistic(false);

 return (
 <button
 disabled={isPending}
 onClick={() => {
 startTransition(async () => {
 setIsPending(true);
 await action();
 });
 }}
 >
 {isPending ? 'Submitting...' : children}
 </button>
 );
}
```

```js src/actions.js hidden
export async function submitForm() {
 await new Promise((res) => setTimeout(res, 1000));
}
```

</Sandpack>

When the button is clicked, `setIsPending(true)` uses optimistic state to immediately show "Submitting..." and disable the button. When the Action is done, `isPending` is rendered as `false` automatically.

This pattern automatically shows a pending state however `action` prop is used with `Button`:

```js
// Show pending state for a state update
<Button action={() => { setState(c => c + 1) }} />

// Show pending state for a navigation
<Button action={() => { navigate('/done') }} />

// Show pending state for a POST
<Button action={async () => { await fetch(/* ... */) }} />

// Show pending state for any combination
<Button action={async () => {
 setState(c => c + 1);
 await fetch(/* ... */);
 navigate('/done');
}} />
```

The pending state will be shown until everything in the `action` prop is finished.

<Note>

You can also use [`useTransition`](/reference/react/useTransition) to get pending state via `isPending`.

The difference is that `useTransition` gives you the `startTransition` function, while `useOptimistic` works with any Transition. Use whichever fits your component's needs.

</Note>

---

### Updating props or state optimistically {/*updating-props-or-state-optimistically*/}

You can wrap props or state in `useOptimistic` to update it immediately while an Action is in progress.

In this example, `LikeButton` receives `isLiked` as a prop and immediately toggles it when clicked:

<Sandpack>

```js src/App.js
import { useState, useOptimistic, startTransition } from 'react';
import { toggleLike } from './actions.js';

export default function App() {
 const [isLiked, setIsLiked] = useState(false);
 const [optimisticIsLiked, setOptimisticIsLiked] = useOptimistic(isLiked);

 function handleClick() {
 startTransition(async () => {
 const newValue = !optimisticIsLiked
 console.log('⏳ setting optimistic state: ' + newValue);

 setOptimisticIsLiked(newValue);
 const updatedValue = await toggleLike(newValue);

 startTransition(() => {
 console.log('⏳ setting real state: ' + updatedValue );
 setIsLiked(updatedValue);
 });
 });
 }

 if (optimisticIsLiked !== isLiked) {
 console.log('✅ rendering optimistic state: ' + optimisticIsLiked);
 } else {
 console.log('✅ rendering real value: ' + optimisticIsLiked);
 }

 return (
 <button onClick={handleClick}>
 {optimisticIsLiked ? '❤️ Unlike' : '🤍 Like'}
 </button>
 );
}
```

```js src/actions.js hidden
export async function toggleLike(value) {
 return await new Promise((res) => setTimeout(() => res(value), 1000));
 // In a real app, this would update the server
}
```

```js src/index.js hidden
import React from 'react';
import {createRoot} from 'react-dom/client';
import './styles.css';

import App from './App';

const root = createRoot(document.getElementById('root'));
// Not using StrictMode so double render logs are not shown.
root.render(<App />);
```

</Sandpack>

When the button is clicked, `setOptimisticIsLiked` immediately updates the displayed state to show the heart as liked. Meanwhile, `await toggleLike` runs in the background. When the `await` completes, `setIsLiked` parent updates the "real" `isLiked` state, and the optimistic state is rendered to match this new value.

<Note>

This example reads from `optimisticIsLiked` to calculate the next value. This works when the base state won't change, but if the base state might change while your Action is pending, you may want to use a state updater or the reducer.

See [Updating state based on the current state](#updating-state-based-on-current-state) for an example.

</Note>

---

### Updating multiple values together {/*updating-multiple-values-together*/}

When an optimistic update affects multiple related values, use a reducer to update them together. This ensures the UI stays consistent.

Here's a follow button that updates both the follow state and follower count:

<Sandpack>

```js src/App.js
import { useState, startTransition } from 'react';
import { followUser, unfollowUser } from './actions.js';
import FollowButton from './FollowButton';

export default function App() {
 const [user, setUser] = useState({
 name: 'React',
 isFollowing: false,
 followerCount: 10500
 });

 async function followAction(shouldFollow) {
 if (shouldFollow) {
 await followUser(user.name);
 } else {
 await unfollowUser(user.name);
 }
 startTransition(() => {
 setUser(current => ({
 ...current,
 isFollowing: shouldFollow,
 followerCount: current.followerCount + (shouldFollow ? 1 : -1)
 }));
 });
 }

 return <FollowButton user={user} followAction={followAction} />;
}
```

```js src/FollowButton.js active
import { useOptimistic, startTransition } from 'react';

export default function FollowButton({ user, followAction }) {
 const [optimisticState, updateOptimistic] = useOptimistic(
 { isFollowing: user.isFollowing, followerCount: user.followerCount },
 (current, isFollowing) => ({
 isFollowing,
 followerCount: current.followerCount + (isFollowing ? 1 : -1)
 })
 );

 function handleClick() {
 const newFollowState = !optimisticState.isFollowing;
 startTransition(async () => {
 updateOptimistic(newFollowState);
 await followAction(newFollowState);
 });
 }

 return (
 <div>
 <p><strong>{user.name}</strong></p>
 <p>{optimisticState.followerCount} followers</p>
 <button onClick={handleClick}>
 {optimisticState.isFollowing ? 'Unfollow' : 'Follow'}
 </button>
 </div>
 );
}
```

```js src/actions.js hidden
export async function followUser(name) {
 await new Promise((res) => setTimeout(res, 1000));
}

export async function unfollowUser(name) {
 await new Promise((res) => setTimeout(res, 1000));
}
```

</Sandpack>

The reducer receives the new `isFollowing` value and calculates both the new follow state and the updated follower count in a single update. This ensures the button text and count always stay in sync.

<DeepDive>

#### Choosing between updaters and reducers {/*choosing-between-updaters-and-reducers*/}

`useOptimistic` supports two patterns for calculating state based on current state:

**Updater functions** work like [useState updaters](/reference/react/useState#updating-state-based-on-the-previous-state). Pass a function to the setter:

```js
const [optimistic, setOptimistic] = useOptimistic(value);
setOptimistic(current => !current);
```

**Reducers** separate the update logic from the setter call:

```js
const [optimistic, dispatch] = useOptimistic(value, (current, action) => {
 // Calculate next state based on current and action
});
dispatch(action);
```

**Use updaters** for calculations where the setter call naturally describes the update. This is similar to using `setState(prev => ...)` with `useState`.

**Use reducers** when you need to pass data to the update (like which item to add) or when handling multiple types of updates with a single hook.

**Why use a reducer?**

Reducers are essential when the base state might change while your Transition is pending. If `todos` changes while your add is pending (for example, another user added a todo), React will re-run your reducer with the new `todos` to recalculate what to show. This ensures your new todo is added to the latest list, not an outdated copy.

An updater function like `setOptimistic(prev => [...prev, newItem])` would only see the state from when the Transition started, missing any updates that happened during the async work.

</DeepDive>

---

### Optimistically adding to a list {/*optimistically-adding-to-a-list*/}

When you need to optimistically add items to a list, use a `reducer`:

<Sandpack>

```js src/App.js
import { useState, startTransition } from 'react';
import { addTodo } from './actions.js';
import TodoList from './TodoList';

export default function App() {
 const [todos, setTodos] = useState([
 { id: 1, text: 'Learn React' }
 ]);

 async function addTodoAction(newTodo) {
 const savedTodo = await addTodo(newTodo);
 startTransition(() => {
 setTodos(todos => [...todos, savedTodo]);
 });
 }

 return <TodoList todos={todos} addTodoAction={addTodoAction} />;
}
```

```js src/TodoList.js active
import { useOptimistic, startTransition } from 'react';

export default function TodoList({ todos, addTodoAction }) {
 const [optimisticTodos, addOptimisticTodo] = useOptimistic(
 todos,
 (currentTodos, newTodo) => [
 ...currentTodos,
 { id: newTodo.id, text: newTodo.text, pending: true }
 ]
 );

 function handleAddTodo(text) {
 const newTodo = { id: crypto.randomUUID(), text: text };
 startTransition(async () => {
 addOptimisticTodo(newTodo);
 await addTodoAction(newTodo);
 });
 }

 return (
 <div>
 <button onClick={() => handleAddTodo('New todo')}>Add Todo</button>
 <ul>
 {optimisticTodos.map(todo => (
 <li key={todo.id}>
 {todo.text} {todo.pending && "(Adding...)"}
 </li>
 ))}
 </ul>
 </div>
 );
}
```

```js src/actions.js hidden
export async function addTodo(todo) {
 await new Promise((res) => setTimeout(res, 1000));
 // In a real app, this would save to the server
 return { ...todo, pending: false };
}
```

</Sandpack>

The `reducer` receives the current list of todos and the new todo to add. This is important because if the `todos` prop changes while your add is pending (for example, another user added a todo), React will update your optimistic state by re-running the reducer with the updated list. This ensures your new todo is added to the latest list, not an outdated copy.

<Note>

Each optimistic item includes a `pending: true` flag so you can show loading state for individual items. When the server responds and the parent updates the canonical `todos` list with the saved item, the optimistic state updates to the confirmed item without the pending flag.

</Note>

---

### Handling multiple `action` types {/*handling-multiple-action-types*/}

When you need to handle multiple types of optimistic updates (like adding and removing items), use a reducer pattern with `action` objects.

This shopping cart example shows how to handle add and remove with a single reducer:

<Sandpack>

```js src/App.js
import { useState, startTransition } from 'react';
import { addToCart, removeFromCart, updateQuantity } from './actions.js';
import ShoppingCart from './ShoppingCart';

export default function App() {
 const [cart, setCart] = useState([]);

 const cartActions = {
 async add(item) {
 await addToCart(item);
 startTransition(() => {
 setCart(current => {
 const exists = current.find(i => i.id === item.id);
 if (exists) {
 return current.map(i =>
 i.id === item.id ? { ...i, quantity: i.quantity + 1 } : i
 );
 }
 return [...current, { ...item, quantity: 1 }];
 });
 });
 },
 async remove(id) {
 await removeFromCart(id);
 startTransition(() => {
 setCart(current => current.filter(item => item.id !== id));
 });
 },
 async updateQuantity(id, quantity) {
 await updateQuantity(id, quantity);
 startTransition(() => {
 setCart(current =>
 current.map(item =>
 item.id === id ? { ...item, quantity } : item
 )
 );
 });
 }
 };

 return <ShoppingCart cart={cart} cartActions={cartActions} />;
}
```

```js src/ShoppingCart.js active
import { useOptimistic, startTransition } from 'react';

export default function ShoppingCart({ cart, cartActions }) {
 const [optimisticCart, dispatch] = useOptimistic(
 cart,
 (currentCart, action) => {
 switch (action.type) {
 case 'add':
 const exists = currentCart.find(item => item.id === action.item.id);
 if (exists) {
 return currentCart.map(item =>
 item.id === action.item.id
 ? { ...item, quantity: item.quantity + 1, pending: true }
 : item
 );
 }
 return [...currentCart, { ...action.item, quantity: 1, pending: true }];
 case 'remove':
 return currentCart.filter(item => item.id !== action.id);
 case 'update_quantity':
 return currentCart.map(item =>
 item.id === action.id
 ? { ...item, quantity: action.quantity, pending: true }
 : item
 );
 default:
 return currentCart;
 }
 }
 );

 function handleAdd(item) {
 startTransition(async () => {
 dispatch({ type: 'add', item });
 await cartActions.add(item);
 });
 }

 function handleRemove(id) {
 startTransition(async () => {
 dispatch({ type: 'remove', id });
 await cartActions.remove(id);
 });
 }

 function handleUpdateQuantity(id, quantity) {
 startTransition(async () => {
 dispatch({ type: 'update_quantity', id, quantity });
 await cartActions.updateQuantity(id, quantity);
 });
 }

 const total = optimisticCart.reduce(
 (sum, item) => sum + item.price * item.quantity,
 0
 );

 return (
 <div>
 <h2>Shopping Cart</h2>
 <div style={{ marginBottom: 16 }}>
 <button onClick={() => handleAdd({
 id: 1, name: 'T-Shirt', price: 25
 })}>
 Add T-Shirt ($25)
 </button>{' '}
 <button onClick={() => handleAdd({
 id: 2, name: 'Mug', price: 15
 })}>
 Add Mug ($15)
 </button>
 </div>
 {optimisticCart.length === 0 ? (
 <p>Your cart is empty</p>
 ) : (
 <ul>
 {optimisticCart.map(item => (
 <li key={item.id}>
 {item.name} - ${item.price} ×
 {item.quantity}
 {' '}= ${item.price * item.quantity}
 <button
 onClick={() => handleRemove(item.id)}
 style={{ marginLeft: 8 }}
 >
 Remove
 </button>
 {item.pending && ' ...'}
 </li>
 ))}
 </ul>
 )}
 <p><strong>Total: ${total}</strong></p>
 </div>
 );
}
```

```js src/actions.js hidden
export async function addToCart(item) {
 await new Promise((res) => setTimeout(res, 800));
}

export async function removeFromCart(id) {
 await new Promise((res) => setTimeout(res, 800));
}

export async function updateQuantity(id, quantity) {
 await new Promise((res) => setTimeout(res, 800));
}
```

</Sandpack>

The reducer handles three `action` types (`add`, `remove`, `update_quantity`) and returns the new optimistic state for each. Each `action` sets a `pending: true` flag so you can show visual feedback while the [Server Function](/reference/rsc/server-functions) runs.

---

### Optimistic delete with error recovery {/*optimistic-delete-with-error-recovery*/}

When deleting items optimistically, you should handle the case where the Action fails.

This example shows how to display an error message when a delete fails, and the UI automatically rolls back to show the item again.

<Sandpack>

```js src/App.js
import { useState, startTransition } from 'react';
import { deleteItem } from './actions.js';
import ItemList from './ItemList';

export default function App() {
 const [items, setItems] = useState([
 { id: 1, name: 'Learn React' },
 { id: 2, name: 'Build an app' },
 { id: 3, name: 'Deploy to production' },
 ]);

 async function deleteAction(id) {
 await deleteItem(id);
 startTransition(() => {
 setItems(current => current.filter(item => item.id !== id));
 });
 }

 return <ItemList items={items} deleteAction={deleteAction} />;
}
```

```js src/ItemList.js active
import { useState, useOptimistic, startTransition } from 'react';

export default function ItemList({ items, deleteAction }) {
 const [error, setError] = useState(null);
 const [optimisticItems, removeItem] = useOptimistic(
 items,
 (currentItems, idToRemove) =>
 currentItems.map(item =>
 item.id === idToRemove
 ? { ...item, deleting: true }
 : item
 )
 );

 function handleDelete(id) {
 setError(null);
 startTransition(async () => {
 removeItem(id);
 try {
 await deleteAction(id);
 } catch (e) {
 setError(e.message);
 }
 });
 }

 return (
 <div>
 <h2>Your Items</h2>
 <ul>
 {optimisticItems.map(item => (
 <li
 key={item.id}
 style={{
 opacity: item.deleting ? 0.5 : 1,
 textDecoration: item.deleting ? 'line-through' : 'none',
 transition: 'opacity 0.2s'
 }}
 >
 {item.name}
 <button
 onClick={() => handleDelete(item.id)}
 disabled={item.deleting}
 style={{ marginLeft: 8 }}
 >
 {item.deleting ? 'Deleting...' : 'Delete'}
 </button>
 </li>
 ))}
 </ul>
 {error && (
 <p style={{ color: 'red', padding: 8, background: '#fee' }}>
 {error}
 </p>
 )}
 </div>
 );
}
```

```js src/actions.js hidden
export async function deleteItem(id) {
 await new Promise((res) => setTimeout(res, 1000));
 // Item 3 always fails to demonstrate error recovery
 if (id === 3) {
 throw new Error('Cannot delete. Permission denied.');
 }
}
```

</Sandpack>

Try deleting 'Deploy to production'. When the delete fails, the item automatically reappears in the list.

---

## Troubleshooting {/*troubleshooting*/}

### I'm getting an error: "An optimistic state update occurred outside a Transition or Action" {/*an-optimistic-state-update-occurred-outside-a-transition-or-action*/}

You may see this error:

<ConsoleBlockMulti>

<ConsoleLogLine level="error">

An optimistic state update occurred outside a Transition or Action. To fix, move the update to an Action, or wrap with `startTransition`.

</ConsoleLogLine>

</ConsoleBlockMulti>

The optimistic setter function must be called inside `startTransition`:

```js
// 🚩 Incorrect: outside a Transition
function handleClick() {
 setOptimistic(newValue); // Warning!
 // ...
}

// ✅ Correct: inside a Transition
function handleClick() {
 startTransition(async () => {
 setOptimistic(newValue);
 // ...
 });
}

// ✅ Also correct: inside an Action prop
function submitAction(formData) {
 setOptimistic(newValue);
 // ...
}
```

When you call the setter outside an Action, the optimistic state will briefly appear and then immediately revert back to the original value. This happens because there's no Transition to "hold" the optimistic state while your Action runs.

### I'm getting an error: "Cannot update optimistic state while rendering" {/*cannot-update-optimistic-state-while-rendering*/}

You may see this error:

<ConsoleBlockMulti>

<ConsoleLogLine level="error">

Cannot update optimistic state while rendering.

</ConsoleLogLine>

</ConsoleBlockMulti>

This error occurs when you call the optimistic setter during the render phase of a component. You can only call it from event handlers, effects, or other callbacks:

```js
// 🚩 Incorrect: calling during render
function MyComponent({ items }) {
 const [isPending, setPending] = useOptimistic(false);

 // This runs during render - not allowed!
 setPending(true);

 // ...
}

// ✅ Correct: calling inside startTransition
function MyComponent({ items }) {
 const [isPending, setPending] = useOptimistic(false);

 function handleClick() {
 startTransition(() => {
 setPending(true);
 // ...
 });
 }

 // ...
}

// ✅ Also correct: calling from an Action
function MyComponent({ items }) {
 const [isPending, setPending] = useOptimistic(false);

 function action() {
 setPending(true);
 // ...
 }

 // ...
}
```

### My optimistic updates show stale values {/*my-optimistic-updates-show-stale-values*/}

If your optimistic state seems to be based on old data, consider using an updater function or reducer to calculate the optimistic state relative to the current state.

```js
// May show stale data if state changes during Action
const [optimistic, setOptimistic] = useOptimistic(count);
setOptimistic(5); // Always sets to 5, even if count changed

// Better: relative updates handle state changes correctly
const [optimistic, adjust] = useOptimistic(count, (current, delta) => current + delta);
adjust(1); // Always adds 1 to whatever the current count is
```

See [Updating state based on the current state](#updating-state-based-on-current-state) for details.

### I don't know if my optimistic update is pending {/*i-dont-know-if-my-optimistic-update-is-pending*/}

To know when `useOptimistic` is pending, you have three options:

1. **Check if `optimisticValue === value`**

```js
const [optimistic, setOptimistic] = useOptimistic(value);
const isPending = optimistic !== value;
```

If the values are not equal, there's a Transition in progress.

2. **Add a `useTransition`**

```js
const [isPending, startTransition] = useTransition();
const [optimistic, setOptimistic] = useOptimistic(value);

//...
startTransition(() => {
 setOptimistic(state);
})
```

Since `useTransition` uses `useOptimistic` for `isPending` under the hood, this is equivalent to option 1.

3. **Add a `pending` flag in your reducer**

```js
const [optimistic, addOptimistic] = useOptimistic(
 items,
 (state, newItem) => [...state, { ...newItem, isPending: true }]
);
```

Since each optimistic item has its own flag, you can show loading state for individual items.

---
title: useReducer
---

<Intro>

`useReducer` is a React Hook that lets you add a [reducer](/learn/extracting-state-logic-into-a-reducer) to your component.

```js
const [state, dispatch] = useReducer(reducer, initialArg, init?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useReducer(reducer, initialArg, init?)` {/*usereducer*/}

Call `useReducer` at the top level of your component to manage its state with a [reducer.](/learn/extracting-state-logic-into-a-reducer)

```js
import { useReducer } from 'react';

function reducer(state, action) {
 // ...
}

function MyComponent() {
 const [state, dispatch] = useReducer(reducer, { age: 42 });
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `reducer`: The reducer function that specifies how the state gets updated. It must be pure, should take the state and action as arguments, and should return the next state. State and action can be of any types.
* `initialArg`: The value from which the initial state is calculated. It can be a value of any type. How the initial state is calculated from it depends on the next `init` argument.
* **optional** `init`: The initializer function that should return the initial state. If it's not specified, the initial state is set to `initialArg`. Otherwise, the initial state is set to the result of calling `init(initialArg)`.

#### Returns {/*returns*/}

`useReducer` returns an array with exactly two values:

1. The current state. During the first render, it's set to `init(initialArg)` or `initialArg` (if there's no `init`).
2. The [`dispatch` function](#dispatch) that lets you update the state to a different value and trigger a re-render.

#### Caveats {/*caveats*/}

* `useReducer` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.
* The `dispatch` function has a stable identity, so you will often see it omitted from Effect dependencies, but including it will not cause the Effect to fire. If the linter lets you omit a dependency without errors, it is safe to do. [Learn more about removing Effect dependencies.](/learn/removing-effect-dependencies#move-dynamic-objects-and-functions-inside-your-effect)
* In Strict Mode, React will **call your reducer and initializer twice** in order to [help you find accidental impurities.](#my-reducer-or-initializer-function-runs-twice) This is development-only behavior and does not affect production. If your reducer and initializer are pure (as they should be), this should not affect your logic. The result from one of the calls is ignored.

---

### `dispatch` function {/*dispatch*/}

The `dispatch` function returned by `useReducer` lets you update the state to a different value and trigger a re-render. You need to pass the action as the only argument to the `dispatch` function:

```js
const [state, dispatch] = useReducer(reducer, { age: 42 });

function handleClick() {
 dispatch({ type: 'incremented_age' });
 // ...
```

React will set the next state to the result of calling the `reducer` function you've provided with the current `state` and the action you've passed to `dispatch`.

#### Parameters {/*dispatch-parameters*/}

* `action`: The action performed by the user. It can be a value of any type. By convention, an action is usually an object with a `type` property identifying it and, optionally, other properties with additional information.

#### Returns {/*dispatch-returns*/}

`dispatch` functions do not have a return value.

#### Caveats {/*setstate-caveats*/}

* The `dispatch` function **only updates the state variable for the *next* render**. If you read the state variable after calling the `dispatch` function, [you will still get the old value](#ive-dispatched-an-action-but-logging-gives-me-the-old-state-value) that was on the screen before your call.

* If the new value you provide is identical to the current `state`, as determined by an [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison, React will **skip re-rendering the component and its children.** This is an optimization. React may still need to call your component before ignoring the result, but it shouldn't affect your code.

* React [batches state updates.](/learn/queueing-a-series-of-state-updates) It updates the screen **after all the event handlers have run** and have called their `set` functions. This prevents multiple re-renders during a single event. In the rare case that you need to force React to update the screen earlier, for example to access the DOM, you can use [`flushSync`.](/reference/react-dom/flushSync)

---

## Usage {/*usage*/}

### Adding a reducer to a component {/*adding-a-reducer-to-a-component*/}

Call `useReducer` at the top level of your component to manage state with a [reducer.](/learn/extracting-state-logic-into-a-reducer)

```js [[1, 8, "state"], [2, 8, "dispatch"], [4, 8, "reducer"], [3, 8, "{ age: 42 }"]]
import { useReducer } from 'react';

function reducer(state, action) {
 // ...
}

function MyComponent() {
 const [state, dispatch] = useReducer(reducer, { age: 42 });
 // ...
```

`useReducer` returns an array with exactly two items:

1. The <CodeStep step={1}>current state</CodeStep> of this state variable, initially set to the <CodeStep step={3}>initial state</CodeStep> you provided.
2. The <CodeStep step={2}>`dispatch` function</CodeStep> that lets you change it in response to interaction.

To update what's on the screen, call <CodeStep step={2}>`dispatch`</CodeStep> with an object representing what the user did, called an *action*:

```js [[2, 2, "dispatch"]]
function handleClick() {
 dispatch({ type: 'incremented_age' });
}
```

React will pass the current state and the action to your <CodeStep step={4}>reducer function</CodeStep>. Your reducer will calculate and return the next state. React will store that next state, render your component with it, and update the UI.

<Sandpack>

```js
import { useReducer } from 'react';

function reducer(state, action) {
 if (action.type === 'incremented_age') {
 return {
 age: state.age + 1
 };
 }
 throw Error('Unknown action.');
}

export default function Counter() {
 const [state, dispatch] = useReducer(reducer, { age: 42 });

 return (
 <>
 <button onClick={() => {
 dispatch({ type: 'incremented_age' })
 }}>
 Increment age
 </button>
 <p>Hello! You are {state.age}.</p>
 </>
 );
}
```

```css
button { display: block; margin-top: 10px; }
```

</Sandpack>

`useReducer` is very similar to [`useState`](/reference/react/useState), but it lets you move the state update logic from event handlers into a single function outside of your component. Read more about [choosing between `useState` and `useReducer`.](/learn/extracting-state-logic-into-a-reducer#comparing-usestate-and-usereducer)

---

### Writing the reducer function {/*writing-the-reducer-function*/}

A reducer function is declared like this:

```js
function reducer(state, action) {
 // ...
}
```

Then you need to fill in the code that will calculate and return the next state. By convention, it is common to write it as a [`switch` statement.](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/switch) For each `case` in the `switch`, calculate and return some next state.

```js {4-7,10-13}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 return {
 name: state.name,
 age: state.age + 1
 };
 }
 case 'changed_name': {
 return {
 name: action.nextName,
 age: state.age
 };
 }
 }
 throw Error('Unknown action: ' + action.type);
}
```

Actions can have any shape. By convention, it's common to pass objects with a `type` property identifying the action. It should include the minimal necessary information that the reducer needs to compute the next state.

```js {5,9-12}
function Form() {
 const [state, dispatch] = useReducer(reducer, { name: 'Taylor', age: 42 });

 function handleButtonClick() {
 dispatch({ type: 'incremented_age' });
 }

 function handleInputChange(e) {
 dispatch({
 type: 'changed_name',
 nextName: e.target.value
 });
 }
 // ...
```

The action type names are local to your component. [Each action describes a single interaction, even if that leads to multiple changes in data.](/learn/extracting-state-logic-into-a-reducer#writing-reducers-well) The shape of the state is arbitrary, but usually it'll be an object or an array.

Read [extracting state logic into a reducer](/learn/extracting-state-logic-into-a-reducer) to learn more.

<Pitfall>

State is read-only. Don't modify any objects or arrays in state:

```js {4,5}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 // 🚩 Don't mutate an object in state like this:
 state.age = state.age + 1;
 return state;
 }
```

Instead, always return new objects from your reducer:

```js {4-8}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 // ✅ Instead, return a new object
 return {
 ...state,
 age: state.age + 1
 };
 }
```

Read [updating objects in state](/learn/updating-objects-in-state) and [updating arrays in state](/learn/updating-arrays-in-state) to learn more.

</Pitfall>

<Recipes titleText="Basic useReducer examples" titleId="examples-basic">

#### Form (object) {/*form-object*/}

In this example, the reducer manages a state object with two fields: `name` and `age`.

<Sandpack>

```js
import { useReducer } from 'react';

function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 return {
 name: state.name,
 age: state.age + 1
 };
 }
 case 'changed_name': {
 return {
 name: action.nextName,
 age: state.age
 };
 }
 }
 throw Error('Unknown action: ' + action.type);
}

const initialState = { name: 'Taylor', age: 42 };

export default function Form() {
 const [state, dispatch] = useReducer(reducer, initialState);

 function handleButtonClick() {
 dispatch({ type: 'incremented_age' });
 }

 function handleInputChange(e) {
 dispatch({
 type: 'changed_name',
 nextName: e.target.value
 });
 }

 return (
 <>
 <input
 value={state.name}
 onChange={handleInputChange}
 />
 <button onClick={handleButtonClick}>
 Increment age
 </button>
 <p>Hello, {state.name}. You are {state.age}.</p>
 </>
 );
}
```

```css
button { display: block; margin-top: 10px; }
```

</Sandpack>

<Solution />

#### Todo list (array) {/*todo-list-array*/}

In this example, the reducer manages an array of tasks. The array needs to be updated [without mutation.](/learn/updating-arrays-in-state)

<Sandpack>

```js src/App.js
import { useReducer } from 'react';
import AddTask from './AddTask.js';
import TaskList from './TaskList.js';

function tasksReducer(tasks, action) {
 switch (action.type) {
 case 'added': {
 return [...tasks, {
 id: action.id,
 text: action.text,
 done: false
 }];
 }
 case 'changed': {
 return tasks.map(t => {
 if (t.id === action.task.id) {
 return action.task;
 } else {
 return t;
 }
 });
 }
 case 'deleted': {
 return tasks.filter(t => t.id !== action.id);
 }
 default: {
 throw Error('Unknown action: ' + action.type);
 }
 }
}

export default function TaskApp() {
 const [tasks, dispatch] = useReducer(
 tasksReducer,
 initialTasks
 );

 function handleAddTask(text) {
 dispatch({
 type: 'added',
 id: nextId++,
 text: text,
 });
 }

 function handleChangeTask(task) {
 dispatch({
 type: 'changed',
 task: task
 });
 }

 function handleDeleteTask(taskId) {
 dispatch({
 type: 'deleted',
 id: taskId
 });
 }

 return (
 <>
 <h1>Prague itinerary</h1>
 <AddTask
 onAddTask={handleAddTask}
 />
 <TaskList
 tasks={tasks}
 onChangeTask={handleChangeTask}
 onDeleteTask={handleDeleteTask}
 />
 </>
 );
}

let nextId = 3;
const initialTasks = [
 { id: 0, text: 'Visit Kafka Museum', done: true },
 { id: 1, text: 'Watch a puppet show', done: false },
 { id: 2, text: 'Lennon Wall pic', done: false }
];
```

```js src/AddTask.js hidden
import { useState } from 'react';

export default function AddTask({ onAddTask }) {
 const [text, setText] = useState('');
 return (
 <>
 <input
 placeholder="Add task"
 value={text}
 onChange={e => setText(e.target.value)}
 />
 <button onClick={() => {
 setText('');
 onAddTask(text);
 }}>Add</button>
 </>
 )
}
```

```js src/TaskList.js hidden
import { useState } from 'react';

export default function TaskList({
 tasks,
 onChangeTask,
 onDeleteTask
}) {
 return (
 <ul>
 {tasks.map(task => (
 <li key={task.id}>
 <Task
 task={task}
 onChange={onChangeTask}
 onDelete={onDeleteTask}
 />
 </li>
 ))}
 </ul>
 );
}

function Task({ task, onChange, onDelete }) {
 const [isEditing, setIsEditing] = useState(false);
 let taskContent;
 if (isEditing) {
 taskContent = (
 <>
 <input
 value={task.text}
 onChange={e => {
 onChange({
 ...task,
 text: e.target.value
 });
 }} />
 <button onClick={() => setIsEditing(false)}>
 Save
 </button>
 </>
 );
 } else {
 taskContent = (
 <>
 {task.text}
 <button onClick={() => setIsEditing(true)}>
 Edit
 </button>
 </>
 );
 }
 return (
 <label>
 <input
 type="checkbox"
 checked={task.done}
 onChange={e => {
 onChange({
 ...task,
 done: e.target.checked
 });
 }}
 />
 {taskContent}
 <button onClick={() => onDelete(task.id)}>
 Delete
 </button>
 </label>
 );
}
```

```css
button { margin: 5px; }
li { list-style-type: none; }
ul, li { margin: 0; padding: 0; }
```

</Sandpack>

<Solution />

#### Writing concise update logic with Immer {/*writing-concise-update-logic-with-immer*/}

If updating arrays and objects without mutation feels tedious, you can use a library like [Immer](https://github.com/immerjs/use-immer#useimmerreducer) to reduce repetitive code. Immer lets you write concise code as if you were mutating objects, but under the hood it performs immutable updates:

<Sandpack>

```js src/App.js
import { useImmerReducer } from 'use-immer';
import AddTask from './AddTask.js';
import TaskList from './TaskList.js';

function tasksReducer(draft, action) {
 switch (action.type) {
 case 'added': {
 draft.push({
 id: action.id,
 text: action.text,
 done: false
 });
 break;
 }
 case 'changed': {
 const index = draft.findIndex(t =>
 t.id === action.task.id
 );
 draft[index] = action.task;
 break;
 }
 case 'deleted': {
 return draft.filter(t => t.id !== action.id);
 }
 default: {
 throw Error('Unknown action: ' + action.type);
 }
 }
}

export default function TaskApp() {
 const [tasks, dispatch] = useImmerReducer(
 tasksReducer,
 initialTasks
 );

 function handleAddTask(text) {
 dispatch({
 type: 'added',
 id: nextId++,
 text: text,
 });
 }

 function handleChangeTask(task) {
 dispatch({
 type: 'changed',
 task: task
 });
 }

 function handleDeleteTask(taskId) {
 dispatch({
 type: 'deleted',
 id: taskId
 });
 }

 return (
 <>
 <h1>Prague itinerary</h1>
 <AddTask
 onAddTask={handleAddTask}
 />
 <TaskList
 tasks={tasks}
 onChangeTask={handleChangeTask}
 onDeleteTask={handleDeleteTask}
 />
 </>
 );
}

let nextId = 3;
const initialTasks = [
 { id: 0, text: 'Visit Kafka Museum', done: true },
 { id: 1, text: 'Watch a puppet show', done: false },
 { id: 2, text: 'Lennon Wall pic', done: false },
];
```

```js src/AddTask.js hidden
import { useState } from 'react';

export default function AddTask({ onAddTask }) {
 const [text, setText] = useState('');
 return (
 <>
 <input
 placeholder="Add task"
 value={text}
 onChange={e => setText(e.target.value)}
 />
 <button onClick={() => {
 setText('');
 onAddTask(text);
 }}>Add</button>
 </>
 )
}
```

```js src/TaskList.js hidden
import { useState } from 'react';

export default function TaskList({
 tasks,
 onChangeTask,
 onDeleteTask
}) {
 return (
 <ul>
 {tasks.map(task => (
 <li key={task.id}>
 <Task
 task={task}
 onChange={onChangeTask}
 onDelete={onDeleteTask}
 />
 </li>
 ))}
 </ul>
 );
}

function Task({ task, onChange, onDelete }) {
 const [isEditing, setIsEditing] = useState(false);
 let taskContent;
 if (isEditing) {
 taskContent = (
 <>
 <input
 value={task.text}
 onChange={e => {
 onChange({
 ...task,
 text: e.target.value
 });
 }} />
 <button onClick={() => setIsEditing(false)}>
 Save
 </button>
 </>
 );
 } else {
 taskContent = (
 <>
 {task.text}
 <button onClick={() => setIsEditing(true)}>
 Edit
 </button>
 </>
 );
 }
 return (
 <label>
 <input
 type="checkbox"
 checked={task.done}
 onChange={e => {
 onChange({
 ...task,
 done: e.target.checked
 });
 }}
 />
 {taskContent}
 <button onClick={() => onDelete(task.id)}>
 Delete
 </button>
 </label>
 );
}
```

```css
button { margin: 5px; }
li { list-style-type: none; }
ul, li { margin: 0; padding: 0; }
```

```json package.json
{
 "dependencies": {
 "immer": "1.7.3",
 "react": "latest",
 "react-dom": "latest",
 "react-scripts": "latest",
 "use-immer": "0.5.1"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Avoiding recreating the initial state {/*avoiding-recreating-the-initial-state*/}

React saves the initial state once and ignores it on the next renders.

```js
function createInitialState(username) {
 // ...
}

function TodoList({ username }) {
 const [state, dispatch] = useReducer(reducer, createInitialState(username));
 // ...
```

Although the result of `createInitialState(username)` is only used for the initial render, you're still calling this function on every render. This can be wasteful if it's creating large arrays or performing expensive calculations.

To solve this, you may **pass it as an _initializer_ function** to `useReducer` as the third argument instead:

```js {6}
function createInitialState(username) {
 // ...
}

function TodoList({ username }) {
 const [state, dispatch] = useReducer(reducer, username, createInitialState);
 // ...
```

Notice that you’re passing `createInitialState`, which is the *function itself*, and not `createInitialState()`, which is the result of calling it. This way, the initial state does not get re-created after initialization.

In the above example, `createInitialState` takes a `username` argument. If your initializer doesn't need any information to compute the initial state, you may pass `null` as the second argument to `useReducer`.

<Recipes titleText="The difference between passing an initializer and passing the initial state directly" titleId="examples-initializer">

#### Passing the initializer function {/*passing-the-initializer-function*/}

This example passes the initializer function, so the `createInitialState` function only runs during initialization. It does not run when component re-renders, such as when you type into the input.

<Sandpack>

```js src/App.js hidden
import TodoList from './TodoList.js';

export default function App() {
 return <TodoList username="Taylor" />;
}
```

```js src/TodoList.js active
import { useReducer } from 'react';

function createInitialState(username) {
 const initialTodos = [];
 for (let i = 0; i < 50; i++) {
 initialTodos.push({
 id: i,
 text: username + "'s task #" + (i + 1)
 });
 }
 return {
 draft: '',
 todos: initialTodos,
 };
}

function reducer(state, action) {
 switch (action.type) {
 case 'changed_draft': {
 return {
 draft: action.nextDraft,
 todos: state.todos,
 };
 };
 case 'added_todo': {
 return {
 draft: '',
 todos: [{
 id: state.todos.length,
 text: state.draft
 }, ...state.todos]
 }
 }
 }
 throw Error('Unknown action: ' + action.type);
}

export default function TodoList({ username }) {
 const [state, dispatch] = useReducer(
 reducer,
 username,
 createInitialState
 );
 return (
 <>
 <input
 value={state.draft}
 onChange={e => {
 dispatch({
 type: 'changed_draft',
 nextDraft: e.target.value
 })
 }}
 />
 <button onClick={() => {
 dispatch({ type: 'added_todo' });
 }}>Add</button>
 <ul>
 {state.todos.map(item => (
 <li key={item.id}>
 {item.text}
 </li>
 ))}
 </ul>
 </>
 );
}
```

</Sandpack>

<Solution />

#### Passing the initial state directly {/*passing-the-initial-state-directly*/}

This example **does not** pass the initializer function, so the `createInitialState` function runs on every render, such as when you type into the input. There is no observable difference in behavior, but this code is less efficient.

<Sandpack>

```js src/App.js hidden
import TodoList from './TodoList.js';

export default function App() {
 return <TodoList username="Taylor" />;
}
```

```js src/TodoList.js active
import { useReducer } from 'react';

function createInitialState(username) {
 const initialTodos = [];
 for (let i = 0; i < 50; i++) {
 initialTodos.push({
 id: i,
 text: username + "'s task #" + (i + 1)
 });
 }
 return {
 draft: '',
 todos: initialTodos,
 };
}

function reducer(state, action) {
 switch (action.type) {
 case 'changed_draft': {
 return {
 draft: action.nextDraft,
 todos: state.todos,
 };
 };
 case 'added_todo': {
 return {
 draft: '',
 todos: [{
 id: state.todos.length,
 text: state.draft
 }, ...state.todos]
 }
 }
 }
 throw Error('Unknown action: ' + action.type);
}

export default function TodoList({ username }) {
 const [state, dispatch] = useReducer(
 reducer,
 createInitialState(username)
 );
 return (
 <>
 <input
 value={state.draft}
 onChange={e => {
 dispatch({
 type: 'changed_draft',
 nextDraft: e.target.value
 })
 }}
 />
 <button onClick={() => {
 dispatch({ type: 'added_todo' });
 }}>Add</button>
 <ul>
 {state.todos.map(item => (
 <li key={item.id}>
 {item.text}
 </li>
 ))}
 </ul>
 </>
 );
}
```

</Sandpack>

<Solution />

</Recipes>

---

## Troubleshooting {/*troubleshooting*/}

### I've dispatched an action, but logging gives me the old state value {/*ive-dispatched-an-action-but-logging-gives-me-the-old-state-value*/}

Calling the `dispatch` function **does not change state in the running code**:

```js {4,5,8}
function handleClick() {
 console.log(state.age); // 42

 dispatch({ type: 'incremented_age' }); // Request a re-render with 43
 console.log(state.age); // Still 42!

 setTimeout(() => {
 console.log(state.age); // Also 42!
 }, 5000);
}
```

This is because [states behaves like a snapshot.](/learn/state-as-a-snapshot) Updating state requests another render with the new state value, but does not affect the `state` JavaScript variable in your already-running event handler.

If you need to guess the next state value, you can calculate it manually by calling the reducer yourself:

```js
const action = { type: 'incremented_age' };
dispatch(action);

const nextState = reducer(state, action);
console.log(state); // { age: 42 }
console.log(nextState); // { age: 43 }
```

---

### I've dispatched an action, but the screen doesn't update {/*ive-dispatched-an-action-but-the-screen-doesnt-update*/}

React will **ignore your update if the next state is equal to the previous state,** as determined by an [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. This usually happens when you change an object or an array in state directly:

```js {4-5,9-10}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 // 🚩 Wrong: mutating existing object
 state.age++;
 return state;
 }
 case 'changed_name': {
 // 🚩 Wrong: mutating existing object
 state.name = action.nextName;
 return state;
 }
 // ...
 }
}
```

You mutated an existing `state` object and returned it, so React ignored the update. To fix this, you need to ensure that you're always [updating objects in state](/learn/updating-objects-in-state) and [updating arrays in state](/learn/updating-arrays-in-state) instead of mutating them:

```js {4-8,11-15}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 // ✅ Correct: creating a new object
 return {
 ...state,
 age: state.age + 1
 };
 }
 case 'changed_name': {
 // ✅ Correct: creating a new object
 return {
 ...state,
 name: action.nextName
 };
 }
 // ...
 }
}
```

---

### A part of my reducer state becomes undefined after dispatching {/*a-part-of-my-reducer-state-becomes-undefined-after-dispatching*/}

Make sure that every `case` branch **copies all of the existing fields** when returning the new state:

```js {5}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 return {
 ...state, // Don't forget this!
 age: state.age + 1
 };
 }
 // ...
```

Without `...state` above, the returned next state would only contain the `age` field and nothing else.

---

### My entire reducer state becomes undefined after dispatching {/*my-entire-reducer-state-becomes-undefined-after-dispatching*/}

If your state unexpectedly becomes `undefined`, you're likely forgetting to `return` state in one of the cases, or your action type doesn't match any of the `case` statements. To find why, throw an error outside the `switch`:

```js {10}
function reducer(state, action) {
 switch (action.type) {
 case 'incremented_age': {
 // ...
 }
 case 'edited_name': {
 // ...
 }
 }
 throw Error('Unknown action: ' + action.type);
}
```

You can also use a static type checker like TypeScript to catch such mistakes.

---

### I'm getting an error: "Too many re-renders" {/*im-getting-an-error-too-many-re-renders*/}

You might get an error that says: `Too many re-renders. React limits the number of renders to prevent an infinite loop.` Typically, this means that you're unconditionally dispatching an action *during render*, so your component enters a loop: render, dispatch (which causes a render), render, dispatch (which causes a render), and so on. Very often, this is caused by a mistake in specifying an event handler:

```js {1-2}
// 🚩 Wrong: calls the handler during render
return <button onClick={handleClick()}>Click me</button>

// ✅ Correct: passes down the event handler
return <button onClick={handleClick}>Click me</button>

// ✅ Correct: passes down an inline function
return <button onClick={(e) => handleClick(e)}>Click me</button>
```

If you can't find the cause of this error, click on the arrow next to the error in the console and look through the JavaScript stack to find the specific `dispatch` function call responsible for the error.

---

### My reducer or initializer function runs twice {/*my-reducer-or-initializer-function-runs-twice*/}

In [Strict Mode](/reference/react/StrictMode), React will call your reducer and initializer functions twice. This shouldn't break your code.

This **development-only** behavior helps you [keep components pure.](/learn/keeping-components-pure) React uses the result of one of the calls, and ignores the result of the other call. As long as your component, initializer, and reducer functions are pure, this shouldn't affect your logic. However, if they are accidentally impure, this helps you notice the mistakes.

For example, this impure reducer function mutates an array in state:

```js {4-6}
function reducer(state, action) {
 switch (action.type) {
 case 'added_todo': {
 // 🚩 Mistake: mutating state
 state.todos.push({ id: nextId++, text: action.text });
 return state;
 }
 // ...
 }
}
```

Because React calls your reducer function twice, you'll see the todo was added twice, so you'll know that there is a mistake. In this example, you can fix the mistake by [replacing the array instead of mutating it](/learn/updating-arrays-in-state#adding-to-an-array):

```js {4-11}
function reducer(state, action) {
 switch (action.type) {
 case 'added_todo': {
 // ✅ Correct: replacing with new state
 return {
 ...state,
 todos: [
 ...state.todos,
 { id: nextId++, text: action.text }
 ]
 };
 }
 // ...
 }
}
```

Now that this reducer function is pure, calling it an extra time doesn't make a difference in behavior. This is why React calling it twice helps you find mistakes. **Only component, initializer, and reducer functions need to be pure.** Event handlers don't need to be pure, so React will never call your event handlers twice.

Read [keeping components pure](/learn/keeping-components-pure) to learn more.

---
title: useRef
---

<Intro>

`useRef` is a React Hook that lets you reference a value that's not needed for rendering.

```js
const ref = useRef(initialValue)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useRef(initialValue)` {/*useref*/}

Call `useRef` at the top level of your component to declare a [ref.](/learn/referencing-values-with-refs)

```js
import { useRef } from 'react';

function MyComponent() {
 const intervalRef = useRef(0);
 const inputRef = useRef(null);
 // ...
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `initialValue`: The value you want the ref object's `current` property to be initially. It can be a value of any type. This argument is ignored after the initial render.

#### Returns {/*returns*/}

`useRef` returns an object with a single property:

* `current`: Initially, it's set to the `initialValue` you have passed. You can later set it to something else. If you pass the ref object to React as a `ref` attribute to a JSX node, React will set its `current` property.

On the next renders, `useRef` will return the same object.

#### Caveats {/*caveats*/}

* You can mutate the `ref.current` property. Unlike state, it is mutable. However, if it holds an object that is used for rendering (for example, a piece of your state), then you shouldn't mutate that object.
* When you change the `ref.current` property, React does not re-render your component. React is not aware of when you change it because a ref is a plain JavaScript object.
* Do not write _or read_ `ref.current` during rendering, except for [initialization.](#avoiding-recreating-the-ref-contents) This makes your component's behavior unpredictable.
* In Strict Mode, React will **call your component function twice** in order to [help you find accidental impurities.](/reference/react/useState#my-initializer-or-updater-function-runs-twice) This is development-only behavior and does not affect production. Each ref object will be created twice, but one of the versions will be discarded. If your component function is pure (as it should be), this should not affect the behavior.

---

## Usage {/*usage*/}

### Referencing a value with a ref {/*referencing-a-value-with-a-ref*/}

Call `useRef` at the top level of your component to declare one or more [refs.](/learn/referencing-values-with-refs)

```js [[1, 4, "intervalRef"], [3, 4, "0"]]
import { useRef } from 'react';

function Stopwatch() {
 const intervalRef = useRef(0);
 // ...
```

`useRef` returns a <CodeStep step={1}>ref object</CodeStep> with a single <CodeStep step={2}>`current` property</CodeStep> initially set to the <CodeStep step={3}>initial value</CodeStep> you provided.

On the next renders, `useRef` will return the same object. You can change its `current` property to store information and read it later. This might remind you of [state](/reference/react/useState), but there is an important difference.

**Changing a ref does not trigger a re-render.** This means refs are perfect for storing information that doesn't affect the visual output of your component. For example, if you need to store an [interval ID](https://developer.mozilla.org/en-US/docs/Web/API/setInterval) and retrieve it later, you can put it in a ref. To update the value inside the ref, you need to manually change its <CodeStep step={2}>`current` property</CodeStep>:

```js [[2, 5, "intervalRef.current"]]
function handleStartClick() {
 const intervalId = setInterval(() => {
 // ...
 }, 1000);
 intervalRef.current = intervalId;
}
```

Later, you can read that interval ID from the ref so that you can call [clear that interval](https://developer.mozilla.org/en-US/docs/Web/API/clearInterval):

```js [[2, 2, "intervalRef.current"]]
function handleStopClick() {
 const intervalId = intervalRef.current;
 clearInterval(intervalId);
}
```

By using a ref, you ensure that:

- You can **store information** between re-renders (unlike regular variables, which reset on every render).
- Changing it **does not trigger a re-render** (unlike state variables, which trigger a re-render).
- The **information is local** to each copy of your component (unlike the variables outside, which are shared).

Changing a ref does not trigger a re-render, so refs are not appropriate for storing information you want to display on the screen. Use state for that instead. Read more about [choosing between `useRef` and `useState`.](/learn/referencing-values-with-refs#differences-between-refs-and-state)

<Recipes titleText="Examples of referencing a value with useRef" titleId="examples-value">

#### Click counter {/*click-counter*/}

This component uses a ref to keep track of how many times the button was clicked. Note that it's okay to use a ref instead of state here because the click count is only read and written in an event handler.

<Sandpack>

```js
import { useRef } from 'react';

export default function Counter() {
 let ref = useRef(0);

 function handleClick() {
 ref.current = ref.current + 1;
 alert('You clicked ' + ref.current + ' times!');
 }

 return (
 <button onClick={handleClick}>
 Click me!
 </button>
 );
}
```

</Sandpack>

If you show `{ref.current}` in the JSX, the number won't update on click. This is because setting `ref.current` does not trigger a re-render. Information that's used for rendering should be state instead.

<Solution />

#### A stopwatch {/*a-stopwatch*/}

This example uses a combination of state and refs. Both `startTime` and `now` are state variables because they are used for rendering. But we also need to hold an [interval ID](https://developer.mozilla.org/en-US/docs/Web/API/setInterval) so that we can stop the interval on button press. Since the interval ID is not used for rendering, it's appropriate to keep it in a ref, and manually update it.

<Sandpack>

```js
import { useState, useRef } from 'react';

export default function Stopwatch() {
 const [startTime, setStartTime] = useState(null);
 const [now, setNow] = useState(null);
 const intervalRef = useRef(null);

 function handleStart() {
 setStartTime(Date.now());
 setNow(Date.now());

 clearInterval(intervalRef.current);
 intervalRef.current = setInterval(() => {
 setNow(Date.now());
 }, 10);
 }

 function handleStop() {
 clearInterval(intervalRef.current);
 }

 let secondsPassed = 0;
 if (startTime != null && now != null) {
 secondsPassed = (now - startTime) / 1000;
 }

 return (
 <>
 <h1>Time passed: {secondsPassed.toFixed(3)}</h1>
 <button onClick={handleStart}>
 Start
 </button>
 <button onClick={handleStop}>
 Stop
 </button>
 </>
 );
}
```

</Sandpack>

<Solution />

</Recipes>

<Pitfall>

**Do not write _or read_ `ref.current` during rendering.**

React expects that the body of your component [behaves like a pure function](/learn/keeping-components-pure):

- If the inputs ([props](/learn/passing-props-to-a-component), [state](/learn/state-a-components-memory), and [context](/learn/passing-data-deeply-with-context)) are the same, it should return exactly the same JSX.
- Calling it in a different order or with different arguments should not affect the results of other calls.

Reading or writing a ref **during rendering** breaks these expectations.

```js {expectedErrors: {'react-compiler': [4]}} {3-4,6-7}
function MyComponent() {
 // ...
 // 🚩 Don't write a ref during rendering
 myRef.current = 123;
 // ...
 // 🚩 Don't read a ref during rendering
 return <h1>{myOtherRef.current}</h1>;
}
```

You can read or write refs **from event handlers or effects instead**.

```js {4-5,9-10}
function MyComponent() {
 // ...
 useEffect(() => {
 // ✅ You can read or write refs in effects
 myRef.current = 123;
 });
 // ...
 function handleClick() {
 // ✅ You can read or write refs in event handlers
 doSomething(myOtherRef.current);
 }
 // ...
}
```

If you *have to* read [or write](/reference/react/useState#storing-information-from-previous-renders) something during rendering, [use state](/reference/react/useState) instead.

When you break these rules, your component might still work, but most of the newer features we're adding to React will rely on these expectations. Read more about [keeping your components pure.](/learn/keeping-components-pure#where-you-_can_-cause-side-effects)

</Pitfall>

---

### Manipulating the DOM with a ref {/*manipulating-the-dom-with-a-ref*/}

It's particularly common to use a ref to manipulate the [DOM.](https://developer.mozilla.org/en-US/docs/Web/API/HTML_DOM_API) React has built-in support for this.

First, declare a <CodeStep step={1}>ref object</CodeStep> with an <CodeStep step={3}>initial value</CodeStep> of `null`:

```js [[1, 4, "inputRef"], [3, 4, "null"]]
import { useRef } from 'react';

function MyComponent() {
 const inputRef = useRef(null);
 // ...
```

Then pass your ref object as the `ref` attribute to the JSX of the DOM node you want to manipulate:

```js [[1, 2, "inputRef"]]
 // ...
 return <input ref={inputRef} />;
```

After React creates the DOM node and puts it on the screen, React will set the <CodeStep step={2}>`current` property</CodeStep> of your ref object to that DOM node. Now you can access the `<input>`'s DOM node and call methods like [`focus()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLElement/focus):

```js [[2, 2, "inputRef.current"]]
 function handleClick() {
 inputRef.current.focus();
 }
```

React will set the `current` property back to `null` when the node is removed from the screen.

Read more about [manipulating the DOM with refs.](/learn/manipulating-the-dom-with-refs)

<Recipes titleText="Examples of manipulating the DOM with useRef" titleId="examples-dom">

#### Focusing a text input {/*focusing-a-text-input*/}

In this example, clicking the button will focus the input:

<Sandpack>

```js
import { useRef } from 'react';

export default function Form() {
 const inputRef = useRef(null);

 function handleClick() {
 inputRef.current.focus();
 }

 return (
 <>
 <input ref={inputRef} />
 <button onClick={handleClick}>
 Focus the input
 </button>
 </>
 );
}
```

</Sandpack>

<Solution />

#### Scrolling an image into view {/*scrolling-an-image-into-view*/}

In this example, clicking the button will scroll an image into view. It uses a ref to the list DOM node, and then calls DOM [`querySelectorAll`](https://developer.mozilla.org/en-US/docs/Web/API/Document/querySelectorAll) API to find the image we want to scroll to.

<Sandpack>

```js
import { useRef } from 'react';

export default function CatFriends() {
 const listRef = useRef(null);

 function scrollToIndex(index) {
 const listNode = listRef.current;
 // This line assumes a particular DOM structure:
 const imgNode = listNode.querySelectorAll('li > img')[index];
 imgNode.scrollIntoView({
 behavior: 'smooth',
 block: 'nearest',
 inline: 'center'
 });
 }

 return (
 <>
 <nav>
 <button onClick={() => scrollToIndex(0)}>
 Neo
 </button>
 <button onClick={() => scrollToIndex(1)}>
 Millie
 </button>
 <button onClick={() => scrollToIndex(2)}>
 Bella
 </button>
 </nav>
 <div>
 <ul ref={listRef}>
 <li>
 <img
 src="https://placecats.com/neo/300/200"
 alt="Neo"
 />
 </li>
 <li>
 <img
 src="https://placecats.com/millie/200/200"
 alt="Millie"
 />
 </li>
 <li>
 <img
 src="https://placecats.com/bella/199/200"
 alt="Bella"
 />
 </li>
 </ul>
 </div>
 </>
 );
}
```

```css
div {
 width: 100%;
 overflow: hidden;
}

nav {
 text-align: center;
}

button {
 margin: .25rem;
}

ul,
li {
 list-style: none;
 white-space: nowrap;
}

li {
 display: inline;
 padding: 0.5rem;
}
```

</Sandpack>

<Solution />

#### Playing and pausing a video {/*playing-and-pausing-a-video*/}

This example uses a ref to call [`play()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/play) and [`pause()`](https://developer.mozilla.org/en-US/docs/Web/API/HTMLMediaElement/pause) on a `<video>` DOM node.

<Sandpack>

```js
import { useState, useRef } from 'react';

export default function VideoPlayer() {
 const [isPlaying, setIsPlaying] = useState(false);
 const ref = useRef(null);

 function handleClick() {
 const nextIsPlaying = !isPlaying;
 setIsPlaying(nextIsPlaying);

 if (nextIsPlaying) {
 ref.current.play();
 } else {
 ref.current.pause();
 }
 }

 return (
 <>
 <button onClick={handleClick}>
 {isPlaying ? 'Pause' : 'Play'}
 </button>
 <video
 width="250"
 ref={ref}
 onPlay={() => setIsPlaying(true)}
 onPause={() => setIsPlaying(false)}
 >
 <source
 src="https://interactive-examples.mdn.mozilla.net/media/cc0-videos/flower.mp4"
 type="video/mp4"
 />
 </video>
 </>
 );
}
```

```css
button { display: block; margin-bottom: 20px; }
```

</Sandpack>

<Solution />

#### Exposing a ref to your own component {/*exposing-a-ref-to-your-own-component*/}

Sometimes, you may want to let the parent component manipulate the DOM inside of your component. For example, maybe you're writing a `MyInput` component, but you want the parent to be able to focus the input (which the parent has no access to). You can create a `ref` in the parent and pass the `ref` as prop to the child component. Read a [detailed walkthrough](/learn/manipulating-the-dom-with-refs#accessing-another-components-dom-nodes) here.

<Sandpack>

```js
import { useRef } from 'react';

function MyInput({ ref }) {
 return <input ref={ref} />;
};

export default function Form() {
 const inputRef = useRef(null);

 function handleClick() {
 inputRef.current.focus();
 }

 return (
 <>
 <MyInput ref={inputRef} />
 <button onClick={handleClick}>
 Focus the input
 </button>
 </>
 );
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Avoiding recreating the ref contents {/*avoiding-recreating-the-ref-contents*/}

React saves the initial ref value once and ignores it on the next renders.

```js
function Video() {
 const playerRef = useRef(new VideoPlayer());
 // ...
```

Although the result of `new VideoPlayer()` is only used for the initial render, you're still calling this function on every render. This can be wasteful if it's creating expensive objects.

To solve it, you may initialize the ref like this instead:

```js
function Video() {
 const playerRef = useRef(null);
 if (playerRef.current === null) {
 playerRef.current = new VideoPlayer();
 }
 // ...
```

Normally, writing or reading `ref.current` during render is not allowed. However, it's fine in this case because the result is always the same, and the condition only executes during initialization so it's fully predictable.

<DeepDive>

#### How to avoid null checks when initializing useRef later {/*how-to-avoid-null-checks-when-initializing-use-ref-later*/}

If you use a type checker and don't want to always check for `null`, you can try a pattern like this instead:

```js
function Video() {
 const playerRef = useRef(null);

 function getPlayer() {
 if (playerRef.current !== null) {
 return playerRef.current;
 }
 const player = new VideoPlayer();
 playerRef.current = player;
 return player;
 }

 // ...
```

Here, the `playerRef` itself is nullable. However, you should be able to convince your type checker that there is no case in which `getPlayer()` returns `null`. Then use `getPlayer()` in your event handlers.

</DeepDive>

---

## Troubleshooting {/*troubleshooting*/}

### I can't get a ref to a custom component {/*i-cant-get-a-ref-to-a-custom-component*/}

If you try to pass a `ref` to your own component like this:

```js
const inputRef = useRef(null);

return <MyInput ref={inputRef} />;
```

You might get an error in the console:

<ConsoleBlock level="error">

TypeError: Cannot read properties of null

</ConsoleBlock>

By default, your own components don't expose refs to the DOM nodes inside them.

To fix this, find the component that you want to get a ref to:

```js
export default function MyInput({ value, onChange }) {
 return (
 <input
 value={value}
 onChange={onChange}
 />
 );
}
```

And then add `ref` to the list of props your component accepts and pass `ref` as a prop to the relevant child [built-in component](/reference/react-dom/components/common) like this:

```js {1,6}
function MyInput({ value, onChange, ref }) {
 return (
 <input
 value={value}
 onChange={onChange}
 ref={ref}
 />
 );
};

export default MyInput;
```

Then the parent component can get a ref to it.

Read more about [accessing another component's DOM nodes.](/learn/manipulating-the-dom-with-refs#accessing-another-components-dom-nodes)

---
title: useState
---

<Intro>

`useState` is a React Hook that lets you add a [state variable](/learn/state-a-components-memory) to your component.

```js
const [state, setState] = useState(initialState)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useState(initialState)` {/*usestate*/}

Call `useState` at the top level of your component to declare a [state variable.](/learn/state-a-components-memory)

```js
import { useState } from 'react';

function MyComponent() {
 const [age, setAge] = useState(28);
 const [name, setName] = useState('Taylor');
 const [todos, setTodos] = useState(() => createTodos());
 // ...
```

The convention is to name state variables like `[something, setSomething]` using [array destructuring.](https://javascript.info/destructuring-assignment)

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `initialState`: The value you want the state to be initially. It can be a value of any type, but there is a special behavior for functions. This argument is ignored after the initial render.
 * If you pass a function as `initialState`, it will be treated as an _initializer function_. It should be pure, should take no arguments, and should return a value of any type. React will call your initializer function when initializing the component, and store its return value as the initial state. [See an example below.](#avoiding-recreating-the-initial-state)

#### Returns {/*returns*/}

`useState` returns an array with exactly two values:

1. The current state. During the first render, it will match the `initialState` you have passed.
2. The [`set` function](#setstate) that lets you update the state to a different value and trigger a re-render.

#### Caveats {/*caveats*/}

* `useState` is a Hook, so you can only call it **at the top level of your component** or your own Hooks. You can't call it inside loops or conditions. If you need that, extract a new component and move the state into it.
* In Strict Mode, React will **call your initializer function twice** in order to [help you find accidental impurities.](#my-initializer-or-updater-function-runs-twice) This is development-only behavior and does not affect production. If your initializer function is pure (as it should be), this should not affect the behavior. The result from one of the calls will be ignored.

---

### `set` functions, like `setSomething(nextState)` {/*setstate*/}

The `set` function returned by `useState` lets you update the state to a different value and trigger a re-render. You can pass the next state directly, or a function that calculates it from the previous state:

```js
const [name, setName] = useState('Edward');

function handleClick() {
 setName('Taylor');
 setAge(a => a + 1);
 // ...
```

#### Parameters {/*setstate-parameters*/}

* `nextState`: The value that you want the state to be. It can be a value of any type, but there is a special behavior for functions.
 * If you pass a function as `nextState`, it will be treated as an _updater function_. It must be pure, should take the pending state as its only argument, and should return the next state. React will put your updater function in a queue and re-render your component. During the next render, React will calculate the next state by applying all of the queued updaters to the previous state. [See an example below.](#updating-state-based-on-the-previous-state)

#### Returns {/*setstate-returns*/}

`set` functions do not have a return value.

#### Caveats {/*setstate-caveats*/}

* The `set` function **only updates the state variable for the *next* render**. If you read the state variable after calling the `set` function, [you will still get the old value](#ive-updated-the-state-but-logging-gives-me-the-old-value) that was on the screen before your call.

* If the new value you provide is identical to the current `state`, as determined by an [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison, React will **skip re-rendering the component and its children.** This is an optimization. Although in some cases React may still need to call your component before skipping the children, it shouldn't affect your code.

* React [batches state updates.](/learn/queueing-a-series-of-state-updates) It updates the screen **after all the event handlers have run** and have called their `set` functions. This prevents multiple re-renders during a single event. In the rare case that you need to force React to update the screen earlier, for example to access the DOM, you can use [`flushSync`.](/reference/react-dom/flushSync)

* The `set` function has a stable identity, so you will often see it omitted from Effect dependencies, but including it will not cause the Effect to fire. If the linter lets you omit a dependency without errors, it is safe to do. [Learn more about removing Effect dependencies.](/learn/removing-effect-dependencies#move-dynamic-objects-and-functions-inside-your-effect)

* Calling the `set` function *during rendering* is only allowed from within the currently rendering component. React will discard its output and immediately attempt to render it again with the new state. This pattern is rarely needed, but you can use it to **store information from the previous renders**. [See an example below.](#storing-information-from-previous-renders)

* In Strict Mode, React will **call your updater function twice** in order to [help you find accidental impurities.](#my-initializer-or-updater-function-runs-twice) This is development-only behavior and does not affect production. If your updater function is pure (as it should be), this should not affect the behavior. The result from one of the calls will be ignored.

---

## Usage {/*usage*/}

### Adding state to a component {/*adding-state-to-a-component*/}

Call `useState` at the top level of your component to declare one or more [state variables.](/learn/state-a-components-memory)

```js [[1, 4, "age"], [2, 4, "setAge"], [3, 4, "42"], [1, 5, "name"], [2, 5, "setName"], [3, 5, "'Taylor'"]]
import { useState } from 'react';

function MyComponent() {
 const [age, setAge] = useState(42);
 const [name, setName] = useState('Taylor');
 // ...
```

The convention is to name state variables like `[something, setSomething]` using [array destructuring.](https://javascript.info/destructuring-assignment)

`useState` returns an array with exactly two items:

1. The <CodeStep step={1}>current state</CodeStep> of this state variable, initially set to the <CodeStep step={3}>initial state</CodeStep> you provided.
2. The <CodeStep step={2}>`set` function</CodeStep> that lets you change it to any other value in response to interaction.

To update what’s on the screen, call the `set` function with some next state:

```js [[2, 2, "setName"]]
function handleClick() {
 setName('Robin');
}
```

React will store the next state, render your component again with the new values, and update the UI.

<Pitfall>

Calling the `set` function [**does not** change the current state in the already executing code](#ive-updated-the-state-but-logging-gives-me-the-old-value):

```js {3}
function handleClick() {
 setName('Robin');
 console.log(name); // Still "Taylor"!
}
```

It only affects what `useState` will return starting from the *next* render.

</Pitfall>

<Recipes titleText="Basic useState examples" titleId="examples-basic">

#### Counter (number) {/*counter-number*/}

In this example, the `count` state variable holds a number. Clicking the button increments it.

<Sandpack>

```js
import { useState } from 'react';

export default function Counter() {
 const [count, setCount] = useState(0);

 function handleClick() {
 setCount(count + 1);
 }

 return (
 <button onClick={handleClick}>
 You pressed me {count} times
 </button>
 );
}
```

</Sandpack>

<Solution />

#### Text field (string) {/*text-field-string*/}

In this example, the `text` state variable holds a string. When you type, `handleChange` reads the latest input value from the browser input DOM element, and calls `setText` to update the state. This allows you to display the current `text` below.

<Sandpack>

```js
import { useState } from 'react';

export default function MyInput() {
 const [text, setText] = useState('hello');

 function handleChange(e) {
 setText(e.target.value);
 }

 return (
 <>
 <input value={text} onChange={handleChange} />
 <p>You typed: {text}</p>
 <button onClick={() => setText('hello')}>
 Reset
 </button>
 </>
 );
}
```

</Sandpack>

<Solution />

#### Checkbox (boolean) {/*checkbox-boolean*/}

In this example, the `liked` state variable holds a boolean. When you click the input, `setLiked` updates the `liked` state variable with whether the browser checkbox input is checked. The `liked` variable is used to render the text below the checkbox.

<Sandpack>

```js
import { useState } from 'react';

export default function MyCheckbox() {
 const [liked, setLiked] = useState(true);

 function handleChange(e) {
 setLiked(e.target.checked);
 }

 return (
 <>
 <label>
 <input
 type="checkbox"
 checked={liked}
 onChange={handleChange}
 />
 I liked this
 </label>
 <p>You {liked ? 'liked' : 'did not like'} this.</p>
 </>
 );
}
```

</Sandpack>

<Solution />

#### Form (two variables) {/*form-two-variables*/}

You can declare more than one state variable in the same component. Each state variable is completely independent.

<Sandpack>

```js
import { useState } from 'react';

export default function Form() {
 const [name, setName] = useState('Taylor');
 const [age, setAge] = useState(42);

 return (
 <>
 <input
 value={name}
 onChange={e => setName(e.target.value)}
 />
 <button onClick={() => setAge(age + 1)}>
 Increment age
 </button>
 <p>Hello, {name}. You are {age}.</p>
 </>
 );
}
```

```css
button { display: block; margin-top: 10px; }
```

</Sandpack>

<Solution />

</Recipes>

---

### Updating state based on the previous state {/*updating-state-based-on-the-previous-state*/}

Suppose the `age` is `42`. This handler calls `setAge(age + 1)` three times:

```js
function handleClick() {
 setAge(age + 1); // setAge(42 + 1)
 setAge(age + 1); // setAge(42 + 1)
 setAge(age + 1); // setAge(42 + 1)
}
```

However, after one click, `age` will only be `43` rather than `45`! This is because calling the `set` function [does not update](/learn/state-as-a-snapshot) the `age` state variable in the already running code. So each `setAge(age + 1)` call becomes `setAge(43)`.

To solve this problem, **you may pass an *updater function*** to `setAge` instead of the next state:

```js [[1, 2, "a", 0], [2, 2, "a + 1"], [1, 3, "a", 0], [2, 3, "a + 1"], [1, 4, "a", 0], [2, 4, "a + 1"]]
function handleClick() {
 setAge(a => a + 1); // setAge(42 => 43)
 setAge(a => a + 1); // setAge(43 => 44)
 setAge(a => a + 1); // setAge(44 => 45)
}
```

Here, `a => a + 1` is your updater function. It takes the <CodeStep step={1}>pending state</CodeStep> and calculates the <CodeStep step={2}>next state</CodeStep> from it.

React puts your updater functions in a [queue.](/learn/queueing-a-series-of-state-updates) Then, during the next render, it will call them in the same order:

1. `a => a + 1` will receive `42` as the pending state and return `43` as the next state.
1. `a => a + 1` will receive `43` as the pending state and return `44` as the next state.
1. `a => a + 1` will receive `44` as the pending state and return `45` as the next state.

There are no other queued updates, so React will store `45` as the current state in the end.

By convention, it's common to name the pending state argument for the first letter of the state variable name, like `a` for `age`. However, you may also call it like `prevAge` or something else that you find clearer.

React may [call your updaters twice](#my-initializer-or-updater-function-runs-twice) in development to verify that they are [pure.](/learn/keeping-components-pure)

<DeepDive>

#### Is using an updater always preferred? {/*is-using-an-updater-always-preferred*/}

You might hear a recommendation to always write code like `setAge(a => a + 1)` if the state you're setting is calculated from the previous state. There is no harm in it, but it is also not always necessary.

In most cases, there is no difference between these two approaches. React always makes sure that for intentional user actions, like clicks, the `age` state variable would be updated before the next click. This means there is no risk of a click handler seeing a "stale" `age` at the beginning of the event handler.

However, if you do multiple updates within the same event, updaters can be helpful. They're also helpful if accessing the state variable itself is inconvenient (you might run into this when optimizing re-renders).

If you prefer consistency over slightly more verbose syntax, it's reasonable to always write an updater if the state you're setting is calculated from the previous state. If it's calculated from the previous state of some *other* state variable, you might want to combine them into one object and [use a reducer.](/learn/extracting-state-logic-into-a-reducer)

</DeepDive>

<Recipes titleText="The difference between passing an updater and passing the next state directly" titleId="examples-updater">

#### Passing the updater function {/*passing-the-updater-function*/}

This example passes the updater function, so the "+3" button works.

<Sandpack>

```js
import { useState } from 'react';

export default function Counter() {
 const [age, setAge] = useState(42);

 function increment() {
 setAge(a => a + 1);
 }

 return (
 <>
 <h1>Your age: {age}</h1>
 <button onClick={() => {
 increment();
 increment();
 increment();
 }}>+3</button>
 <button onClick={() => {
 increment();
 }}>+1</button>
 </>
 );
}
```

```css
button { display: block; margin: 10px; font-size: 20px; }
h1 { display: block; margin: 10px; }
```

</Sandpack>

<Solution />

#### Passing the next state directly {/*passing-the-next-state-directly*/}

This example **does not** pass the updater function, so the "+3" button **doesn't work as intended**.

<Sandpack>

```js
import { useState } from 'react';

export default function Counter() {
 const [age, setAge] = useState(42);

 function increment() {
 setAge(age + 1);
 }

 return (
 <>
 <h1>Your age: {age}</h1>
 <button onClick={() => {
 increment();
 increment();
 increment();
 }}>+3</button>
 <button onClick={() => {
 increment();
 }}>+1</button>
 </>
 );
}
```

```css
button { display: block; margin: 10px; font-size: 20px; }
h1 { display: block; margin: 10px; }
```

</Sandpack>

<Solution />

</Recipes>

---

### Updating objects and arrays in state {/*updating-objects-and-arrays-in-state*/}

You can put objects and arrays into state. In React, state is considered read-only, so **you should *replace* it rather than *mutate* your existing objects**. For example, if you have a `form` object in state, don't mutate it:

```js
// 🚩 Don't mutate an object in state like this:
form.firstName = 'Taylor';
```

Instead, replace the whole object by creating a new one:

```js
// ✅ Replace state with a new object
setForm({
 ...form,
 firstName: 'Taylor'
});
```

Read [updating objects in state](/learn/updating-objects-in-state) and [updating arrays in state](/learn/updating-arrays-in-state) to learn more.

<Recipes titleText="Examples of objects and arrays in state" titleId="examples-objects">

#### Form (object) {/*form-object*/}

In this example, the `form` state variable holds an object. Each input has a change handler that calls `setForm` with the next state of the entire form. The `{ ...form }` spread syntax ensures that the state object is replaced rather than mutated.

<Sandpack>

```js
import { useState } from 'react';

export default function Form() {
 const [form, setForm] = useState({
 firstName: 'Barbara',
 lastName: 'Hepworth',
 email: 'bhepworth@sculpture.com',
 });

 return (
 <>
 <label>
 First name:
 <input
 value={form.firstName}
 onChange={e => {
 setForm({
 ...form,
 firstName: e.target.value
 });
 }}
 />
 </label>
 <label>
 Last name:
 <input
 value={form.lastName}
 onChange={e => {
 setForm({
 ...form,
 lastName: e.target.value
 });
 }}
 />
 </label>
 <label>
 Email:
 <input
 value={form.email}
 onChange={e => {
 setForm({
 ...form,
 email: e.target.value
 });
 }}
 />
 </label>
 <p>
 {form.firstName}{' '}
 {form.lastName}{' '}
 ({form.email})
 </p>
 </>
 );
}
```

```css
label { display: block; }
input { margin-left: 5px; }
```

</Sandpack>

<Solution />

#### Form (nested object) {/*form-nested-object*/}

In this example, the state is more nested. When you update nested state, you need to create a copy of the object you're updating, as well as any objects "containing" it on the way upwards. Read [updating a nested object](/learn/updating-objects-in-state#updating-a-nested-object) to learn more.

<Sandpack>

```js
import { useState } from 'react';

export default function Form() {
 const [person, setPerson] = useState({
 name: 'Niki de Saint Phalle',
 artwork: {
 title: 'Blue Nana',
 city: 'Hamburg',
 image: 'https://react.dev/images/docs/scientists/Sd1AgUOm.jpg',
 }
 });

 function handleNameChange(e) {
 setPerson({
 ...person,
 name: e.target.value
 });
 }

 function handleTitleChange(e) {
 setPerson({
 ...person,
 artwork: {
 ...person.artwork,
 title: e.target.value
 }
 });
 }

 function handleCityChange(e) {
 setPerson({
 ...person,
 artwork: {
 ...person.artwork,
 city: e.target.value
 }
 });
 }

 function handleImageChange(e) {
 setPerson({
 ...person,
 artwork: {
 ...person.artwork,
 image: e.target.value
 }
 });
 }

 return (
 <>
 <label>
 Name:
 <input
 value={person.name}
 onChange={handleNameChange}
 />
 </label>
 <label>
 Title:
 <input
 value={person.artwork.title}
 onChange={handleTitleChange}
 />
 </label>
 <label>
 City:
 <input
 value={person.artwork.city}
 onChange={handleCityChange}
 />
 </label>
 <label>
 Image:
 <input
 value={person.artwork.image}
 onChange={handleImageChange}
 />
 </label>
 <p>
 <i>{person.artwork.title}</i>
 {' by '}
 {person.name}
 <br />
 (located in {person.artwork.city})
 </p>
 <img
 src={person.artwork.image}
 alt={person.artwork.title}
 />
 </>
 );
}
```

```css
label { display: block; }
input { margin-left: 5px; margin-bottom: 5px; }
img { width: 200px; height: 200px; }
```

</Sandpack>

<Solution />

#### List (array) {/*list-array*/}

In this example, the `todos` state variable holds an array. Each button handler calls `setTodos` with the next version of that array. The `[...todos]` spread syntax, `todos.map()` and `todos.filter()` ensure the state array is replaced rather than mutated.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import AddTodo from './AddTodo.js';
import TaskList from './TaskList.js';

let nextId = 3;
const initialTodos = [
 { id: 0, title: 'Buy milk', done: true },
 { id: 1, title: 'Eat tacos', done: false },
 { id: 2, title: 'Brew tea', done: false },
];

export default function TaskApp() {
 const [todos, setTodos] = useState(initialTodos);

 function handleAddTodo(title) {
 setTodos([
 ...todos,
 {
 id: nextId++,
 title: title,
 done: false
 }
 ]);
 }

 function handleChangeTodo(nextTodo) {
 setTodos(todos.map(t => {
 if (t.id === nextTodo.id) {
 return nextTodo;
 } else {
 return t;
 }
 }));
 }

 function handleDeleteTodo(todoId) {
 setTodos(
 todos.filter(t => t.id !== todoId)
 );
 }

 return (
 <>
 <AddTodo
 onAddTodo={handleAddTodo}
 />
 <TaskList
 todos={todos}
 onChangeTodo={handleChangeTodo}
 onDeleteTodo={handleDeleteTodo}
 />
 </>
 );
}
```

```js src/AddTodo.js
import { useState } from 'react';

export default function AddTodo({ onAddTodo }) {
 const [title, setTitle] = useState('');
 return (
 <>
 <input
 placeholder="Add todo"
 value={title}
 onChange={e => setTitle(e.target.value)}
 />
 <button onClick={() => {
 setTitle('');
 onAddTodo(title);
 }}>Add</button>
 </>
 )
}
```

```js src/TaskList.js
import { useState } from 'react';

export default function TaskList({
 todos,
 onChangeTodo,
 onDeleteTodo
}) {
 return (
 <ul>
 {todos.map(todo => (
 <li key={todo.id}>
 <Task
 todo={todo}
 onChange={onChangeTodo}
 onDelete={onDeleteTodo}
 />
 </li>
 ))}
 </ul>
 );
}

function Task({ todo, onChange, onDelete }) {
 const [isEditing, setIsEditing] = useState(false);
 let todoContent;
 if (isEditing) {
 todoContent = (
 <>
 <input
 value={todo.title}
 onChange={e => {
 onChange({
 ...todo,
 title: e.target.value
 });
 }} />
 <button onClick={() => setIsEditing(false)}>
 Save
 </button>
 </>
 );
 } else {
 todoContent = (
 <>
 {todo.title}
 <button onClick={() => setIsEditing(true)}>
 Edit
 </button>
 </>
 );
 }
 return (
 <label>
 <input
 type="checkbox"
 checked={todo.done}
 onChange={e => {
 onChange({
 ...todo,
 done: e.target.checked
 });
 }}
 />
 {todoContent}
 <button onClick={() => onDelete(todo.id)}>
 Delete
 </button>
 </label>
 );
}
```

```css
button { margin: 5px; }
li { list-style-type: none; }
ul, li { margin: 0; padding: 0; }
```

</Sandpack>

<Solution />

#### Writing concise update logic with Immer {/*writing-concise-update-logic-with-immer*/}

If updating arrays and objects without mutation feels tedious, you can use a library like [Immer](https://github.com/immerjs/use-immer) to reduce repetitive code. Immer lets you write concise code as if you were mutating objects, but under the hood it performs immutable updates:

<Sandpack>

```js
import { useState } from 'react';
import { useImmer } from 'use-immer';

let nextId = 3;
const initialList = [
 { id: 0, title: 'Big Bellies', seen: false },
 { id: 1, title: 'Lunar Landscape', seen: false },
 { id: 2, title: 'Terracotta Army', seen: true },
];

export default function BucketList() {
 const [list, updateList] = useImmer(initialList);

 function handleToggle(artworkId, nextSeen) {
 updateList(draft => {
 const artwork = draft.find(a =>
 a.id === artworkId
 );
 artwork.seen = nextSeen;
 });
 }

 return (
 <>
 <h1>Art Bucket List</h1>
 <h2>My list of art to see:</h2>
 <ItemList
 artworks={list}
 onToggle={handleToggle} />
 </>
 );
}

function ItemList({ artworks, onToggle }) {
 return (
 <ul>
 {artworks.map(artwork => (
 <li key={artwork.id}>
 <label>
 <input
 type="checkbox"
 checked={artwork.seen}
 onChange={e => {
 onToggle(
 artwork.id,
 e.target.checked
 );
 }}
 />
 {artwork.title}
 </label>
 </li>
 ))}
 </ul>
 );
}
```

```json package.json
{
 "dependencies": {
 "immer": "1.7.3",
 "react": "latest",
 "react-dom": "latest",
 "react-scripts": "latest",
 "use-immer": "0.5.1"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Avoiding recreating the initial state {/*avoiding-recreating-the-initial-state*/}

React saves the initial state once and ignores it on the next renders.

```js
function TodoList() {
 const [todos, setTodos] = useState(createInitialTodos());
 // ...
```

Although the result of `createInitialTodos()` is only used for the initial render, you're still calling this function on every render. This can be wasteful if it's creating large arrays or performing expensive calculations.

To solve this, you may **pass it as an _initializer_ function** to `useState` instead:

```js
function TodoList() {
 const [todos, setTodos] = useState(createInitialTodos);
 // ...
```

Notice that you’re passing `createInitialTodos`, which is the *function itself*, and not `createInitialTodos()`, which is the result of calling it. If you pass a function to `useState`, React will only call it during initialization.

React may [call your initializers twice](#my-initializer-or-updater-function-runs-twice) in development to verify that they are [pure.](/learn/keeping-components-pure)

<Recipes titleText="The difference between passing an initializer and passing the initial state directly" titleId="examples-initializer">

#### Passing the initializer function {/*passing-the-initializer-function*/}

This example passes the initializer function, so the `createInitialTodos` function only runs during initialization. It does not run when component re-renders, such as when you type into the input.

<Sandpack>

```js
import { useState } from 'react';

function createInitialTodos() {
 const initialTodos = [];
 for (let i = 0; i < 50; i++) {
 initialTodos.push({
 id: i,
 text: 'Item ' + (i + 1)
 });
 }
 return initialTodos;
}

export default function TodoList() {
 const [todos, setTodos] = useState(createInitialTodos);
 const [text, setText] = useState('');

 return (
 <>
 <input
 value={text}
 onChange={e => setText(e.target.value)}
 />
 <button onClick={() => {
 setText('');
 setTodos([{
 id: todos.length,
 text: text
 }, ...todos]);
 }}>Add</button>
 <ul>
 {todos.map(item => (
 <li key={item.id}>
 {item.text}
 </li>
 ))}
 </ul>
 </>
 );
}
```

</Sandpack>

<Solution />

#### Passing the initial state directly {/*passing-the-initial-state-directly*/}

This example **does not** pass the initializer function, so the `createInitialTodos` function runs on every render, such as when you type into the input. There is no observable difference in behavior, but this code is less efficient.

<Sandpack>

```js
import { useState } from 'react';

function createInitialTodos() {
 const initialTodos = [];
 for (let i = 0; i < 50; i++) {
 initialTodos.push({
 id: i,
 text: 'Item ' + (i + 1)
 });
 }
 return initialTodos;
}

export default function TodoList() {
 const [todos, setTodos] = useState(createInitialTodos());
 const [text, setText] = useState('');

 return (
 <>
 <input
 value={text}
 onChange={e => setText(e.target.value)}
 />
 <button onClick={() => {
 setText('');
 setTodos([{
 id: todos.length,
 text: text
 }, ...todos]);
 }}>Add</button>
 <ul>
 {todos.map(item => (
 <li key={item.id}>
 {item.text}
 </li>
 ))}
 </ul>
 </>
 );
}
```

</Sandpack>

<Solution />

</Recipes>

---

### Resetting state with a key {/*resetting-state-with-a-key*/}

You'll often encounter the `key` attribute when [rendering lists.](/learn/rendering-lists) However, it also serves another purpose.

You can **reset a component's state by passing a different `key` to a component.** In this example, the Reset button changes the `version` state variable, which we pass as a `key` to the `Form`. When the `key` changes, React re-creates the `Form` component (and all of its children) from scratch, so its state gets reset.

Read [preserving and resetting state](/learn/preserving-and-resetting-state) to learn more.

<Sandpack>

```js src/App.js
import { useState } from 'react';

export default function App() {
 const [version, setVersion] = useState(0);

 function handleReset() {
 setVersion(version + 1);
 }

 return (
 <>
 <button onClick={handleReset}>Reset</button>
 <Form key={version} />
 </>
 );
}

function Form() {
 const [name, setName] = useState('Taylor');

 return (
 <>
 <input
 value={name}
 onChange={e => setName(e.target.value)}
 />
 <p>Hello, {name}.</p>
 </>
 );
}
```

```css
button { display: block; margin-bottom: 20px; }
```

</Sandpack>

---

### Storing information from previous renders {/*storing-information-from-previous-renders*/}

Usually, you will update state in event handlers. However, in rare cases you might want to adjust state in response to rendering -- for example, you might want to change a state variable when a prop changes.

In most cases, you don't need this:

* **If the value you need can be computed entirely from the current props or other state, [remove that redundant state altogether.](/learn/choosing-the-state-structure#avoid-redundant-state)** If you're worried about recomputing too often, the [`useMemo` Hook](/reference/react/useMemo) can help.
* If you want to reset the entire component tree's state, [pass a different `key` to your component.](#resetting-state-with-a-key)
* If you can, update all the relevant state in the event handlers.

In the rare case that none of these apply, there is a pattern you can use to update state based on the values that have been rendered so far, by calling a `set` function while your component is rendering.

Here's an example. This `CountLabel` component displays the `count` prop passed to it:

```js src/CountLabel.js
export default function CountLabel({ count }) {
 return <h1>{count}</h1>
}
```

Say you want to show whether the counter has *increased or decreased* since the last change. The `count` prop doesn't tell you this -- you need to keep track of its previous value. Add the `prevCount` state variable to track it. Add another state variable called `trend` to hold whether the count has increased or decreased. Compare `prevCount` with `count`, and if they're not equal, update both `prevCount` and `trend`. Now you can show both the current count prop and *how it has changed since the last render*.

<Sandpack>

```js src/App.js
import { useState } from 'react';
import CountLabel from './CountLabel.js';

export default function App() {
 const [count, setCount] = useState(0);
 return (
 <>
 <button onClick={() => setCount(count + 1)}>
 Increment
 </button>
 <button onClick={() => setCount(count - 1)}>
 Decrement
 </button>
 <CountLabel count={count} />
 </>
 );
}
```

```js src/CountLabel.js active
import { useState } from 'react';

export default function CountLabel({ count }) {
 const [prevCount, setPrevCount] = useState(count);
 const [trend, setTrend] = useState(null);
 if (prevCount !== count) {
 setPrevCount(count);
 setTrend(count > prevCount ? 'increasing' : 'decreasing');
 }
 return (
 <>
 <h1>{count}</h1>
 {trend && <p>The count is {trend}</p>}
 </>
 );
}
```

```css
button { margin-bottom: 10px; }
```

</Sandpack>

Note that if you call a `set` function while rendering, it must be inside a condition like `prevCount !== count`, and there must be a call like `setPrevCount(count)` inside of the condition. Otherwise, your component would re-render in a loop until it crashes. Also, you can only update the state of the *currently rendering* component like this. Calling the `set` function of *another* component during rendering is an error. Finally, your `set` call should still [update state without mutation](#updating-objects-and-arrays-in-state) -- this doesn't mean you can break other rules of [pure functions.](/learn/keeping-components-pure)

This pattern can be hard to understand and is usually best avoided. However, it's better than updating state in an effect. When you call the `set` function during render, React will re-render that component immediately after your component exits with a `return` statement, and before rendering the children. This way, children don't need to render twice. The rest of your component function will still execute (and the result will be thrown away). If your condition is below all the Hook calls, you may add an early `return;` to restart rendering earlier.

---

## Troubleshooting {/*troubleshooting*/}

### I've updated the state, but logging gives me the old value {/*ive-updated-the-state-but-logging-gives-me-the-old-value*/}

Calling the `set` function **does not change state in the running code**:

```js {4,5,8}
function handleClick() {
 console.log(count); // 0

 setCount(count + 1); // Request a re-render with 1
 console.log(count); // Still 0!

 setTimeout(() => {
 console.log(count); // Also 0!
 }, 5000);
}
```

This is because [states behaves like a snapshot.](/learn/state-as-a-snapshot) Updating state requests another render with the new state value, but does not affect the `count` JavaScript variable in your already-running event handler.

If you need to use the next state, you can save it in a variable before passing it to the `set` function:

```js
const nextCount = count + 1;
setCount(nextCount);

console.log(count); // 0
console.log(nextCount); // 1
```

---

### I've updated the state, but the screen doesn't update {/*ive-updated-the-state-but-the-screen-doesnt-update*/}

React will **ignore your update if the next state is equal to the previous state,** as determined by an [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is) comparison. This usually happens when you change an object or an array in state directly:

```js
obj.x = 10; // 🚩 Wrong: mutating existing object
setObj(obj); // 🚩 Doesn't do anything
```

You mutated an existing `obj` object and passed it back to `setObj`, so React ignored the update. To fix this, you need to ensure that you're always [_replacing_ objects and arrays in state instead of _mutating_ them](#updating-objects-and-arrays-in-state):

```js
// ✅ Correct: creating a new object
setObj({
 ...obj,
 x: 10
});
```

---

### I'm getting an error: "Too many re-renders" {/*im-getting-an-error-too-many-re-renders*/}

You might get an error that says: `Too many re-renders. React limits the number of renders to prevent an infinite loop.` Typically, this means that you're unconditionally setting state *during render*, so your component enters a loop: render, set state (which causes a render), render, set state (which causes a render), and so on. Very often, this is caused by a mistake in specifying an event handler:

```js {1-2}
// 🚩 Wrong: calls the handler during render
return <button onClick={handleClick()}>Click me</button>

// ✅ Correct: passes down the event handler
return <button onClick={handleClick}>Click me</button>

// ✅ Correct: passes down an inline function
return <button onClick={(e) => handleClick(e)}>Click me</button>
```

If you can't find the cause of this error, click on the arrow next to the error in the console and look through the JavaScript stack to find the specific `set` function call responsible for the error.

---

### My initializer or updater function runs twice {/*my-initializer-or-updater-function-runs-twice*/}

In [Strict Mode](/reference/react/StrictMode), React will call some of your functions twice instead of once:

```js {2,5-6,11-12}
function TodoList() {
 // This component function will run twice for every render.

 const [todos, setTodos] = useState(() => {
 // This initializer function will run twice during initialization.
 return createTodos();
 });

 function handleClick() {
 setTodos(prevTodos => {
 // This updater function will run twice for every click.
 return [...prevTodos, createTodo()];
 });
 }
 // ...
```

This is expected and shouldn't break your code.

This **development-only** behavior helps you [keep components pure.](/learn/keeping-components-pure) React uses the result of one of the calls, and ignores the result of the other call. As long as your component, initializer, and updater functions are pure, this shouldn't affect your logic. However, if they are accidentally impure, this helps you notice the mistakes.

For example, this impure updater function mutates an array in state:

```js {2,3}
setTodos(prevTodos => {
 // 🚩 Mistake: mutating state
 prevTodos.push(createTodo());
});
```

Because React calls your updater function twice, you'll see the todo was added twice, so you'll know that there is a mistake. In this example, you can fix the mistake by [replacing the array instead of mutating it](#updating-objects-and-arrays-in-state):

```js {2,3}
setTodos(prevTodos => {
 // ✅ Correct: replacing with new state
 return [...prevTodos, createTodo()];
});
```

Now that this updater function is pure, calling it an extra time doesn't make a difference in behavior. This is why React calling it twice helps you find mistakes. **Only component, initializer, and updater functions need to be pure.** Event handlers don't need to be pure, so React will never call your event handlers twice.

Read [keeping components pure](/learn/keeping-components-pure) to learn more.

---

### I'm trying to set state to a function, but it gets called instead {/*im-trying-to-set-state-to-a-function-but-it-gets-called-instead*/}

You can't put a function into state like this:

```js
const [fn, setFn] = useState(someFunction);

function handleClick() {
 setFn(someOtherFunction);
}
```

Because you're passing a function, React assumes that `someFunction` is an [initializer function](#avoiding-recreating-the-initial-state), and that `someOtherFunction` is an [updater function](#updating-state-based-on-the-previous-state), so it tries to call them and store the result. To actually *store* a function, you have to put `() =>` before them in both cases. Then React will store the functions you pass.

```js {1,4}
const [fn, setFn] = useState(() => someFunction);

function handleClick() {
 setFn(() => someOtherFunction);
}
```

---
title: useSyncExternalStore
---

<Intro>

`useSyncExternalStore` is a React Hook that lets you subscribe to an external store.

```js
const snapshot = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot?)
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot?)` {/*usesyncexternalstore*/}

Call `useSyncExternalStore` at the top level of your component to read a value from an external data store.

```js
import { useSyncExternalStore } from 'react';
import { todosStore } from './todoStore.js';

function TodosApp() {
 const todos = useSyncExternalStore(todosStore.subscribe, todosStore.getSnapshot);
 // ...
}
```

It returns the snapshot of the data in the store. You need to pass two functions as arguments:

1. The `subscribe` function should subscribe to the store and return a function that unsubscribes.
2. The `getSnapshot` function should read a snapshot of the data from the store.

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

* `subscribe`: A function that takes a single `callback` argument and subscribes it to the store. When the store changes, it should invoke the provided `callback`, which will cause React to re-call `getSnapshot` and (if needed) re-render the component. The `subscribe` function should return a function that cleans up the subscription.

* `getSnapshot`: A function that returns a snapshot of the data in the store that's needed by the component. While the store has not changed, repeated calls to `getSnapshot` must return the same value. If the store changes and the returned value is different (as compared by [`Object.is`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Object/is)), React re-renders the component.

* **optional** `getServerSnapshot`: A function that returns the initial snapshot of the data in the store. It will be used only during server rendering and during hydration of server-rendered content on the client. The server snapshot must be the same between the client and the server, and is usually serialized and passed from the server to the client. If you omit this argument, rendering the component on the server will throw an error.

#### Returns {/*returns*/}

The current snapshot of the store which you can use in your rendering logic.

#### Caveats {/*caveats*/}

* The store snapshot returned by `getSnapshot` must be immutable. If the underlying store has mutable data, return a new immutable snapshot if the data has changed. Otherwise, return a cached last snapshot.

* If a different `subscribe` function is passed during a re-render, React will re-subscribe to the store using the newly passed `subscribe` function. You can prevent this by declaring `subscribe` outside the component.

* If the store is mutated during a [non-blocking Transition update](/reference/react/useTransition), React will fall back to performing that update as blocking. Specifically, for every Transition update, React will call `getSnapshot` a second time just before applying changes to the DOM. If it returns a different value than when it was called originally, React will restart the update from scratch, this time applying it as a blocking update, to ensure that every component on screen is reflecting the same version of the store.

* It's not recommended to _suspend_ a render based on a store value returned by `useSyncExternalStore`. The reason is that mutations to the external store cannot be marked as [non-blocking Transition updates](/reference/react/useTransition), so they will trigger the nearest [`Suspense` fallback](/reference/react/Suspense), replacing already-rendered content on screen with a loading spinner, which typically makes a poor UX.

 For example, the following are discouraged:

 ```js
 const LazyProductDetailPage = lazy(() => import('./ProductDetailPage.js'));

 function ShoppingApp() {
 const selectedProductId = useSyncExternalStore(...);

 // ❌ Calling `use` with a Promise dependent on `selectedProductId`
 const data = use(fetchItem(selectedProductId))

 // ❌ Conditionally rendering a lazy component based on `selectedProductId`
 return selectedProductId != null ? <LazyProductDetailPage /> : <FeaturedProducts />;
 }
 ```

---

## Usage {/*usage*/}

### Subscribing to an external store {/*subscribing-to-an-external-store*/}

Most of your React components will only read data from their [props,](/learn/passing-props-to-a-component) [state,](/reference/react/useState) and [context.](/reference/react/useContext) However, sometimes a component needs to read some data from some store outside of React that changes over time. This includes:

* Third-party state management libraries that hold state outside of React.
* Browser APIs that expose a mutable value and events to subscribe to its changes.

Call `useSyncExternalStore` at the top level of your component to read a value from an external data store.

```js [[1, 5, "todosStore.subscribe"], [2, 5, "todosStore.getSnapshot"], [3, 5, "todos", 0]]
import { useSyncExternalStore } from 'react';
import { todosStore } from './todoStore.js';

function TodosApp() {
 const todos = useSyncExternalStore(todosStore.subscribe, todosStore.getSnapshot);
 // ...
}
```

It returns the <CodeStep step={3}>snapshot</CodeStep> of the data in the store. You need to pass two functions as arguments:

1. The <CodeStep step={1}>`subscribe` function</CodeStep> should subscribe to the store and return a function that unsubscribes.
2. The <CodeStep step={2}>`getSnapshot` function</CodeStep> should read a snapshot of the data from the store.

React will use these functions to keep your component subscribed to the store and re-render it on changes.

For example, in the sandbox below, `todosStore` is implemented as an external store that stores data outside of React. The `TodosApp` component connects to that external store with the `useSyncExternalStore` Hook.

<Sandpack>

```js
import { useSyncExternalStore } from 'react';
import { todosStore } from './todoStore.js';

export default function TodosApp() {
 const todos = useSyncExternalStore(todosStore.subscribe, todosStore.getSnapshot);
 return (
 <>
 <button onClick={() => todosStore.addTodo()}>Add todo</button>
 <hr />
 <ul>
 {todos.map(todo => (
 <li key={todo.id}>{todo.text}</li>
 ))}
 </ul>
 </>
 );
}
```

```js src/todoStore.js
// This is an example of a third-party store
// that you might need to integrate with React.

// If your app is fully built with React,
// we recommend using React state instead.

let nextId = 0;
let todos = [{ id: nextId++, text: 'Todo #1' }];
let listeners = [];

export const todosStore = {
 addTodo() {
 todos = [...todos, { id: nextId++, text: 'Todo #' + nextId }]
 emitChange();
 },
 subscribe(listener) {
 listeners = [...listeners, listener];
 return () => {
 listeners = listeners.filter(l => l !== listener);
 };
 },
 getSnapshot() {
 return todos;
 }
};

function emitChange() {
 for (let listener of listeners) {
 listener();
 }
}
```

</Sandpack>

<Note>

When possible, we recommend using built-in React state with [`useState`](/reference/react/useState) and [`useReducer`](/reference/react/useReducer) instead. The `useSyncExternalStore` API is mostly useful if you need to integrate with existing non-React code.

</Note>

---

### Subscribing to a browser API {/*subscribing-to-a-browser-api*/}

Another reason to add `useSyncExternalStore` is when you want to subscribe to some value exposed by the browser that changes over time. For example, suppose that you want your component to display whether the network connection is active. The browser exposes this information via a property called [`navigator.onLine`.](https://developer.mozilla.org/en-US/docs/Web/API/Navigator/onLine)

This value can change without React's knowledge, so you should read it with `useSyncExternalStore`.

```js
import { useSyncExternalStore } from 'react';

function ChatIndicator() {
 const isOnline = useSyncExternalStore(subscribe, getSnapshot);
 // ...
}
```

To implement the `getSnapshot` function, read the current value from the browser API:

```js
function getSnapshot() {
 return navigator.onLine;
}
```

Next, you need to implement the `subscribe` function. For example, when `navigator.onLine` changes, the browser fires the [`online`](https://developer.mozilla.org/en-US/docs/Web/API/Window/online_event) and [`offline`](https://developer.mozilla.org/en-US/docs/Web/API/Window/offline_event) events on the `window` object. You need to subscribe the `callback` argument to the corresponding events, and then return a function that cleans up the subscriptions:

```js
function subscribe(callback) {
 window.addEventListener('online', callback);
 window.addEventListener('offline', callback);
 return () => {
 window.removeEventListener('online', callback);
 window.removeEventListener('offline', callback);
 };
}
```

Now React knows how to read the value from the external `navigator.onLine` API and how to subscribe to its changes. Disconnect your device from the network and notice that the component re-renders in response:

<Sandpack>

```js
import { useSyncExternalStore } from 'react';

export default function ChatIndicator() {
 const isOnline = useSyncExternalStore(subscribe, getSnapshot);
 return <h1>{isOnline ? '✅ Online' : '❌ Disconnected'}</h1>;
}

function getSnapshot() {
 return navigator.onLine;
}

function subscribe(callback) {
 window.addEventListener('online', callback);
 window.addEventListener('offline', callback);
 return () => {
 window.removeEventListener('online', callback);
 window.removeEventListener('offline', callback);
 };
}
```

</Sandpack>

---

### Extracting the logic to a custom Hook {/*extracting-the-logic-to-a-custom-hook*/}

Usually you won't write `useSyncExternalStore` directly in your components. Instead, you'll typically call it from your own custom Hook. This lets you use the same external store from different components.

For example, this custom `useOnlineStatus` Hook tracks whether the network is online:

```js {3,6}
import { useSyncExternalStore } from 'react';

export function useOnlineStatus() {
 const isOnline = useSyncExternalStore(subscribe, getSnapshot);
 return isOnline;
}

function getSnapshot() {
 // ...
}

function subscribe(callback) {
 // ...
}
```

Now different components can call `useOnlineStatus` without repeating the underlying implementation:

<Sandpack>

```js
import { useOnlineStatus } from './useOnlineStatus.js';

function StatusBar() {
 const isOnline = useOnlineStatus();
 return <h1>{isOnline ? '✅ Online' : '❌ Disconnected'}</h1>;
}

function SaveButton() {
 const isOnline = useOnlineStatus();

 function handleSaveClick() {
 console.log('✅ Progress saved');
 }

 return (
 <button disabled={!isOnline} onClick={handleSaveClick}>
 {isOnline ? 'Save progress' : 'Reconnecting...'}
 </button>
 );
}

export default function App() {
 return (
 <>
 <SaveButton />
 <StatusBar />
 </>
 );
}
```

```js src/useOnlineStatus.js
import { useSyncExternalStore } from 'react';

export function useOnlineStatus() {
 const isOnline = useSyncExternalStore(subscribe, getSnapshot);
 return isOnline;
}

function getSnapshot() {
 return navigator.onLine;
}

function subscribe(callback) {
 window.addEventListener('online', callback);
 window.addEventListener('offline', callback);
 return () => {
 window.removeEventListener('online', callback);
 window.removeEventListener('offline', callback);
 };
}
```

</Sandpack>

---

### Adding support for server rendering {/*adding-support-for-server-rendering*/}

If your React app uses [server rendering,](/reference/react-dom/server) your React components will also run outside the browser environment to generate the initial HTML. This creates a few challenges when connecting to an external store:

- If you're connecting to a browser-only API, it won't work because it does not exist on the server.
- If you're connecting to a third-party data store, you'll need its data to match between the server and client.

To solve these issues, pass a `getServerSnapshot` function as the third argument to `useSyncExternalStore`:

```js {4,12-14}
import { useSyncExternalStore } from 'react';

export function useOnlineStatus() {
 const isOnline = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
 return isOnline;
}

function getSnapshot() {
 return navigator.onLine;
}

function getServerSnapshot() {
 return true; // Always show "Online" for server-generated HTML
}

function subscribe(callback) {
 // ...
}
```

The `getServerSnapshot` function is similar to `getSnapshot`, but it runs only in two situations:

- It runs on the server when generating the HTML.
- It runs on the client during [hydration](/reference/react-dom/client/hydrateRoot), i.e. when React takes the server HTML and makes it interactive.

This lets you provide the initial snapshot value which will be used before the app becomes interactive. If there is no meaningful initial value for the server rendering, omit this argument to [force rendering on the client.](/reference/react/Suspense#providing-a-fallback-for-server-errors-and-client-only-content)

<Note>

Make sure that `getServerSnapshot` returns the same exact data on the initial client render as it returned on the server. For example, if `getServerSnapshot` returned some prepopulated store content on the server, you need to transfer this content to the client. One way to do this is to emit a `<script>` tag during server rendering that sets a global like `window.MY_STORE_DATA`, and read from that global on the client in `getServerSnapshot`. Your external store should provide instructions on how to do that.

</Note>

---

## Troubleshooting {/*troubleshooting*/}

### I'm getting an error: "The result of `getSnapshot` should be cached" {/*im-getting-an-error-the-result-of-getsnapshot-should-be-cached*/}

This error means your `getSnapshot` function returns a new object every time it's called, for example:

```js {2-5}
function getSnapshot() {
 // 🔴 Do not return always different objects from getSnapshot
 return {
 todos: myStore.todos
 };
}
```

React will re-render the component if `getSnapshot` return value is different from the last time. This is why, if you always return a different value, you will enter an infinite loop and get this error.

Your `getSnapshot` object should only return a different object if something has actually changed. If your store contains immutable data, you can return that data directly:

```js {2-3}
function getSnapshot() {
 // ✅ You can return immutable data
 return myStore.todos;
}
```

If your store data is mutable, your `getSnapshot` function should return an immutable snapshot of it. This means it *does* need to create new objects, but it shouldn't do this for every single call. Instead, it should store the last calculated snapshot, and return the same snapshot as the last time if the data in the store has not changed. How you determine whether mutable data has changed depends on your mutable store.

---

### My `subscribe` function gets called after every re-render {/*my-subscribe-function-gets-called-after-every-re-render*/}

This `subscribe` function is defined *inside* a component so it is different on every re-render:

```js {2-5}
function ChatIndicator() {
 // 🚩 Always a different function, so React will resubscribe on every re-render
 function subscribe() {
 // ...
 }

 const isOnline = useSyncExternalStore(subscribe, getSnapshot);

 // ...
}
```

React will resubscribe to your store if you pass a different `subscribe` function between re-renders. If this causes performance issues and you'd like to avoid resubscribing, move the `subscribe` function outside:

```js {1-4}
// ✅ Always the same function, so React won't need to resubscribe
function subscribe() {
 // ...
}

function ChatIndicator() {
 const isOnline = useSyncExternalStore(subscribe, getSnapshot);
 // ...
}
```

Alternatively, wrap `subscribe` into [`useCallback`](/reference/react/useCallback) to only resubscribe when some argument changes:

```js {2-5}
function ChatIndicator({ userId }) {
 // ✅ Same function as long as userId doesn't change
 const subscribe = useCallback(() => {
 // ...
 }, [userId]);

 const isOnline = useSyncExternalStore(subscribe, getSnapshot);

 // ...
}
```

---
title: useTransition
---

<Intro>

`useTransition` is a React Hook that lets you render a part of the UI in the background.

```js
const [isPending, startTransition] = useTransition()
```

</Intro>

<InlineToc />

---

## Reference {/*reference*/}

### `useTransition()` {/*usetransition*/}

Call `useTransition` at the top level of your component to mark some state updates as Transitions.

```js
import { useTransition } from 'react';

function TabContainer() {
 const [isPending, startTransition] = useTransition();
 // ...
}
```

[See more examples below.](#usage)

#### Parameters {/*parameters*/}

`useTransition` does not take any parameters.

#### Returns {/*returns*/}

`useTransition` returns an array with exactly two items:

1. The `isPending` flag that tells you whether there is a pending Transition.
2. The [`startTransition` function](#starttransition) that lets you mark updates as a Transition.

---

### `startTransition(action)` {/*starttransition*/}

The `startTransition` function returned by `useTransition` lets you mark an update as a Transition.

```js {6,8}
function TabContainer() {
 const [isPending, startTransition] = useTransition();
 const [tab, setTab] = useState('about');

 function selectTab(nextTab) {
 startTransition(() => {
 setTab(nextTab);
 });
 }
 // ...
}
```

<Note>
#### Functions called in `startTransition` are called "Actions". {/*functions-called-in-starttransition-are-called-actions*/}

The function passed to `startTransition` is called an "Action". By convention, any callback called inside `startTransition` (such as a callback prop) should be named `action` or include the "Action" suffix:

```js {1,9}
function SubmitButton({ submitAction }) {
 const [isPending, startTransition] = useTransition();

 return (
 <button
 disabled={isPending}
 onClick={() => {
 startTransition(async () => {
 await submitAction();
 });
 }}
 >
 Submit
 </button>
 );
}

```

</Note>

#### Parameters {/*starttransition-parameters*/}

* `action`: A function that updates some state by calling one or more [`set` functions](/reference/react/useState#setstate). React calls `action` immediately with no parameters and marks all state updates scheduled synchronously during the `action` function call as Transitions. Any async calls that are awaited in the `action` will be included in the Transition, but currently require wrapping any `set` functions after the `await` in an additional `startTransition` (see [Troubleshooting](#react-doesnt-treat-my-state-update-after-await-as-a-transition)). State updates marked as Transitions will be [non-blocking](#perform-non-blocking-updates-with-actions) and [will not display unwanted loading indicators](#preventing-unwanted-loading-indicators).

#### Returns {/*starttransition-returns*/}

`startTransition` does not return anything.

#### Caveats {/*starttransition-caveats*/}

* `useTransition` is a Hook, so it can only be called inside components or custom Hooks. If you need to start a Transition somewhere else (for example, from a data library), call the standalone [`startTransition`](/reference/react/startTransition) instead.

* You can wrap an update into a Transition only if you have access to the `set` function of that state. If you want to start a Transition in response to some prop or a custom Hook value, try [`useDeferredValue`](/reference/react/useDeferredValue) instead.

* The function you pass to `startTransition` is called immediately, marking all state updates that happen while it executes as Transitions. If you try to perform state updates in a `setTimeout`, for example, they won't be marked as Transitions.

* You must wrap any state updates after any async requests in another `startTransition` to mark them as Transitions. This is a known limitation that we will fix in the future (see [Troubleshooting](#react-doesnt-treat-my-state-update-after-await-as-a-transition)).

* The `startTransition` function has a stable identity, so you will often see it omitted from Effect dependencies, but including it will not cause the Effect to fire. If the linter lets you omit a dependency without errors, it is safe to do. [Learn more about removing Effect dependencies.](/learn/removing-effect-dependencies#move-dynamic-objects-and-functions-inside-your-effect)

* A state update marked as a Transition will be interrupted by other state updates. For example, if you update a chart component inside a Transition, but then start typing into an input while the chart is in the middle of a re-render, React will restart the rendering work on the chart component after handling the input update.

* Transition updates can't be used to control text inputs.

* If there are multiple ongoing Transitions, React currently batches them together. This is a limitation that may be removed in a future release.

## Usage {/*usage*/}

### Perform non-blocking updates with Actions {/*perform-non-blocking-updates-with-actions*/}

Call `useTransition` at the top of your component to create Actions, and access the pending state:

```js [[1, 4, "isPending"], [2, 4, "startTransition"]]
import {useState, useTransition} from 'react';

function CheckoutForm() {
 const [isPending, startTransition] = useTransition();
 // ...
}
```

`useTransition` returns an array with exactly two items:

1. The <CodeStep step={1}>`isPending` flag</CodeStep> that tells you whether there is a pending Transition.
2. The <CodeStep step={2}>`startTransition` function</CodeStep> that lets you create an Action.

To start a Transition, pass a function to `startTransition` like this:

```js
import {useState, useTransition} from 'react';
import {updateQuantity} from './api';

function CheckoutForm() {
 const [isPending, startTransition] = useTransition();
 const [quantity, setQuantity] = useState(1);

 function onSubmit(newQuantity) {
 startTransition(async function () {
 const savedQuantity = await updateQuantity(newQuantity);
 startTransition(() => {
 setQuantity(savedQuantity);
 });
 });
 }
 // ...
}
```

The function passed to `startTransition` is called the "Action". You can update state and (optionally) perform side effects within an Action, and the work will be done in the background without blocking user interactions on the page. A Transition can include multiple Actions, and while a Transition is in progress, your UI stays responsive. For example, if the user clicks a tab but then changes their mind and clicks another tab, the second click will be immediately handled without waiting for the first update to finish.

To give the user feedback about in-progress Transitions, the `isPending` state switches to `true` at the first call to `startTransition`, and stays `true` until all Actions complete and the final state is shown to the user. Transitions ensure side effects in Actions to complete in order to [prevent unwanted loading indicators](#preventing-unwanted-loading-indicators), and you can provide immediate feedback while the Transition is in progress with `useOptimistic`.

<Recipes titleText="The difference between Actions and regular event handling">

#### Updating the quantity in an Action {/*updating-the-quantity-in-an-action*/}

In this example, the `updateQuantity` function simulates a request to the server to update the item's quantity in the cart. This function is *artificially slowed down* so that it takes at least a second to complete the request.

Update the quantity multiple times quickly. Notice that the pending "Total" state is shown while any requests are in progress, and the "Total" updates only after the final request is complete. Because the update is in an Action, the "quantity" can continue to be updated while the request is in progress.

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "react": "beta",
 "react-dom": "beta"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js src/App.js
import { useState, useTransition } from "react";
import { updateQuantity } from "./api";
import Item from "./Item";
import Total from "./Total";

export default function App({}) {
 const [quantity, setQuantity] = useState(1);
 const [isPending, startTransition] = useTransition();

 const updateQuantityAction = async newQuantity => {
 // To access the pending state of a transition,
 // call startTransition again.
 startTransition(async () => {
 const savedQuantity = await updateQuantity(newQuantity);
 startTransition(() => {
 setQuantity(savedQuantity);
 });
 });
 };

 return (
 <div>
 <h1>Checkout</h1>
 <Item action={updateQuantityAction}/>
 <hr />
 <Total quantity={quantity} isPending={isPending} />
 </div>
 );
}
```

```js src/Item.js
import { startTransition } from "react";

export default function Item({action}) {
 function handleChange(event) {
 // To expose an action prop, await the callback in startTransition.
 startTransition(async () => {
 await action(event.target.value);
 })
 }
 return (
 <div className="item">
 <span>Eras Tour Tickets</span>
 <label htmlFor="name">Quantity: </label>
 <input
 type="number"
 onChange={handleChange}
 defaultValue={1}
 min={1}
 />
 </div>
 )
}
```

```js src/Total.js
const intl = new Intl.NumberFormat("en-US", {
 style: "currency",
 currency: "USD"
});

export default function Total({quantity, isPending}) {
 return (
 <div className="total">
 <span>Total:</span>
 <span>
 {isPending ? "🌀 Updating..." : `${intl.format(quantity * 9999)}`}
 </span>
 </div>
 )
}
```

```js src/api.js
export async function updateQuantity(newQuantity) {
 return new Promise((resolve, reject) => {
 // Simulate a slow network request.
 setTimeout(() => {
 resolve(newQuantity);
 }, 2000);
 });
}
```

```css
.item {
 display: flex;
 align-items: center;
 justify-content: start;
}

.item label {
 flex: 1;
 text-align: right;
}

.item input {
 margin-left: 4px;
 width: 60px;
 padding: 4px;
}

.total {
 height: 50px;
 line-height: 25px;
 display: flex;
 align-content: center;
 justify-content: space-between;
}
```

</Sandpack>

This is a basic example to demonstrate how Actions work, but this example does not handle requests completing out of order. When updating the quantity multiple times, it's possible for the previous requests to finish after later requests causing the quantity to update out of order. This is a known limitation that we will fix in the future (see [Troubleshooting](#my-state-updates-in-transitions-are-out-of-order) below).

For common use cases, React provides built-in abstractions such as:
- [`useActionState`](/reference/react/useActionState)
- [`<form>` actions](/reference/react-dom/components/form)
- [Server Functions](/reference/rsc/server-functions)

These solutions handle request ordering for you. When using Transitions to build your own custom hooks or libraries that manage async state transitions, you have greater control over the request ordering, but you must handle it yourself.

<Solution />

#### Updating the quantity without an Action {/*updating-the-users-name-without-an-action*/}

In this example, the `updateQuantity` function also simulates a request to the server to update the item's quantity in the cart. This function is *artificially slowed down* so that it takes at least a second to complete the request.

Update the quantity multiple times quickly. Notice that the pending "Total" state is shown while any requests is in progress, but the "Total" updates multiple times for each time the "quantity" was clicked:

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "react": "beta",
 "react-dom": "beta"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js src/App.js
import { useState } from "react";
import { updateQuantity } from "./api";
import Item from "./Item";
import Total from "./Total";

export default function App({}) {
 const [quantity, setQuantity] = useState(1);
 const [isPending, setIsPending] = useState(false);

 const onUpdateQuantity = async newQuantity => {
 // Manually set the isPending State.
 setIsPending(true);
 const savedQuantity = await updateQuantity(newQuantity);
 setIsPending(false);
 setQuantity(savedQuantity);
 };

 return (
 <div>
 <h1>Checkout</h1>
 <Item onUpdateQuantity={onUpdateQuantity}/>
 <hr />
 <Total quantity={quantity} isPending={isPending} />
 </div>
 );
}

```

```js src/Item.js
export default function Item({onUpdateQuantity}) {
 function handleChange(event) {
 onUpdateQuantity(event.target.value);
 }
 return (
 <div className="item">
 <span>Eras Tour Tickets</span>
 <label htmlFor="name">Quantity: </label>
 <input
 type="number"
 onChange={handleChange}
 defaultValue={1}
 min={1}
 />
 </div>
 )
}
```

```js src/Total.js
const intl = new Intl.NumberFormat("en-US", {
 style: "currency",
 currency: "USD"
});

export default function Total({quantity, isPending}) {
 return (
 <div className="total">
 <span>Total:</span>
 <span>
 {isPending ? "🌀 Updating..." : `${intl.format(quantity * 9999)}`}
 </span>
 </div>
 )
}
```

```js src/api.js
export async function updateQuantity(newQuantity) {
 return new Promise((resolve, reject) => {
 // Simulate a slow network request.
 setTimeout(() => {
 resolve(newQuantity);
 }, 2000);
 });
}
```

```css
.item {
 display: flex;
 align-items: center;
 justify-content: start;
}

.item label {
 flex: 1;
 text-align: right;
}

.item input {
 margin-left: 4px;
 width: 60px;
 padding: 4px;
}

.total {
 height: 50px;
 line-height: 25px;
 display: flex;
 align-content: center;
 justify-content: space-between;
}
```

</Sandpack>

A common solution to this problem is to prevent the user from making changes while the quantity is updating:

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "react": "beta",
 "react-dom": "beta"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js src/App.js
import { useState, useTransition } from "react";
import { updateQuantity } from "./api";
import Item from "./Item";
import Total from "./Total";

export default function App({}) {
 const [quantity, setQuantity] = useState(1);
 const [isPending, setIsPending] = useState(false);

 const onUpdateQuantity = async event => {
 const newQuantity = event.target.value;
 // Manually set the isPending state.
 setIsPending(true);
 const savedQuantity = await updateQuantity(newQuantity);
 setIsPending(false);
 setQuantity(savedQuantity);
 };

 return (
 <div>
 <h1>Checkout</h1>
 <Item isPending={isPending} onUpdateQuantity={onUpdateQuantity}/>
 <hr />
 <Total quantity={quantity} isPending={isPending} />
 </div>
 );
}

```

```js src/Item.js
export default function Item({isPending, onUpdateQuantity}) {
 return (
 <div className="item">
 <span>Eras Tour Tickets</span>
 <label htmlFor="name">Quantity: </label>
 <input
 type="number"
 disabled={isPending}
 onChange={onUpdateQuantity}
 defaultValue={1}
 min={1}
 />
 </div>
 )
}
```

```js src/Total.js
const intl = new Intl.NumberFormat("en-US", {
 style: "currency",
 currency: "USD"
});

export default function Total({quantity, isPending}) {
 return (
 <div className="total">
 <span>Total:</span>
 <span>
 {isPending ? "🌀 Updating..." : `${intl.format(quantity * 9999)}`}
 </span>
 </div>
 )
}
```

```js src/api.js
export async function updateQuantity(newQuantity) {
 return new Promise((resolve, reject) => {
 // Simulate a slow network request.
 setTimeout(() => {
 resolve(newQuantity);
 }, 2000);
 });
}
```

```css
.item {
 display: flex;
 align-items: center;
 justify-content: start;
}

.item label {
 flex: 1;
 text-align: right;
}

.item input {
 margin-left: 4px;
 width: 60px;
 padding: 4px;
}

.total {
 height: 50px;
 line-height: 25px;
 display: flex;
 align-content: center;
 justify-content: space-between;
}
```

</Sandpack>

This solution makes the app feel slow, because the user must wait each time they update the quantity. It's possible to add more complex handling manually to allow the user to interact with the UI while the quantity is updating, but Actions handle this case with a straight-forward built-in API.

<Solution />

</Recipes>

---

### Exposing `action` prop from components {/*exposing-action-props-from-components*/}

You can expose an `action` prop from a component to allow a parent to call an Action.

For example, this `TabButton` component wraps its `onClick` logic in an `action` prop:

```js {8-12}
export default function TabButton({ action, children, isActive }) {
 const [isPending, startTransition] = useTransition();
 if (isActive) {
 return <b>{children}</b>
 }
 return (
 <button onClick={() => {
 startTransition(async () => {
 // await the action that's passed in.
 // This allows it to be either sync or async.
 await action();
 });
 }}>
 {children}
 </button>
 );
}
```

Because the parent component updates its state inside the `action`, that state update gets marked as a Transition. This means you can click on "Posts" and then immediately click "Contact" and it does not block user interactions:

<Sandpack>

```js
import { useState } from 'react';
import TabButton from './TabButton.js';
import AboutTab from './AboutTab.js';
import PostsTab from './PostsTab.js';
import ContactTab from './ContactTab.js';

export default function TabContainer() {
 const [tab, setTab] = useState('about');
 return (
 <>
 <TabButton
 isActive={tab === 'about'}
 action={() => setTab('about')}
 >
 About
 </TabButton>
 <TabButton
 isActive={tab === 'posts'}
 action={() => setTab('posts')}
 >
 Posts (slow)
 </TabButton>
 <TabButton
 isActive={tab === 'contact'}
 action={() => setTab('contact')}
 >
 Contact
 </TabButton>
 <hr />
 {tab === 'about' && <AboutTab />}
 {tab === 'posts' && <PostsTab />}
 {tab === 'contact' && <ContactTab />}
 </>
 );
}
```

```js src/TabButton.js active
import { useTransition } from 'react';

export default function TabButton({ action, children, isActive }) {
 const [isPending, startTransition] = useTransition();
 if (isActive) {
 return <b>{children}</b>
 }
 if (isPending) {
 return <b className="pending">{children}</b>;
 }
 return (
 <button onClick={async () => {
 startTransition(async () => {
 // await the action that's passed in.
 // This allows it to be either sync or async.
 await action();
 });
 }}>
 {children}
 </button>
 );
}
```

```js src/AboutTab.js
export default function AboutTab() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js {expectedErrors: {'react-compiler': [19, 20]}} src/PostsTab.js
import { memo } from 'react';

const PostsTab = memo(function PostsTab() {
 // Log once. The actual slowdown is inside SlowPost.
 console.log('[ARTIFICIALLY SLOW] Rendering 500 <SlowPost />');

 let items = [];
 for (let i = 0; i < 500; i++) {
 items.push(<SlowPost key={i} index={i} />);
 }
 return (
 <ul className="items">
 {items}
 </ul>
 );
});

function SlowPost({ index }) {
 let startTime = performance.now();
 while (performance.now() - startTime < 1) {
 // Do nothing for 1 ms per item to emulate extremely slow code
 }

 return (
 <li className="item">
 Post #{index + 1}
 </li>
 );
}

export default PostsTab;
```

```js src/ContactTab.js
export default function ContactTab() {
 return (
 <>
 <p>
 You can find me online here:
 </p>
 <ul>
 <li>admin@mysite.com</li>
 <li>+123456789</li>
 </ul>
 </>
 );
}
```

```css
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
.items {
 max-height: 300px;
 overflow: auto;
}
```

</Sandpack>

<Note>

When exposing an `action` prop from a component, you should `await` it inside the transition.

This allows the `action` callback to be either synchronous or asynchronous without requiring an additional `startTransition` to wrap the `await` in the action.

</Note>

---

### Displaying a pending visual state {/*displaying-a-pending-visual-state*/}

You can use the `isPending` boolean value returned by `useTransition` to indicate to the user that a Transition is in progress. For example, the tab button can have a special "pending" visual state:

```js {4-6}
function TabButton({ action, children, isActive }) {
 const [isPending, startTransition] = useTransition();
 // ...
 if (isPending) {
 return <b className="pending">{children}</b>;
 }
 // ...
```

Notice how clicking "Posts" now feels more responsive because the tab button itself updates right away:

<Sandpack>

```js
import { useState } from 'react';
import TabButton from './TabButton.js';
import AboutTab from './AboutTab.js';
import PostsTab from './PostsTab.js';
import ContactTab from './ContactTab.js';

export default function TabContainer() {
 const [tab, setTab] = useState('about');
 return (
 <>
 <TabButton
 isActive={tab === 'about'}
 action={() => setTab('about')}
 >
 About
 </TabButton>
 <TabButton
 isActive={tab === 'posts'}
 action={() => setTab('posts')}
 >
 Posts (slow)
 </TabButton>
 <TabButton
 isActive={tab === 'contact'}
 action={() => setTab('contact')}
 >
 Contact
 </TabButton>
 <hr />
 {tab === 'about' && <AboutTab />}
 {tab === 'posts' && <PostsTab />}
 {tab === 'contact' && <ContactTab />}
 </>
 );
}
```

```js src/TabButton.js active
import { useTransition } from 'react';

export default function TabButton({ action, children, isActive }) {
 const [isPending, startTransition] = useTransition();
 if (isActive) {
 return <b>{children}</b>
 }
 if (isPending) {
 return <b className="pending">{children}</b>;
 }
 return (
 <button onClick={() => {
 startTransition(async () => {
 await action();
 });
 }}>
 {children}
 </button>
 );
}
```

```js src/AboutTab.js
export default function AboutTab() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js {expectedErrors: {'react-compiler': [19, 20]}} src/PostsTab.js
import { memo } from 'react';

const PostsTab = memo(function PostsTab() {
 // Log once. The actual slowdown is inside SlowPost.
 console.log('[ARTIFICIALLY SLOW] Rendering 500 <SlowPost />');

 let items = [];
 for (let i = 0; i < 500; i++) {
 items.push(<SlowPost key={i} index={i} />);
 }
 return (
 <ul className="items">
 {items}
 </ul>
 );
});

function SlowPost({ index }) {
 let startTime = performance.now();
 while (performance.now() - startTime < 1) {
 // Do nothing for 1 ms per item to emulate extremely slow code
 }

 return (
 <li className="item">
 Post #{index + 1}
 </li>
 );
}

export default PostsTab;
```

```js src/ContactTab.js
export default function ContactTab() {
 return (
 <>
 <p>
 You can find me online here:
 </p>
 <ul>
 <li>admin@mysite.com</li>
 <li>+123456789</li>
 </ul>
 </>
 );
}
```

```css
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
.items {
 max-height: 300px;
 overflow: auto;
}
```

</Sandpack>

---

### Preventing unwanted loading indicators {/*preventing-unwanted-loading-indicators*/}

In this example, the `PostsTab` component fetches some data using [use](/reference/react/use). When you click the "Posts" tab, the `PostsTab` component *suspends*, causing the closest loading fallback to appear:

<Sandpack>

```js
import { Suspense, useState } from 'react';
import TabButton from './TabButton.js';
import AboutTab from './AboutTab.js';
import PostsTab from './PostsTab.js';
import ContactTab from './ContactTab.js';

export default function TabContainer() {
 const [tab, setTab] = useState('about');
 return (
 <Suspense fallback={<h1>🌀 Loading...</h1>}>
 <TabButton
 isActive={tab === 'about'}
 action={() => setTab('about')}
 >
 About
 </TabButton>
 <TabButton
 isActive={tab === 'posts'}
 action={() => setTab('posts')}
 >
 Posts
 </TabButton>
 <TabButton
 isActive={tab === 'contact'}
 action={() => setTab('contact')}
 >
 Contact
 </TabButton>
 <hr />
 {tab === 'about' && <AboutTab />}
 {tab === 'posts' && <PostsTab />}
 {tab === 'contact' && <ContactTab />}
 </Suspense>
 );
}
```

```js src/TabButton.js
export default function TabButton({ action, children, isActive }) {
 if (isActive) {
 return <b>{children}</b>
 }
 return (
 <button onClick={() => {
 action();
 }}>
 {children}
 </button>
 );
}
```

```js src/AboutTab.js hidden
export default function AboutTab() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/PostsTab.js hidden
import {use} from 'react';
import { fetchData } from './data.js';

function PostsTab() {
 const posts = use(fetchData('/posts'));
 return (
 <ul className="items">
 {posts.map(post =>
 <Post key={post.id} title={post.title} />
 )}
 </ul>
 );
}

function Post({ title }) {
 return (
 <li className="item">
 {title}
 </li>
 );
}

export default PostsTab;
```

```js src/ContactTab.js hidden
export default function ContactTab() {
 return (
 <>
 <p>
 You can find me online here:
 </p>
 <ul>
 <li>admin@mysite.com</li>
 <li>+123456789</li>
 </ul>
 </>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/posts')) {
 return await getPosts();
 } else {
 throw Error('Not implemented');
 }
}

async function getPosts() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });
 let posts = [];
 for (let i = 0; i < 500; i++) {
 posts.push({
 id: i,
 title: 'Post #' + (i + 1)
 });
 }
 return posts;
}
```

```css
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
```

</Sandpack>

Hiding the entire tab container to show a loading indicator leads to a jarring user experience. If you add `useTransition` to `TabButton`, you can instead display the pending state in the tab button instead.

Notice that clicking "Posts" no longer replaces the entire tab container with a spinner:

<Sandpack>

```js
import { Suspense, useState } from 'react';
import TabButton from './TabButton.js';
import AboutTab from './AboutTab.js';
import PostsTab from './PostsTab.js';
import ContactTab from './ContactTab.js';

export default function TabContainer() {
 const [tab, setTab] = useState('about');
 return (
 <Suspense fallback={<h1>🌀 Loading...</h1>}>
 <TabButton
 isActive={tab === 'about'}
 action={() => setTab('about')}
 >
 About
 </TabButton>
 <TabButton
 isActive={tab === 'posts'}
 action={() => setTab('posts')}
 >
 Posts
 </TabButton>
 <TabButton
 isActive={tab === 'contact'}
 action={() => setTab('contact')}
 >
 Contact
 </TabButton>
 <hr />
 {tab === 'about' && <AboutTab />}
 {tab === 'posts' && <PostsTab />}
 {tab === 'contact' && <ContactTab />}
 </Suspense>
 );
}
```

```js src/TabButton.js active
import { useTransition } from 'react';

export default function TabButton({ action, children, isActive }) {
 const [isPending, startTransition] = useTransition();
 if (isActive) {
 return <b>{children}</b>
 }
 if (isPending) {
 return <b className="pending">{children}</b>;
 }
 return (
 <button onClick={() => {
 startTransition(async () => {
 await action();
 });
 }}>
 {children}
 </button>
 );
}
```

```js src/AboutTab.js hidden
export default function AboutTab() {
 return (
 <p>Welcome to my profile!</p>
 );
}
```

```js src/PostsTab.js hidden
import {use} from 'react';
import { fetchData } from './data.js';

function PostsTab() {
 const posts = use(fetchData('/posts'));
 return (
 <ul className="items">
 {posts.map(post =>
 <Post key={post.id} title={post.title} />
 )}
 </ul>
 );
}

function Post({ title }) {
 return (
 <li className="item">
 {title}
 </li>
 );
}

export default PostsTab;
```

```js src/ContactTab.js hidden
export default function ContactTab() {
 return (
 <>
 <p>
 You can find me online here:
 </p>
 <ul>
 <li>admin@mysite.com</li>
 <li>+123456789</li>
 </ul>
 </>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url.startsWith('/posts')) {
 return await getPosts();
 } else {
 throw Error('Not implemented');
 }
}

async function getPosts() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 1000);
 });
 let posts = [];
 for (let i = 0; i < 500; i++) {
 posts.push({
 id: i,
 title: 'Post #' + (i + 1)
 });
 }
 return posts;
}
```

```css
button { margin-right: 10px }
b { display: inline-block; margin-right: 10px; }
.pending { color: #777; }
```

</Sandpack>

[Read more about using Transitions with Suspense.](/reference/react/Suspense#preventing-already-revealed-content-from-hiding)

<Note>

Transitions only "wait" long enough to avoid hiding *already revealed* content (like the tab container). If the Posts tab had a [nested `<Suspense>` boundary,](/reference/react/Suspense#revealing-nested-content-as-it-loads) the Transition would not "wait" for it.

</Note>

---

### Building a Suspense-enabled router {/*building-a-suspense-enabled-router*/}

If you're building a React framework or a router, we recommend marking page navigations as Transitions.

```js {3,6,8}
function Router() {
 const [page, setPage] = useState('/');
 const [isPending, startTransition] = useTransition();

 function navigate(url) {
 startTransition(() => {
 setPage(url);
 });
 }
 // ...
```

This is recommended for three reasons:

- [Transitions are interruptible,](#perform-non-blocking-updates-with-actions) which lets the user click away without waiting for the re-render to complete.
- [Transitions prevent unwanted loading indicators,](#preventing-unwanted-loading-indicators) which lets the user avoid jarring jumps on navigation.
- [Transitions wait for all pending actions](#perform-non-blocking-updates-with-actions) which lets the user wait for side effects to complete before the new page is shown.

Here is a simplified router example using Transitions for navigations.

<Sandpack>

```js src/App.js
import { Suspense, useState, useTransition } from 'react';
import IndexPage from './IndexPage.js';
import ArtistPage from './ArtistPage.js';
import Layout from './Layout.js';

export default function App() {
 return (
 <Suspense fallback={<BigSpinner />}>
 <Router />
 </Suspense>
 );
}

function Router() {
 const [page, setPage] = useState('/');
 const [isPending, startTransition] = useTransition();

 function navigate(url) {
 startTransition(() => {
 setPage(url);
 });
 }

 let content;
 if (page === '/') {
 content = (
 <IndexPage navigate={navigate} />
 );
 } else if (page === '/the-beatles') {
 content = (
 <ArtistPage
 artist={{
 id: 'the-beatles',
 name: 'The Beatles',
 }}
 />
 );
 }
 return (
 <Layout isPending={isPending}>
 {content}
 </Layout>
 );
}

function BigSpinner() {
 return <h2>🌀 Loading...</h2>;
}
```

```js src/Layout.js
export default function Layout({ children, isPending }) {
 return (
 <div className="layout">
 <section className="header" style={{
 opacity: isPending ? 0.7 : 1
 }}>
 Music Browser
 </section>
 <main>
 {children}
 </main>
 </div>
 );
}
```

```js src/IndexPage.js
export default function IndexPage({ navigate }) {
 return (
 <button onClick={() => navigate('/the-beatles')}>
 Open The Beatles artist page
 </button>
 );
}
```

```js src/ArtistPage.js
import { Suspense } from 'react';
import Albums from './Albums.js';
import Biography from './Biography.js';
import Panel from './Panel.js';

export default function ArtistPage({ artist }) {
 return (
 <>
 <h1>{artist.name}</h1>
 <Biography artistId={artist.id} />
 <Suspense fallback={<AlbumsGlimmer />}>
 <Panel>
 <Albums artistId={artist.id} />
 </Panel>
 </Suspense>
 </>
 );
}

function AlbumsGlimmer() {
 return (
 <div className="glimmer-panel">
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 <div className="glimmer-line" />
 </div>
 );
}
```

```js src/Albums.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Albums({ artistId }) {
 const albums = use(fetchData(`/${artistId}/albums`));
 return (
 <ul>
 {albums.map(album => (
 <li key={album.id}>
 {album.title} ({album.year})
 </li>
 ))}
 </ul>
 );
}
```

```js src/Biography.js
import {use} from 'react';
import { fetchData } from './data.js';

export default function Biography({ artistId }) {
 const bio = use(fetchData(`/${artistId}/bio`));
 return (
 <section>
 <p className="bio">{bio}</p>
 </section>
 );
}
```

```js src/Panel.js
export default function Panel({ children }) {
 return (
 <section className="panel">
 {children}
 </section>
 );
}
```

```js src/data.js hidden
// Note: the way you would do data fetching depends on
// the framework that you use together with Suspense.
// Normally, the caching logic would be inside a framework.

let cache = new Map();

export function fetchData(url) {
 if (!cache.has(url)) {
 cache.set(url, getData(url));
 }
 return cache.get(url);
}

async function getData(url) {
 if (url === '/the-beatles/albums') {
 return await getAlbums();
 } else if (url === '/the-beatles/bio') {
 return await getBio();
 } else {
 throw Error('Not implemented');
 }
}

async function getBio() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 500);
 });

 return `The Beatles were an English rock band,
 formed in Liverpool in 1960, that comprised
 John Lennon, Paul McCartney, George Harrison
 and Ringo Starr.`;
}

async function getAlbums() {
 // Add a fake delay to make waiting noticeable.
 await new Promise(resolve => {
 setTimeout(resolve, 3000);
 });

 return [{
 id: 13,
 title: 'Let It Be',
 year: 1970
 }, {
 id: 12,
 title: 'Abbey Road',
 year: 1969
 }, {
 id: 11,
 title: 'Yellow Submarine',
 year: 1969
 }, {
 id: 10,
 title: 'The Beatles',
 year: 1968
 }, {
 id: 9,
 title: 'Magical Mystery Tour',
 year: 1967
 }, {
 id: 8,
 title: 'Sgt. Pepper\'s Lonely Hearts Club Band',
 year: 1967
 }, {
 id: 7,
 title: 'Revolver',
 year: 1966
 }, {
 id: 6,
 title: 'Rubber Soul',
 year: 1965
 }, {
 id: 5,
 title: 'Help!',
 year: 1965
 }, {
 id: 4,
 title: 'Beatles For Sale',
 year: 1964
 }, {
 id: 3,
 title: 'A Hard Day\'s Night',
 year: 1964
 }, {
 id: 2,
 title: 'With The Beatles',
 year: 1963
 }, {
 id: 1,
 title: 'Please Please Me',
 year: 1963
 }];
}
```

```css
main {
 min-height: 200px;
 padding: 10px;
}

.layout {
 border: 1px solid black;
}

.header {
 background: #222;
 padding: 10px;
 text-align: center;
 color: white;
}

.bio { font-style: italic; }

.panel {
 border: 1px solid #aaa;
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-panel {
 border: 1px dashed #aaa;
 background: linear-gradient(90deg, rgba(221,221,221,1) 0%, rgba(255,255,255,1) 100%);
 border-radius: 6px;
 margin-top: 20px;
 padding: 10px;
}

.glimmer-line {
 display: block;
 width: 60%;
 height: 20px;
 margin: 10px;
 border-radius: 4px;
 background: #f0f0f0;
}
```

</Sandpack>

<Note>

[Suspense-enabled](/reference/react/Suspense) routers are expected to wrap the navigation updates into Transitions by default.

</Note>

---

### Displaying an error to users with an error boundary {/*displaying-an-error-to-users-with-error-boundary*/}

If a function passed to `startTransition` throws an error, you can display an error to your user with an [error boundary](/reference/react/Component#catching-rendering-errors-with-an-error-boundary). To use an error boundary, wrap the component where you are calling the `useTransition` in an error boundary. Once the function passed to `startTransition` errors, the fallback for the error boundary will be displayed.

<Sandpack>

```js src/AddCommentContainer.js active
import { useTransition } from "react";
import { ErrorBoundary } from "react-error-boundary";

export function AddCommentContainer() {
 return (
 <ErrorBoundary fallback={<p>⚠️Something went wrong</p>}>
 <AddCommentButton />
 </ErrorBoundary>
 );
}

function addComment(comment) {
 // For demonstration purposes to show Error Boundary
 if (comment == null) {
 throw new Error("Example Error: An error thrown to trigger error boundary");
 }
}

function AddCommentButton() {
 const [pending, startTransition] = useTransition();

 return (
 <button
 disabled={pending}
 onClick={() => {
 startTransition(() => {
 // Intentionally not passing a comment
 // so error gets thrown
 addComment();
 });
 }}
 >
 Add comment
 </button>
 );
}
```

```js src/App.js hidden
import { AddCommentContainer } from "./AddCommentContainer.js";

export default function App() {
 return <AddCommentContainer />;
}
```

```js src/index.js hidden
import React, { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './styles.css';
import App from './App';

const root = createRoot(document.getElementById('root'));
root.render(
 <StrictMode>
 <App />
 </StrictMode>
);
```

```json package.json hidden
{
 "dependencies": {
 "react": "19.0.0-rc-3edc000d-20240926",
 "react-dom": "19.0.0-rc-3edc000d-20240926",
 "react-scripts": "^5.0.0",
 "react-error-boundary": "4.0.3"
 },
 "main": "/index.js"
}
```
</Sandpack>

---

## Troubleshooting {/*troubleshooting*/}

### Updating an input in a Transition doesn't work {/*updating-an-input-in-a-transition-doesnt-work*/}

You can't use a Transition for a state variable that controls an input:

```js {4,10}
const [text, setText] = useState('');
// ...
function handleChange(e) {
 // ❌ Can't use Transitions for controlled input state
 startTransition(() => {
 setText(e.target.value);
 });
}
// ...
return <input value={text} onChange={handleChange} />;
```

This is because Transitions are non-blocking, but updating an input in response to the change event should happen synchronously. If you want to run a Transition in response to typing, you have two options:

1. You can declare two separate state variables: one for the input state (which always updates synchronously), and one that you will update in a Transition. This lets you control the input using the synchronous state, and pass the Transition state variable (which will "lag behind" the input) to the rest of your rendering logic.
2. Alternatively, you can have one state variable, and add [`useDeferredValue`](/reference/react/useDeferredValue) which will "lag behind" the real value. It will trigger non-blocking re-renders to "catch up" with the new value automatically.

---

### React doesn't treat my state update as a Transition {/*react-doesnt-treat-my-state-update-as-a-transition*/}

When you wrap a state update in a Transition, make sure that it happens *during* the `startTransition` call:

```js
startTransition(() => {
 // ✅ Setting state *during* startTransition call
 setPage('/about');
});
```

The function you pass to `startTransition` must be synchronous. You can't mark an update as a Transition like this:

```js
startTransition(() => {
 // ❌ Setting state *after* startTransition call
 setTimeout(() => {
 setPage('/about');
 }, 1000);
});
```

Instead, you could do this:

```js
setTimeout(() => {
 startTransition(() => {
 // ✅ Setting state *during* startTransition call
 setPage('/about');
 });
}, 1000);
```

---

### React doesn't treat my state update after `await` as a Transition {/*react-doesnt-treat-my-state-update-after-await-as-a-transition*/}

When you use `await` inside a `startTransition` function, the state updates that happen after the `await` are not marked as Transitions. You must wrap state updates after each `await` in a `startTransition` call:

```js
startTransition(async () => {
 await someAsyncFunction();
 // ❌ Not using startTransition after await
 setPage('/about');
});
```

However, this works instead:

```js
startTransition(async () => {
 await someAsyncFunction();
 // ✅ Using startTransition *after* await
 startTransition(() => {
 setPage('/about');
 });
});
```

This is a JavaScript limitation due to React losing the scope of the async context. In the future, when [AsyncContext](https://github.com/tc39/proposal-async-context) is available, this limitation will be removed.

---

### I want to call `useTransition` from outside a component {/*i-want-to-call-usetransition-from-outside-a-component*/}

You can't call `useTransition` outside a component because it's a Hook. In this case, use the standalone [`startTransition`](/reference/react/startTransition) method instead. It works the same way, but it doesn't provide the `isPending` indicator.

---

### The function I pass to `startTransition` executes immediately {/*the-function-i-pass-to-starttransition-executes-immediately*/}

If you run this code, it will print 1, 2, 3:

```js {1,3,6}
console.log(1);
startTransition(() => {
 console.log(2);
 setPage('/about');
});
console.log(3);
```

**It is expected to print 1, 2, 3.** The function you pass to `startTransition` does not get delayed. Unlike with the browser `setTimeout`, it does not run the callback later. React executes your function immediately, but any state updates scheduled *while it is running* are marked as Transitions. You can imagine that it works like this:

```js
// A simplified version of how React works

let isInsideTransition = false;

function startTransition(scope) {
 isInsideTransition = true;
 scope();
 isInsideTransition = false;
}

function setState() {
 if (isInsideTransition) {
 // ... schedule a Transition state update ...
 } else {
 // ... schedule an urgent state update ...
 }
}
```

### My state updates in Transitions are out of order {/*my-state-updates-in-transitions-are-out-of-order*/}

If you `await` inside `startTransition`, you might see the updates happen out of order.

In this example, the `updateQuantity` function simulates a request to the server to update the item's quantity in the cart. This function *artificially returns every other request after the previous* to simulate race conditions for network requests.

Try updating the quantity once, then update it quickly multiple times. You might see the incorrect total:

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "react": "beta",
 "react-dom": "beta"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js src/App.js
import { useState, useTransition } from "react";
import { updateQuantity } from "./api";
import Item from "./Item";
import Total from "./Total";

export default function App({}) {
 const [quantity, setQuantity] = useState(1);
 const [isPending, startTransition] = useTransition();
 // Store the actual quantity in separate state to show the mismatch.
 const [clientQuantity, setClientQuantity] = useState(1);

 const updateQuantityAction = newQuantity => {
 setClientQuantity(newQuantity);

 // Access the pending state of the transition,
 // by wrapping in startTransition again.
 startTransition(async () => {
 const savedQuantity = await updateQuantity(newQuantity);
 startTransition(() => {
 setQuantity(savedQuantity);
 });
 });
 };

 return (
 <div>
 <h1>Checkout</h1>
 <Item action={updateQuantityAction}/>
 <hr />
 <Total clientQuantity={clientQuantity} savedQuantity={quantity} isPending={isPending} />
 </div>
 );
}

```

```js src/Item.js
import {startTransition} from 'react';

export default function Item({action}) {
 function handleChange(e) {
 // Update the quantity in an Action.
 startTransition(async () => {
 await action(e.target.value);
 });
 }
 return (
 <div className="item">
 <span>Eras Tour Tickets</span>
 <label htmlFor="name">Quantity: </label>
 <input
 type="number"
 onChange={handleChange}
 defaultValue={1}
 min={1}
 />
 </div>
 )
}
```

```js src/Total.js
const intl = new Intl.NumberFormat("en-US", {
 style: "currency",
 currency: "USD"
});

export default function Total({ clientQuantity, savedQuantity, isPending }) {
 return (
 <div className="total">
 <span>Total:</span>
 <div>
 <div>
 {isPending
 ? "🌀 Updating..."
 : `${intl.format(savedQuantity * 9999)}`}
 </div>
 <div className="error">
 {!isPending &&
 clientQuantity !== savedQuantity &&
 `Wrong total, expected: ${intl.format(clientQuantity * 9999)}`}
 </div>
 </div>
 </div>
 );
}
```

```js src/api.js
let firstRequest = true;
export async function updateQuantity(newName) {
 return new Promise((resolve, reject) => {
 if (firstRequest === true) {
 firstRequest = false;
 setTimeout(() => {
 firstRequest = true;
 resolve(newName);
 // Simulate every other request being slower
 }, 1000);
 } else {
 setTimeout(() => {
 resolve(newName);
 }, 50);
 }
 });
}
```

```css
.item {
 display: flex;
 align-items: center;
 justify-content: start;
}

.item label {
 flex: 1;
 text-align: right;
}

.item input {
 margin-left: 4px;
 width: 60px;
 padding: 4px;
}

.total {
 height: 50px;
 line-height: 25px;
 display: flex;
 align-content: center;
 justify-content: space-between;
}

.total div {
 display: flex;
 flex-direction: column;
 align-items: flex-end;
}

.error {
 color: red;
}
```

</Sandpack>

When clicking multiple times, it's possible for previous requests to finish after later requests. When this happens, React currently has no way to know the intended order. This is because the updates are scheduled asynchronously, and React loses context of the order across the async boundary.

This is expected, because Actions within a Transition do not guarantee execution order. For common use cases, React provides higher-level abstractions like [`useActionState`](/reference/react/useActionState) and [`<form>` actions](/reference/react-dom/components/form) that handle ordering for you. For advanced use cases, you'll need to implement your own queuing and abort logic to handle this.

Example of `useActionState` handling execution order:

<Sandpack>

```json package.json hidden
{
 "dependencies": {
 "react": "beta",
 "react-dom": "beta"
 },
 "scripts": {
 "start": "react-scripts start",
 "build": "react-scripts build",
 "test": "react-scripts test --env=jsdom",
 "eject": "react-scripts eject"
 }
}
```

```js src/App.js
import { useState, useActionState } from "react";
import { updateQuantity } from "./api";
import Item from "./Item";
import Total from "./Total";

export default function App({}) {
 // Store the actual quantity in separate state to show the mismatch.
 const [clientQuantity, setClientQuantity] = useState(1);
 const [quantity, updateQuantityAction, isPending] = useActionState(
 async (prevState, payload) => {
 setClientQuantity(payload);
 const savedQuantity = await updateQuantity(payload);
 return savedQuantity; // Return the new quantity to update the state
 },
 1 // Initial quantity
 );

 return (
 <div>
 <h1>Checkout</h1>
 <Item action={updateQuantityAction}/>
 <hr />
 <Total clientQuantity={clientQuantity} savedQuantity={quantity} isPending={isPending} />
 </div>
 );
}

```

```js src/Item.js
import {startTransition} from 'react';

export default function Item({action}) {
 function handleChange(e) {
 // Update the quantity in an Action.
 startTransition(() => {
 action(e.target.value);
 });
 }
 return (
 <div className="item">
 <span>Eras Tour Tickets</span>
 <label htmlFor="name">Quantity: </label>
 <input
 type="number"
 onChange={handleChange}
 defaultValue={1}
 min={1}
 />
 </div>
 )
}
```

```js src/Total.js
const intl = new Intl.NumberFormat("en-US", {
 style: "currency",
 currency: "USD"
});

export default function Total({ clientQuantity, savedQuantity, isPending }) {
 return (
 <div className="total">
 <span>Total:</span>
 <div>
 <div>
 {isPending
 ? "🌀 Updating..."
 : `${intl.format(savedQuantity * 9999)}`}
 </div>
 <div className="error">
 {!isPending &&
 clientQuantity !== savedQuantity &&
 `Wrong total, expected: ${intl.format(clientQuantity * 9999)}`}
 </div>
 </div>
 </div>
 );
}
```

```js src/api.js
let firstRequest = true;
export async function updateQuantity(newName) {
 return new Promise((resolve, reject) => {
 if (firstRequest === true) {
 firstRequest = false;
 setTimeout(() => {
 firstRequest = true;
 resolve(newName);
 // Simulate every other request being slower
 }, 1000);
 } else {
 setTimeout(() => {
 resolve(newName);
 }, 50);
 }
 });
}
```

```css
.item {
 display: flex;
 align-items: center;
 justify-content: start;
}

.item label {
 flex: 1;
 text-align: right;
}

.item input {
 margin-left: 4px;
 width: 60px;
 padding: 4px;
}

.total {
 height: 50px;
 line-height: 25px;
 display: flex;
 align-content: center;
 justify-content: space-between;
}

.total div {
 display: flex;
 flex-direction: column;
 align-items: flex-end;
}

.error {
 color: red;
}
```

</Sandpack>
