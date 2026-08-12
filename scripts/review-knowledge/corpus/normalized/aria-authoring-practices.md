-

 ##

 Accordion (Sections With Show/Hide Functionality)

 An accordion is a vertically stacked set of interactive headings that each contain a title, content snippet, or thumbnail representing a section of content.

-

 ##

 Alert

 An alert is an element that displays a brief, important message in a way that attracts the user's attention without interrupting the user's task.

-

 ##

 Alert and Message Dialogs

 An alert dialog is a modal dialog that interrupts the user's workflow to communicate an important message and acquire a response.

-

 ##

 Breadcrumb

 A breadcrumb trail consists of a list of links to the parent pages of the current page in hierarchical order.

-

 ##

 Button

 A button is a widget that enables users to trigger an action or event, such as submitting a form, opening a dialog, canceling an action, or performing a delete operation.

-

 ##

 Carousel (Slide Show or Image Rotator)

 A carousel presents a set of items, referred to as slides, by sequentially displaying a subset of one or more slides.

-

 ##

 Checkbox

 WAI-ARIA supports two types of checkbox widgets: dual-state checkboxes toggle between two choices -- checked and not checked -- and tri-state checkboxes, which allow an additional third state known as partially checked.

-

 ##

 Combobox

 A combobox is an input widget that has an associated popup.

-

 ##

 Dialog (Modal)

 A dialog is a window overlaid on either the primary window or another dialog window.

-

 ##

 Disclosure (Show/Hide)

 A disclosure is a widget that enables content to be either collapsed (hidden) or expanded (visible).

-

 ##

 Feed

 A feed is a section of a page that automatically loads new sections of content as the user scrolls.

-

 ##

 Grid (Interactive Tabular Data and Layout Containers)

 A grid widget is a container that enables users to navigate the information or interactive elements it contains using directional navigation keys, such as arrow keys, Home, and End.

-

 ##

 Landmarks

 Landmarks are a set of eight roles that identify the major sections of a page.

-

 ##

 Link

 A link widget provides an interactive reference to a resource.

-

 ##

 Listbox

 A listbox widget presents a list of options and allows a user to select one or more of them.

-

 ##

 Menu and Menubar

 A menu is a widget that offers a list of choices to the user, such as a set of actions or functions.

-

 ##

 Menu Button

 A menu button is a button that opens a menu as described in the Menu and Menubar Pattern.

-

 ##

 Meter

 A meter is a graphical display of a numeric value that varies within a defined range.

-

 ##

 Radio Group

 A radio group is a set of checkable buttons, known as radio buttons, where no more than one of the buttons can be checked at a time.

-

 ##

 Slider

 A slider is an input where the user selects a value from within a given range.

-

 ##

 Slider (Multi-Thumb)

 A multi-thumb slider implements the Slider Pattern but includes two or more thumbs, often on a single rail.

-

 ##

 Spinbutton

 A spinbutton is an input widget that restricts its value to a set or range of discrete values.

-

 ##

 Switch

 A switch is an input widget that allows users to choose one of two values: on or off.

-

 ##

 Table

 Like an HTML table element, a WAI-ARIA table is a static tabular structure containing one or more rows that each contain one or more cells; it is not an interactive widget.

-

 ##

 Tabs

 Tabs are a set of layered sections of content, known as tab panels, that display one panel of content at a time.

-

 ##

 Toolbar

 A toolbar is a container for grouping a set of controls, such as buttons, menubuttons, or checkboxes.

-

 ##

 Tooltip

 A tooltip is a popup that displays information related to an element when the element receives keyboard focus or the mouse hovers over it.

-

 ##

 Tree View

 A tree view widget presents a hierarchical list.

-

 ##

 Treegrid

 A treegrid widget presents a hierarchical data grid consisting of tabular information that is editable or interactive.

-

 ##

 Window Splitter

 A window splitter is a moveable separator between two sections, or panes, of a window that enables users to change the relative size of the panes.

-
 ## No results found.

# Accordion Pattern (Sections With Show/Hide Functionality)

## About This Pattern

An accordion is a vertically stacked set of interactive headings that each contain a title, content snippet, or thumbnail representing a section of content. The headings function as controls that enable users to reveal or hide their associated sections of content. Accordions are commonly used to reduce the need to scroll when presenting multiple sections of content on a single page.

Terms for understanding accordions include:

- Accordion Header:
- Label for or thumbnail representing a section of content that also serves as a control for showing, and in some implementations, hiding the section of content.
- Accordion Panel:
- Section of content associated with an accordion header.

In some accordions, there are additional elements that are always visible adjacent to the accordion header. For instance, a menubutton may accompany each accordion header to provide access to actions that apply to that section. And, in some cases, a snippet of the hidden content may also be visually persistent.

## Example

Accordion Example: demonstrates using an accordion to divide a form into three expandable sections.

## Keyboard Interaction

-
 `Enter`or`Space`:- When focus is on the accordion header for a collapsed panel, expands the associated panel. If the implementation allows only one panel to be expanded, and if another panel is expanded, collapses that panel.
- When focus is on the accordion header for an expanded panel, collapses the panel if the implementation supports collapsing. Some implementations require one panel to be expanded at all times and allow only one panel to be expanded; so, they do not support a collapse function.

- `Tab`: Moves focus to the next focusable element; all focusable elements in the accordion are included in the page- `Tab`sequence.
- `Shift`+- `Tab`: Moves focus to the previous focusable element; all focusable elements in the accordion are included in the page- `Tab`sequence.

## WAI-ARIA Roles, States, and Properties

- The title of each accordion header is contained in an element with role button.
-
 Each accordion header `button`is wrapped in an element with role heading that has a value set for aria-level that is appropriate for the information architecture of the page.- If the native host language has an element with an implicit `heading`and`aria-level`, such as an HTML heading tag, a native host language element may be used.
- The `button`element is the only element inside the`heading`element. That is, if there are other visually persistent elements, they are not included inside the`heading`element.

- If the native host language has an element with an implicit
- If the accordion panel associated with an accordion header is visible, the header `button`element has aria-expanded set to`true`. If the panel is not visible, aria-expanded is set to`false`.
- The accordion header `button`element has aria-controls set to the ID of the element containing the accordion panel content.
- If the accordion panel associated with an accordion header is visible, and if the accordion does not permit the panel to be collapsed, the header `button`element has aria-disabled set to`true`.
-
 Optionally, each element that serves as a container for panel content has role region and aria-labelledby with a value that refers to the button that controls display of the panel.
 - Avoid using the `region`role in circumstances that create landmark region proliferation, e.g., in an accordion that contains more than approximately 6 panels that can be expanded at the same time.
- Role `region`is especially helpful to the perception of structure by screen reader users when panels contain heading elements or a nested accordion.

- Avoid using the

# Alert Pattern

## About This Pattern

An alert is an element that displays a brief, important message in a way that attracts the user's attention without interrupting the user's task. Dynamically rendered alerts are automatically announced by most screen readers, and in some operating systems, they may trigger an alert sound. It is important to note that, at this time, screen readers do not inform users of alerts that are present on the page before page load completes.

Because alerts are intended to provide important and potentially time-sensitive information without interfering with the user's ability to continue working, it is crucial they do not affect keyboard focus. The Alert Dialog Pattern is designed for situations where interrupting work flow is necessary.

It is also important to avoid designing alerts that disappear automatically. An alert that disappears too quickly can lead to failure to meet WCAG 2.0 success criterion 2.2.3. Another critical design consideration is the frequency of interruption caused by alerts. Frequent interruptions inhibit usability for people with visual and cognitive disabilities, which makes meeting the requirements of WCAG 2.0 success criterion 2.2.4 more difficult.

## Example

## Keyboard Interaction

Not applicable.

## WAI-ARIA Roles, States, and Properties

The widget has a role of alert.

# Alert and Message Dialogs Pattern

## About This Pattern

An alert dialog is a modal dialog that interrupts the user's workflow to communicate an important message and acquire a response. Examples include action confirmation prompts and error message confirmations. The alertdialog role enables assistive technologies and browsers to distinguish alert dialogs from other dialogs so they have the option of giving alert dialogs special treatment, such as playing a system alert sound.

## Example

Alert Dialog Example: A confirmation prompt that demonstrates an alert dialog.

## Keyboard Interaction

See the keyboard interaction section for the modal dialog pattern.

## WAI-ARIA Roles, States, and Properties

- The element that contains all elements of the dialog, including the alert message and any dialog buttons, has role alertdialog.
- The dialog container element has aria-modal set to `true`.
-
 The element with role `alertdialog`has either:- A value for aria-labelledby that refers to the element containing the title of the dialog if the dialog has a visible label.
- A value for aria-label if the dialog does not have a visible label.

-
 The element with role `alertdialog`has a value set for aria-describedby that refers to the element containing the alert message.

### Note

-
 Because marking a dialog modal by setting aria-modal to `true`can prevent users of some assistive technologies from perceiving content outside the dialog, users of those technologies will experience severe negative ramifications if a dialog is marked modal but does not behave as a modal for other users. So, mark a dialog modal**only when both:**- Application code prevents all users from interacting in any way with content outside of it.
- Visual styling obscures the content outside of it.

-
 The `aria-modal`property introduced by ARIA 1.1 replaces aria-hidden for informing assistive technologies that content outside a dialog is inert. However, in legacy dialog implementations where`aria-hidden`is used to make content outside a dialog inert for assistive technology users, it is important that:- `aria-hidden`is set to- `true`on each element containing a portion of the inert layer.
- The dialog element is not a descendant of any element that has `aria-hidden`set to`true`.

# Breadcrumb Pattern

## About This Pattern

A breadcrumb trail consists of a list of links to the parent pages of the current page in hierarchical order. It helps users find their place within a website or web application. Breadcrumbs are often placed horizontally before a page's main content.

## Example

## Keyboard Interaction

Not applicable.

## WAI-ARIA Roles, States, and Properties

- Breadcrumb trail is contained within a navigation landmark region.
- The landmark region is labelled via aria-label or aria-labelledby.
-
 The link to the current page has aria-current set to `page`. If the element representing the current page is not a link, aria-current is optional.

# Button Pattern

## About This Pattern

A button is a widget that enables users to trigger an action or event, such as submitting a form, opening a dialog, canceling an action, or performing a delete operation. A common convention for informing users that a button launches a dialog is to append "…" (ellipsis) to the button label, e.g., "Save as…".

In addition to the ordinary button widget, WAI-ARIA supports 2 other types of buttons:

-
 Toggle button: A two-state button that can be either off (not pressed) or on (pressed).
 To tell assistive technologies that a button is a toggle button, specify a value for the attribute aria-pressed.
 For example, a button labelled mute in an audio player could indicate that sound is muted by setting the pressed state true.
 **Important:**it is critical the label on a toggle does not change when its state changes. In this example, when the pressed state is true, the label remains "Mute" so a screen reader would say something like "Mute toggle button pressed". Alternatively, if the design were to call for the button label to change from "Mute" to "Unmute," the aria-pressed attribute would not be needed.
- Menu button: as described in the menu button pattern, a button is revealed to assistive technologies as a menu button if it has the property aria-haspopup set to either `menu`or`true`.

### Note

 The types of actions performed by buttons are distinctly different from the function of a link (see link pattern).
 It is important that both the appearance and role of a widget match the function it provides.
 Nevertheless, elements occasionally have the visual style of a link but perform the action of a button.
 In such cases, giving the element role `button` helps assistive technology users understand the function of the element.
 However, a better solution is to adjust the visual design so it matches the function and ARIA role.

## Examples

- Button Examples: Examples of clickable HTML `div`and`span`elements made into accessible command and toggle buttons.
-
 Button Examples (IDL): Examples of clickable HTML `div`and`span`elements made into accessible command and toggle buttons. This example uses the IDL Interface.

## Keyboard Interaction

When the button has focus:

- `Space`: Activates the button.
- `Enter`: Activates the button.
-
 Following button activation, focus is set depending on the type of action the button performs. For example:
 - If activating the button opens a dialog, the focus moves inside the dialog. (see dialog pattern)
- If activating the button closes a dialog, focus typically returns to the button that opened the dialog unless the function performed in the dialog context logically leads to a different element. For example, activating a cancel button in a dialog returns focus to the button that opened the dialog. However, if the dialog were confirming the action of deleting the page from which it was opened, the focus would logically move to a new context.
- If activating the button does not dismiss the current context, then focus typically remains on the button after activation, e.g., an Apply or Recalculate button.
- If the button action indicates a context change, such as move to next step in a wizard or add another search criteria, then it is often appropriate to move focus to the starting point for that action.
-
 If the button is activated with a shortcut key, the focus usually remains in the context from which the shortcut key was activated.
 For example, if `Alt`+`U`were assigned to an "Up" button that moves the currently focused item in a list one position higher in the list, pressing`Alt`+`U`when the focus is in the list would not move the focus from the list.

## WAI-ARIA Roles, States, and Properties

- The button has role of button.
-
 The `button`has an accessible label. By default, the accessible name is computed from any text content inside the button element. However, it can also be provided with aria-labelledby or aria-label.
- If a description of the button's function is present, the button element has aria-describedby set to the ID of the element containing the description.
- When the action associated with a button is unavailable, the button has aria-disabled set to `true`.
-
 If the button is a toggle button, it has an aria-pressed state.
 When the button is toggled on, the value of this state is `true`, and when toggled off, the state is`false`.

# Carousel (Slide Show or Image Rotator) Pattern

## About This Pattern

A carousel presents a set of items, referred to as slides, by sequentially displaying a subset of one or more slides. Typically, one slide is displayed at a time, and users can activate a next or previous slide control that hides the current slide and "rotates" the next or previous slide into view. In some implementations, rotation automatically starts when the page loads, and it may also automatically stop once all the slides have been displayed. While a slide may contain any type of content, image carousels where each slide contains nothing more than a single image are common.

Ensuring all users can easily control and are not adversely affected by slide rotation is an essential aspect of making carousels accessible. For instance, the screen reader experience can be confusing and disorienting if slides that are not visible on screen are incorrectly hidden, e.g., displayed off-screen. Similarly, if slides rotate automatically and a screen reader user is not aware of the rotation, the user may read an element on slide 1, execute the screen reader command for next element, and, instead of hearing the next element on slide 1, hear an element from slide 2 without any knowledge that the element just announced is from an entirely new context.

Features needed to provide sufficient rotation control include:

- Buttons for displaying the previous and next slides.
- Optionally, a control, or group of controls, for choosing a specific slide to display. For example, slide picker controls can be marked up as tabs in a tablist with the slide represented by a tabpanel element.
-
 If the carousel can automatically rotate, it also:
 - Has a button for stopping and restarting rotation. This is particularly important for supporting assistive technologies operating in a mode that does not move either keyboard focus or the mouse.
- Stops rotating when keyboard focus enters the carousel. It does not restart unless the user explicitly requests it to do so.
- Stops rotating whenever the mouse is hovering over the carousel.

## Examples

- Auto-Rotating Image Carousel with Buttons for Slide Control: A basic image carousel that demonstrates the accessibility features necessary for carousels that rotate automatically on page load and also enables users to choose which slide is displayed with buttons for previous and next slide.
- Auto-Rotating Image Carousel with Tabs for Slide Control: An image carousel that demonstrates accessibility features necessary for carousels that rotate automatically on page load and also enables users to choose which slide is displayed with a set of tabs.

## Terms

The following terms are used to describe components of a carousel.

- Slide
- A single content container within a set of content containers that hold the content to be presented by the carousel.
- Rotation Control
- An interactive element that stops and starts automatic slide rotation.
- Next Slide Control
- An interactive element, often styled as an arrow, that displays the next slide in the rotation sequence.
- Previous Slide Control
- An interactive element, often styled as an arrow, that displays the previous slide in the rotation sequence.
- Slide Picker Controls
- A group of elements, often styled as small dots, that enable the user to pick a specific slide in the rotation sequence to display.

## Keyboard Interaction

- If the carousel has an auto-rotate feature, automatic slide rotation stops when any element in the carousel receives keyboard focus. It does not resume unless the user activates the rotation control.
- `Tab`and- `Shift`+- `Tab`: Move focus through the interactive elements of the carousel as specified by the page tab sequence -- scripting for- `Tab`is not necessary.
- Button elements implement the keyboard interaction defined in the button pattern. Note: Activating the rotation control, next slide, and previous slide do not move focus, so users may easily repetitively activate them as many times as desired.
-
 If present, the rotation control is the first element in the `Tab`sequence inside the carousel. It is essential that it precede the rotating content so it can be easily located.
- If tab elements are used for slide picker controls, they implement the keyboard interaction defined in the Tabs Pattern.

## WAI-ARIA Roles, States, and Properties

This section describes the element composition for three styles of carousels:

- Basic: Has rotation, previous slide, and next slide controls but no slide picker controls.
- Tabbed: Has basic controls plus a single tab stop for slide picker controls implemented using the tabs pattern.
- Grouped: Has basic controls plus a series of tab stops in a group of slide picker controls where each control implements the button pattern. Because each slide selector button adds an element to the page tab sequence, this style is the least friendly for keyboard users.

### Basic carousel elements

- A carousel container element that encompasses all components of the carousel, including both carousel controls and slides, has either role region or role group. The most appropriate role for the carousel container depends on the information architecture of the page. See the Landmark Regions Practice to determine whether the carousel warrants being designated as a landmark region.
- The carousel container has the aria-roledescription property set to `carousel`.
-
 If the carousel has a visible label, its accessible label is provided by the property aria-labelledby on the carousel container set to the ID of the element containing the visible label.
 Otherwise, an accessible label is provided by the property aria-label set on the carousel container.
 Note that since the `aria-roledescription`is set to "carousel", the label does not contain the word "carousel".
- The rotation control, next slide control, and previous slide control are either native button elements (recommended) or implement the button pattern.
-
 The rotation control has an accessible label provided by either its inner text or aria-label.
 The label changes to match the action the button will perform, e.g., "Stop slide rotation" or "Start slide rotation".
 A label that changes when the button is activated clearly communicates both that slide content can change automatically and when it is doing so.
 Note that since the label changes, the rotation control does not have any states, e.g., `aria-pressed`, specified.
- Each slide container has role group with the property aria-roledescription set to `slide`.
-
 Each slide has an accessible name:
 - If a slide has a visible label, its accessible label is provided by the property aria-labelledby on the slide container set to the ID of the element containing the visible label.
- Otherwise, an accessible label is provided by the property aria-label set on the slide container.
- If unique names that identify the slide content are not available, a number and set size can serve as a meaningful alternative, e.g., "3 of 10". Note: Normally, including set position and size information in an accessible name is not appropriate. An exception is helpful in this implementation because group elements do not support aria-setsize or aria-posinset. The tabbed carousel implementation pattern does not have this limitation.
- Note that since the `aria-roledescription`is set to "slide", the label does not contain the word "slide."

-
 Optionally, an element wrapping the set of slide elements has aria-atomic set to `false`and aria-live set to:- `off`: if the carousel is automatically rotating.
- `polite`: if the carousel is- **NOT**automatically rotating.

### Tabbed Carousel Elements

The structure of a tabbed carousel is the same as a basic carousel except that:

- Each slide container has role tabpanel in lieu of `group`, and it does not have the`aria-roledescription`property.
-
 It has slide picker controls implemented using the tabs pattern where:
 - Each control is a `tab`element, so activating a tab displays the slide associated with that tab.
-
 The accessible name of each `tab`indicates which slide it will display by including the name or number of the slide, e.g., "Slide 3". Slide names are preferable if each slide has a unique name.
- The set of controls is grouped in a `tablist`element with an accessible name provided by the value of aria-label that identifies the purpose of the tabs, e.g., "Choose slide to display."
- The `tab`,`tablist`, and`tabpanel`implement the properties specified in the tabs pattern.

- Each control is a

### Grouped Carousel Elements

A grouped carousel has the same structure as a basic carousel, but it also includes slide picker controls where:

- The set of slide picker controls is contained in an element with role group.
- The group containing the picker controls has an accessible label provided by the value of aria-label that identifies the purpose of the controls, e.g., "Choose slide to display."
- Each picker control is a native button element (recommended) or implements the button pattern.
-
 The accessible name of each picker button matches the name of the slide it displays.
 One technique for accomplishing this is to set aria-labelledby to a value that references the slide `group`element.
-
 The picker button representing the currently displayed slide has the property aria-disabled set to `true`. Note:`aria-disabled`is preferable to the HTML`disabled`attribute because this is a circumstance where screen reader users benefit from the disabled button being included in the page`Tab`sequence.

# Checkbox Pattern

## About This Pattern

WAI-ARIA supports two types of checkbox widgets: dual-state checkboxes toggle between two choices -- checked and not checked -- and tri-state checkboxes, which allow an additional third state known as partially checked.

One common use of a tri-state checkbox can be found in software installers where a single tri-state checkbox is used to represent and control the state of an entire group of install options. And, each option in the group can be individually turned on or off with a dual state checkbox.

- If all options in the group are checked, the overall state is represented by the tri-state checkbox displaying as checked.
- If some of the options in the group are checked, the overall state is represented with the tri-state checkbox displaying as partially checked.
- If none of the options in the group are checked, the overall state of the group is represented with the tri-state checkbox displaying as not checked.

The user can use the tri-state checkbox to change all options in the group with a single action:

- Checking the overall checkbox checks all options in the group.
- Unchecking the overall checkbox will uncheck all options in the group.
- And, In some implementations, the system may remember which options were checked the last time the overall status was partially checked. If this feature is provided, activating the overall checkbox a third time recreates that partially checked state where only some options in the group are checked.

## Examples

- Checkbox (Two-State) Example: Demonstrates a simple 2-state checkbox.
- Checkbox (Mixed-State) Example: Demonstrates a checkbox that uses the mixed value for aria-checked to reflect and control checked states within a group of two-state HTML checkboxes contained in an HTML `fieldset`.

## Keyboard Interaction

When the checkbox has focus, pressing the `Space` key changes the state of the checkbox.

## WAI-ARIA Roles, States, and Properties

- The checkbox has role checkbox.
-
 The checkbox has an accessible label provided by one of the following:
 - Visible text content contained within the element with role `checkbox`.
- A visible label referenced by the value of aria-labelledby set on the element with role `checkbox`.
- aria-label set on the element with role `checkbox`.

- Visible text content contained within the element with role
- When checked, the checkbox element has state aria-checked set to `true`.
- When not checked, it has state aria-checked set to `false`.
- When partially checked, it has state aria-checked set to `mixed`.
- If a set of checkboxes is presented as a logical group with a visible label, the checkboxes are included in an element with role group that has the property aria-labelledby set to the ID of the element containing the label.
- If the presentation includes additional descriptive static text relevant to a checkbox or checkbox group, the checkbox or checkbox group has the property aria-describedby set to the ID of the element containing the description.

# Combobox Pattern

## About This Pattern

A combobox is an input widget that has an associated popup. The popup enables users to choose a value for the input from a collection. The popup may be a listbox, grid, tree, or dialog.

 In some implementations, the popup presents allowed values, while in other implementations, the popup presents suggested values.
 Many implementations also include a third optional element -- a graphical Open

 button adjacent to the combobox, which indicates availability of the popup.
 Activating the Open

 button displays the popup if suggestions are available.

 The combobox pattern supports several optional behaviors.
 The one that most shapes interaction is text input.
 Some comboboxes allow users to type and edit text in the combobox and others do not.
 If a combobox does not support text input, it is referred to as select-only, meaning the only way users can set its value is by selecting a value in the popup.
 For example, in some browsers, an HTML `select` element with `size="1"` is presented to assistive technologies as a combobox.
 Alternatively, if a combobox supports text input, it is referred to as editable.
 An editable combobox may either allow users to input any arbitrary value, or it may restrict its value to a discrete set of allowed values, in which case typing input serves to filter suggestions presented in the popup.

The popup is hidden by default, i.e., the default state is collapsed. The conditions that trigger expansion -- display of the popup --are specific to each implementation. Some possible conditions that trigger expansion include:

-
 It is displayed when the `Down Arrow`key is pressed or theOpen button is activated. Optionally, if the combobox is editable and contains a string that does not match an allowed value, expansion does not occur.
- It is displayed when the combobox receives focus even if the combobox is editable and empty.
- If the combobox is editable, the popup is displayed only if a certain number of characters are typed in the combobox and those characters match some portion of one of the suggested values.

Combobox widgets are useful for acquiring user input in either of two scenarios:

- The value must be one of a predefined set of allowed values, e.g., a location field must contain a valid location name. Note that the listbox and menu button patterns are also useful in this scenario; differences between combobox and alternative patterns are described below.
- An arbitrary value is allowed, but it is advantageous to suggest possible values to users. For example, a search field may suggest similar or previous searches to save the user time.

The nature of possible values presented by a popup and the way they are presented is called the autocomplete behavior. Comboboxes can have one of four forms of autocomplete:

-
 **No autocomplete:**The combobox is editable, and when the popup is triggered, the suggested values it contains are the same regardless of the characters typed in the combobox. For example, the popup suggests a set of recently entered values, and the suggestions do not change as the user types.
-
 **List autocomplete with manual selection:**When the popup is triggered, it presents suggested values. If the combobox is editable, the suggested values complete or logically correspond to the characters typed in the combobox. The character string the user has typed will become the value of the combobox unless the user selects a value in the popup.
-
 **List autocomplete with automatic selection:**The combobox is editable, and when the popup is triggered, it presents suggested values that complete or logically correspond to the characters typed in the combobox, and the first suggestion is automatically highlighted as selected. The automatically selected suggestion becomes the value of the combobox when the combobox loses focus unless the user chooses a different suggestion or changes the character string in the combobox.
-
 **List with inline autocomplete:**This is the same as list with automatic selection with one additional feature. The portion of the selected suggestion that has not been typed by the user, a completion string, appears inline after the input cursor in the combobox. The inline completion string is visually highlighted and has a selected state.

If a combobox is editable and has any form of list autocomplete, the popup may appear and disappear as the user types. For example, if the user types a two character string that triggers five suggestions to be displayed but then types a third character that forms a string that does not have any matching suggestions, the popup may close and, if present, the inline completion string disappears.

 Two other widgets that are also visually compact and enable users to make a single choice from a set of discrete choices are listbox and menu button.
 One feature that distinguishes combobox from both listbox and menu button is that the user's choice can be presented as a value in an editable field, which gives users the ability to select some or all of the value for copying to the clipboard.
 Comboboxes and menu buttons can be implemented so users can explore the set of allowed choices without losing a previously made choice.
 That is, users can navigate the set of available choices in a combobox popup or menu and then press `escape`, which closes the popup or menu without changing previous input.
 In contrast, navigating among options in a single-select listbox immediately changes its value, and `Escape` does not provide an undo mechanism.
 Comboboxes and listboxes can be marked as required with `aria-required="true"`, and they have an accessible name that is distinct from their value.
 Thus, when assistive technology users focus either a combobox or listbox in its default state, they can perceive both a name and value for the widget.
 However, a menu button cannot be marked required, and while it has an accessible name, it does not have a value so is not suitable for conveying the user's choice in its collapsed state.

## Examples

- Select-Only Combobox: A single-select combobox with no text input that is functionally similar to an HTML `select`element.
- Editable Combobox with Both List and Inline Autocomplete: An editable combobox that demonstrates the autocomplete behavior known as list with inline autocomplete .
- Editable Combobox with List Autocomplete: An editable combobox that demonstrates the autocomplete behavior known as list with manual selection .
- Editable Combobox Without Autocomplete: An editable combobox that demonstrates the behavior associated with `aria-autocomplete=none`.
- Editable Combobox with Grid Popup: An editable combobox that presents suggestions in a grid, enabling users to navigate descriptive information about each suggestion.
- Date Picker Combobox: An editable date input combobox that opens a dialog containing a calendar grid and buttons for navigating by month and year.

## Keyboard Interaction

- `Tab`: The combobox is in the page- `Tab`sequence.
- Note: The popup indicator icon or button (if present), the popup, and the popup descendants are excluded from the page `Tab`sequence.

### Combobox Keyboard Interaction

When focus is in the combobox:

-
 `Down Arrow`: If the popup is available, moves focus into the popup:- If the autocomplete behavior automatically selected a suggestion before `Down Arrow`was pressed, focus is placed on the suggestion following the automatically selected suggestion.
- Otherwise, places focus on the first focusable element in the popup.

- If the autocomplete behavior automatically selected a suggestion before
- `Up Arrow`(Optional): If the popup is available, places focus on the last focusable element in the popup.
-
 `Escape`: Dismisses the popup if it is visible. Optionally, if the popup is hidden before`Escape`is pressed, clears the combobox.
-
 `Enter`: If the combobox is editable and an autocomplete suggestion is selected in the popup, accepts the suggestion either by placing the input cursor at the end of the accepted value in the combobox or by performing a default action on the value. For example, in a messaging application, the default action may be to add the accepted value to a list of message recipients and then clear the combobox so the user can add another recipient.
-
 Printable Characters:
 - If the combobox is editable, type characters in the combobox. Note that some implementations may regard certain characters as invalid and prevent their input.
- If the combobox is not editable, optionally moves focus to a value that starts with the typed characters.

- If the combobox is editable, it supports standard single line text editing keys appropriate for the device platform (see note below).
- `Alt`+- `Down Arrow`(Optional): If the popup is available but not displayed, displays the popup without moving focus.
-
 `Alt`+`Up Arrow`(Optional): If the popup is displayed:- If the popup contains focus, returns focus to the combobox.
- Closes the popup.

#### Note

Standard single line text editing keys appropriate for the device platform:

- include keys for input, cursor movement, selection, and text manipulation.
- Standard key assignments for editing functions depend on the device operating system.
- The most robust approach for providing text editing functions is to rely on browsers, which supply them for HTML inputs with type text and for elements with the `contenteditable`HTML attribute.
- **IMPORTANT:**Ensure JavaScript does not interfere with browser-provided text editing functions by capturing key events for the keys used to perform them.

### Listbox Popup Keyboard Interaction

When focus is in a listbox popup:

- `Enter`: Accepts the focused option in the listbox by closing the popup, placing the accepted value in the combobox, and if the combobox is editable, placing the input cursor at the end of the value.
-
 `Escape`: Closes the popup and returns focus to the combobox. Optionally, if the combobox is editable, clears the contents of the combobox.
-
 `Down Arrow`: Moves focus to and selects the next option. If focus is on the last option, either returns focus to the combobox or does nothing.
-
 `Up Arrow`: Moves focus to and selects the previous option. If focus is on the first option, either returns focus to the combobox or does nothing.
-
 `Right Arrow`: If the combobox is editable, returns focus to the combobox without closing the popup and moves the input cursor one character to the right. If the input cursor is on the right-most character, the cursor does not move.
-
 `Left Arrow`: If the combobox is editable, returns focus to the combobox without closing the popup and moves the input cursor one character to the left. If the input cursor is on the left-most character, the cursor does not move.
- `Home`(Optional): Either moves focus to and selects the first option or, if the combobox is editable, returns focus to the combobox and places the cursor on the first character.
- `End`(Optional): Either moves focus to the last option or, if the combobox is editable, returns focus to the combobox and places the cursor after the last character.
-
 Any printable character:
 - If the combobox is editable, returns the focus to the combobox without closing the popup and types the character.
- Otherwise, moves focus to the next option with a name that starts with the characters typed.

- `Backspace`(Optional): If the combobox is editable, returns focus to the combobox and deletes the character prior to the cursor.
- `Delete`(Optional): If the combobox is editable, returns focus to the combobox, removes the selected state if a suggestion was selected, and removes the inline autocomplete string if present.

#### Note

-
 DOM Focus is maintained on the combobox and the assistive technology focus is moved within the listbox using `aria-activedescendant`as described in Managing Focus in Composites Using aria-activedescendant.
- Selection follows focus in the listbox; the listbox allows only one suggested value to be selected at a time for the combobox value.

### Grid Popup Keyboard Interaction

In a grid popup, each suggested value may be represented by either a single cell or an entire row. See notes below for how this aspect of grid design effects the keyboard interaction design and the way that selection moves in response to focus movements.

- `Enter`: Accepts the currently selected suggested value by closing the popup, placing the selected value in the combobox, and if the combobox is editable, placing the input cursor at the end of the value.
-
 `Escape`: Closes the popup and returns focus to the combobox. Optionally, if the combobox is editable, clears the contents of the combobox.
-
 `Right Arrow`: Moves focus one cell to the right. Optionally, if focus is on the right-most cell in the row, focus may move to the first cell in the following row. If focus is on the last cell in the grid, either does nothing or returns focus to the combobox.
-
 `Left Arrow`: Moves focus one cell to the left. Optionally, if focus is on the left-most cell in the row, focus may move to the last cell in the previous row. If focus is on the first cell in the grid, either does nothing or returns focus to the combobox.
-
 `Down Arrow`: Moves focus one cell down. If focus is in the last row of the grid, either does nothing or returns focus to the combobox.
-
 `Up Arrow`: Moves focus one cell up. If focus is in the first row of the grid, either does nothing or returns focus to the combobox.
-
 `Page Down`(Optional): Moves focus down an author-determined number of rows, typically scrolling so the bottom row in the currently visible set of rows becomes one of the first visible rows. If focus is in the last row of the grid, focus does not move.
-
 `Page Up`(Optional): Moves focus up an author-determined number of rows, typically scrolling so the top row in the currently visible set of rows becomes one of the last visible rows. If focus is in the first row of the grid, focus does not move.
-
 `Home`(Optional):**Either:**- Moves focus to the first cell in the row that contains focus. Or, if the grid has fewer than three cells per row or multiple suggested values per row, focus may move to the first cell in the grid.
- If the combobox is editable, returns focus to the combobox and places the cursor on the first character.

-
 `End`(Optional):**Either:**- Moves focus to the last cell in the row that contains focus. Or, if the grid has fewer than three cells per row or multiple suggested values per row, focus may move to the last cell in the grid.
- If the combobox is editable, returns focus to the combobox and places the cursor after the last character.

- `Control`+- `Home`(optional): moves focus to the first row.
- `Control`+- `End`(Optional): moves focus to the last row.
- Any printable character: If the combobox is editable, returns the focus to the combobox without closing the popup and types the character.
- `Backspace`(Optional): If the combobox is editable, returns focus to the combobox and deletes the character prior to the cursor.
- `Delete`(Optional): If the combobox is editable, returns focus to the combobox, removes the selected state if a suggestion was selected, and removes the inline autocomplete string if present.

#### Note

-
 DOM Focus is maintained on the combobox and the assistive technology focus is moved within the grid using `aria-activedescendant`as described in Managing Focus in Composites Using aria-activedescendant.
- The grid allows only one suggested value to be selected at a time for the combobox value.
-
 In a grid popup, each suggested value may be represented by either a single cell or an entire row.
 This aspect of design effects focus and selection movement:
 -
 If every cell contains a different suggested value:
 - Selection follows focus so that the cell containing focus is selected.
- Horizontal arrow key navigation typically wraps from one row to another.
- Vertical arrow key navigation typically wraps from one column to another.

-
 If all cells in a row contain information about the same suggested value:
 - Either the row containing focus is selected or a cell containing a suggested value is selected when any cell in the same row contains focus.
- Horizontal key navigation may wrap from one row to another.
- Vertical arrow key navigation **does not**wrap from one column to another.

-
 If every cell contains a different suggested value:

### Tree Popup Keyboard Interaction

In some implementations of tree popups, some or all parent nodes may serve as suggestion category labels so may not be selectable values. See notes below for how this aspect of the design effects the way selection moves in response to focus movements.

When focus is in a vertically oriented tree popup:

- `Enter`: Accepts the currently selected suggested value by closing the popup, placing the selected value in the combobox, and if the combobox is editable, placing the input cursor at the end of the value.
-
 `Escape`: Closes the popup and returns focus to the combobox. Optionally, clears the contents of the combobox.
-
 `Right arrow`:- When focus is on a closed node, opens the node; focus and selection do not move.
- When focus is on a open node, moves focus to the first child node and selects it if it is selectable.
- When focus is on an end node, does nothing.

-
 `Left arrow`:- When focus is on an open node, closes the node.
- When focus is on a child node that is also either an end node or a closed node, moves focus to its parent node and selects it if it is selectable.
- When focus is on a root node that is also either an end node or a closed node, does nothing.

- `Down Arrow`: Moves focus to the next node that is focusable without opening or closing a node and selects it if it is selectable.
- `Up Arrow`: Moves focus to the previous node that is focusable without opening or closing a node and selects it if it is selectable.
- `Home`: Moves focus to the first focusable node in the tree without opening or closing a node and selects it if it is selectable.
- `End`: Moves focus to the last focusable node in the tree without opening or closing a node and selects it if it is selectable.
-
 Any printable character:
 - If the combobox is editable, returns the focus to the combobox without closing the popup and types the character.
- Otherwise, moves focus to the next suggested value with a name that starts with the characters typed.

#### Note

-
 DOM Focus is maintained on the combobox and the assistive technology focus is moved within the tree using `aria-activedescendant`as described in Managing Focus in Composites Using aria-activedescendant.
- The tree allows only one suggested value to be selected at a time for the combobox value.
-
 In a tree popup, some or all parent nodes may not be selectable values; they may serve as category labels for suggested values.
 If focus moves to a node that is not a selectable value, either:
 - The previously selected node, if any, remains selected until focus moves to a node that is selectable.
- There is no selected value.
- In either case, focus is visually distinct from selection so users can readily see if a value is selected or not.

-
 If nodes in the tree are arranged horizontally (aria-orientation is set to `horizontal`):- `Down Arrow`performs as- `Right Arrow`is described above, and vice versa.
- `Up Arrow`performs as- `Left Arrow`is described above, and vice versa.

### Dialog Popup Keyboard Interaction

When focus is in a dialog popup:

-
 There are two ways to close the popup and return focus to the combobox:
 - Perform an action in the dialog, such as activate a button, that specifies a value for the combobox.
-
 Cancel out of the dialog, e.g., press `Escape`or activate the cancel button in the dialog. Canceling either returns focus to the text box without changing the combobox value or returns focus to the combobox and clears the combobox.

- The dialog implements the keyboard interaction defined in the modal dialog pattern.

#### Note

Unlike other combobox popups, dialogs do not support `aria-activedescendant` so DOM focus moves into the dialog from the combobox.

## WAI-ARIA Roles, States, and Properties

- The element that serves as an input and displays the combobox value has role combobox.
-
 The combobox element has aria-controls set to a value that refers to the element that serves as the popup.
 Note that `aria-controls`only needs to be set when the popup is visible. However, it is valid to reference an element that is not visible.
- The popup is an element that has role listbox, tree, grid, or dialog.
-
 If the popup has a role other than `listbox`, the element with role`combobox`has aria-haspopup set to a value that corresponds to the popup type. That is,`aria-haspopup`is set to`grid`,`tree`, or`dialog`. Note that elements with role`combobox`have an implicit`aria-haspopup`value of`listbox`.
-
 When the combobox popup is not visible, the element with role `combobox`has aria-expanded set to`false`. When the popup element is visible,`aria-expanded`is set to`true`. Note that elements with role`combobox`have a default value for`aria-expanded`of`false`.
- When a combobox receives focus, DOM focus is placed on the combobox element.
- When a descendant of a listbox, grid, or tree popup is focused, DOM focus remains on the combobox and the combobox has aria-activedescendant set to a value that refers to the focused element within the popup.
- For a combobox that controls a listbox, grid, or tree popup, when a suggested value is visually indicated as the currently selected value, the `option`,`gridcell`,`row`, or`treeitem`containing that value has aria-selected set to`true`.
-
 If the combobox has a visible label and the combobox element is an HTML element that can be labelled using the HTML `label`element (e.g., the`input`element), it is labeled using the`label`element. Otherwise, if it has a visible label, the combobox element has aria-labelledby set to a value that refers to the labelling element. Otherwise, the combobox element has a label provided by aria-label.
-
 The combobox element has aria-autocomplete set to a value that corresponds to its autocomplete behavior:
 - `none`: When the popup is displayed, the suggested values it contains are the same regardless of the characters typed in the combobox.
-
 `list`: When the popup is triggered, it presents suggested values. If the combobox is editable, the values complete or logically correspond to the characters typed in the combobox.
-
 `both`: When the popup is triggered, it presents suggested values that complete or logically correspond to the characters typed in the combobox. In addition, the portion of the selected suggestion that has not been typed by the user, known as thecompletion string , appears inline after the input cursor in the combobox. The inline completion string is visually highlighted and has a selected state.

### Note

- In ARIA 1.0, the combobox referenced its popup with aria-owns instead of aria-controls. While user agents might support comboboxes with aria-owns for backwards compatibility with legacy content, it is strongly recommended that authors use aria-controls.
- When referring to the roles, states, and properties documentation for the below list of patterns used for popups, keep in mind that a combobox is a single-select widget where selection follows focus in the popup.
- The roles, states, and properties for popup elements are defined in their respective design patterns:

# Dialog (Modal) Pattern

## About This Pattern

A dialog is a window overlaid on either the primary window or another dialog window. Windows under a modal dialog are inert. That is, users cannot interact with content outside an active dialog window. Inert content outside an active dialog is typically visually obscured or dimmed so it is difficult to discern, and in some implementations, attempts to interact with the inert content cause the dialog to close.

 Like non-modal dialogs, modal dialogs contain their tab sequence.
 That is, `Tab` and `Shift` + `Tab` do not move focus outside the dialog.
 However, unlike most non-modal dialogs, modal dialogs do not provide means for moving keyboard focus outside the dialog window without closing the dialog.

The alertdialog role is a special-case dialog role designed specifically for dialogs that divert users' attention to a brief, important message. Its usage is described in the Alert Dialog Pattern.

## Examples

- Modal Dialog Example: Demonstrates multiple layers of modal dialogs with both small and large amounts of content.
- Date Picker Dialog Example: Demonstrates a dialog containing a calendar grid for choosing a date.

## Keyboard Interaction

 In the following description, the term tabbable element

 refers to any element with a `tabindex` value of zero or greater.
 Note that values greater than 0 are strongly discouraged.

- When a dialog opens, focus moves to an element inside the dialog. See notes below regarding initial focus placement.
-
 `Tab`:- Moves focus to the next tabbable element inside the dialog.
- If focus is on the last tabbable element inside the dialog, moves focus to the first tabbable element inside the dialog.

-
 `Shift`+`Tab`:- Moves focus to the previous tabbable element inside the dialog.
- If focus is on the first tabbable element inside the dialog, moves focus to the last tabbable element inside the dialog.

- `Escape`: Closes the dialog.

### Note

-
 When a dialog opens, focus moves to an element contained in the dialog.
 Generally, focus is initially set on the first focusable element.
 However, the most appropriate focus placement will depend on the nature and size of the content.
 Examples include:
 -
 If the dialog content includes semantic structures, such as lists, tables, or multiple paragraphs, that need to be perceived in order to easily understand the content, i.e., if the content would be difficult to understand when announced as a single unbroken string, then it is advisable to add `tabindex="-1"`to a static element at the start of the content and initially focus that element. This makes it easier for assistive technology users to read the content by navigating the semantic structures. Additionally, it is advisable to omit applying`aria-describedby`to the dialog container in this scenario.
- If content is large enough that focusing the first interactive element could cause the beginning of content to scroll out of view, it is advisable to add `tabindex="-1"`to a static element at the top of the dialog, such as the dialog title or first paragraph, and initially focus that element.
- If a dialog contains the final step in a process that is not easily reversible, such as deleting data or completing a financial transaction, it may be advisable to set focus on the least destructive action, especially if undoing the action is difficult or impossible. The Alert Dialog Pattern is often employed in such circumstances.
- If a dialog is limited to interactions that either provide additional information or continue processing, it may be advisable to set focus to the element that is likely to be most frequently used, such as an OK orContinue button.

-
 If the dialog content includes semantic structures, such as lists, tables, or multiple paragraphs, that need to be perceived in order to easily understand the content, i.e., if the content would be difficult to understand when announced as a single unbroken string, then it is advisable to add
-
 When a dialog closes, focus returns to the element that invoked the dialog unless either:
 - The invoking element no longer exists. Then, focus is set on another element that provides logical work flow.
-
 The work flow design includes the following conditions that can occasionally make focusing a different element a more logical choice:
 - It is very unlikely users need to immediately re-invoke the dialog.
- The task completed in the dialog is directly related to a subsequent step in the work flow.

- It is strongly recommended that the tab sequence of all dialogs include a visible element with role `button`that closes the dialog, such as a close icon or cancel button.

## WAI-ARIA Roles, States, and Properties

- The element that serves as the dialog container has a role of dialog.
- All elements required to operate the dialog are descendants of the element that has role `dialog`.
- The dialog container element has aria-modal set to `true`.
-
 The dialog has either:
 - A value set for the aria-labelledby property that refers to a visible dialog title.
- A label specified by aria-label.

-
 Optionally, the aria-describedby property is set on the element with the `dialog`role to indicate which element or elements in the dialog contain content that describes the primary purpose or message of the dialog. Specifying descriptive elements enables screen readers to announce the description along with the dialog title and initially focused element when the dialog opens, which is typically helpful only when the descriptive content is simple and can easily be understood without structural information. It is advisable to omit specifying`aria-describedby`if the dialog content includes semantic structures, such as lists, tables, or multiple paragraphs, that need to be perceived in order to easily understand the content, i.e., if the content would be difficult to understand when announced as a single unbroken string.

### Note

-
 Because marking a dialog modal by setting aria-modal to `true`can prevent users of some assistive technologies from perceiving content outside the dialog, users of those technologies will experience severe negative ramifications if a dialog is marked modal but does not behave as a modal for other users. So, mark a dialog modal**only when both:**- Application code prevents all users from interacting in any way with content outside of it.
- Visual styling obscures the content outside of it.

-
 The `aria-modal`property introduced by ARIA 1.1 replaces aria-hidden for informing assistive technologies that content outside a dialog is inert. However, in legacy dialog implementations where`aria-hidden`is used to make content outside a dialog inert for assistive technology users, it is important that:- `aria-hidden`is set to- `true`on each element containing a portion of the inert layer.
- The dialog element is not a descendant of any element that has `aria-hidden`set to`true`.

# Disclosure (Show/Hide) Pattern

## About This Pattern

A disclosure is a widget that enables content to be either collapsed (hidden) or expanded (visible). It has two elements: a disclosure button and a section of content whose visibility is controlled by the button. When the controlled content is hidden, the button is often styled as a typical push button with a right-pointing arrow or triangle to hint that activating the button will display additional content. When the content is visible, the arrow or triangle typically points down.

## Examples

## Keyboard Interaction

When the disclosure control has focus:

- `Enter`: activates the disclosure control and toggles the visibility of the disclosure content.
- `Space`: activates the disclosure control and toggles the visibility of the disclosure content.

## WAI-ARIA Roles, States, and Properties

- The element that shows and hides the content has role button.
-
 When the content is visible, the element with role `button`has aria-expanded set to`true`. When the content area is hidden, it is set to`false`.
-
 Optionally, the element with role `button`has a value specified for aria-controls that refers to the element that contains all the content that is shown or hidden.

# Feed Pattern

## About This Pattern

A feed is a section of a page that automatically loads new sections of content as the user scrolls. The sections of content in a feed are presented in article elements. So, a feed can be thought of as a dynamic list of articles that often appears to scroll infinitely.

The feature that most distinguishes feed from other ARIA patterns that support loading data as users scroll, e.g., a grid, is that a feed is a structure, not a widget. Consequently, assistive technologies with a reading mode, such as screen readers, default to reading mode when interacting with feed content. However, unlike most other WAI-ARIA structures, a feed establishes an interoperability contract between the web page and assistive technologies. The contract governs scroll interactions so that assistive technology users can read articles, jump forward and backward by article, and reliably trigger new articles to load while in reading mode.

For example, a product page on a shopping site may have a related products section that displays five products at a time. As the user scrolls, more products are requested and loaded into the DOM. While a static design might include a next button for loading five more products, a dynamic implementation that automatically loads more data as the user scrolls simplifies the user experience and reduces the inertia associated with viewing more than the first five product suggestions. But, unfortunately when web pages load content dynamically based on scroll events, it can cause usability and interoperability difficulties for users of assistive technologies.

The feed pattern enables reliable assistive technology reading mode interaction by establishing the following interoperability agreement between the web page and assistive technologies:

-
 In the context of a feed, the web page code is responsible for:
 - Appropriate visual scrolling of the content based on which article contains DOM focus.
- Loading or removing feed articles based on which article contains DOM focus.

-
 In the context of a feed, assistive technologies with a reading mode are responsible for:
 - Indicating which article contains the reading cursor by ensuring the article element or one of its descendants has DOM focus.
- providing reading mode keys that move DOM focus to the next and previous articles.
- Providing reading mode keys for moving the reading cursor and DOM focus past the end and before the start of the feed.

Thus, implementing the feed pattern allows a screen reader to reliably read and trigger the loading of feed content while staying in its reading mode.

Another feature of the feed pattern is its ability to facilitate skim reading for assistive technology users. Web page authors may provide both an accessible name and description for each article. By identifying the elements inside of an article that provide the title and the primary content, assistive technologies can provide functions that enable users to jump from article to article and efficiently discern which articles may be worthy of more attention.

## Example

## Keyboard Interaction

 The feed pattern is not based on a desktop GUI widget so the `feed` role is not associated with any well-established keyboard conventions.
 Supporting the following, or a similar, interface is recommended.

When focus is inside the feed:

- `Page Down`: Move focus to next article.
- `Page Up`: Move focus to previous article.
- `Control`+- `End`: Move focus to the first focusable element after the feed.
- `Control`+- `Home`: Move focus to the first focusable element before the feed.

### Note

- Due to the lack of convention, providing easily discoverable keyboard interface documentation is especially important.
-
 In some cases, a feed may contain a nested feed.
 For example, an article in a social media feed may contain a feed of comments on that article.
 To navigate the nested feed, users first move focus inside the nested feed.
 Options for supporting nested feed navigation include:
 -
 Users move focus into the nested feed from the content of the containing article with `Tab`. This may be slow if the article contains a significant number of links, buttons, or other widgets.
- Provide a key for moving focus from the elements in the containing article to the first item in the nested feed, e.g., `Alt`+`Page Down`.
- To continue reading the outer feed, `Control`+`End`moves focus to the next article in the outer feed.

-
 Users move focus into the nested feed from the content of the containing article with
- In the rare circumstance that a feed article contains a widget that uses the above suggested keys, the feed navigation key will operate the contained widget, and the user needs to move focus to an element that does not utilize the feed navigation keys in order to navigate the feed.

## WAI-ARIA Roles, States, and Properties

- The element that contains the set of feed articles has role feed.
-
 If the feed has a visible label, the `feed`element has aria-labelledby referring to the element containing the title. Otherwise, the`feed`element has a label specified with aria-label.
-
 Each unit of content in a feed is contained in an element with role article.
 All content inside the feed is contained in an `article`element.
- Each `article`element has aria-labelledby referring to elements inside the article that can serve as a distinguishing label.
- It is optional but strongly recommended for each `article`element to have aria-describedby referring to one or more elements inside the article that serve as the primary content of the article.
- Each `article`element has aria-posinset set to a value that represents its position in the feed.
-
 Each `article`element has aria-setsize set to a value that represents either the total number of articles that have been loaded or the total number in the feed, depending on which value is deemed more helpful to users. If the total number in the feed is undetermined, it can be represented by a`aria-setsize`value of`-1`.
-
 When `article`elements are being added to or removed from the`feed`container, and if the operation requires multiple DOM operations, the`feed`element has aria-busy set to`true`during the update operation. Note that it is extremely important that`aria-busy`is set to`false`when the operation is complete or the changes may not become visible to some assistive technology users.

# Grid (Interactive Tabular Data and Layout Containers) Pattern

## About This Pattern

 A grid widget is a container that enables users to navigate the information or interactive elements it contains using directional navigation keys, such as arrow keys, `Home`, and `End`.
 As a generic container widget that offers flexible keyboard navigation, it can serve a wide variety of needs.
 It can be used for purposes as simple as grouping a collection of checkboxes or navigation links or as complex as creating a full-featured spreadsheet application.
 While the words "row" and "column" are used in the names of WAI-ARIA attributes and by assistive technologies when describing and presenting the logical structure of elements with the `grid` role, using the `grid` role on an element does not necessarily imply that its visual presentation is tabular.

When presenting content that is tabular, consider the following factors when choosing between implementing this `grid` pattern or the table pattern.

-
 A `grid`is a composite widget so it:- Always contains multiple focusable elements.
- Only one of the focusable elements contained by the grid is included in the page tab sequence.
- Requires the author to provide code that manages focus movement inside it.

- All focusable elements contained in a table are included in the page tab sequence.

 Uses of the `grid` pattern broadly fall into two categories: presenting tabular information (data grids) and grouping other widgets (layout grids).
 Even though both data grids and layout grids employ the same ARIA roles, states, and properties, differences in their content and purpose surface factors that are important to consider in keyboard interaction design.
 To address these factors, the following two sections describe separate keyboard interaction patterns for data and layout grids.

## Examples

- Layout Grid Examples: Three example implementations of grids that are used to lay out widgets, including a collection of navigation links, a message recipients list, and a set of search results.
- Data Grid Examples: Three example implementations of grid that include features relevant to presenting tabular information, such as content editing, sort, and column hiding.
- Advanced Data Grid Example: Example of a grid with behaviors and features similar to a typical spreadsheet, including cell and row selection.

## Data Grids For Presenting Tabular Information

 A `grid` can be used to present tabular information that has column titles, row titles, or both.
 The `grid` pattern is particularly useful if the tabular information is editable or interactive.
 For example, when data elements are links to more information, rather than presenting them in a static table and including the links in the tab sequence, implementing the `grid` pattern provides users with intuitive and efficient keyboard navigation of the grid contents as well as a shorter tab sequence for the page.
 A `grid` may also offer functions, such as cell content editing, selection, cut, copy, and paste.

In a grid, every cell contains a focusable element or is itself focusable, regardless of whether the cell content is editable or interactive. There is one exception: if column or row header cells do not provide functions, such as sort or filter, they do not need to be focusable. One reason it is important for all cells to be able to receive or contain keyboard focus is that screen readers will typically be in their application reading mode, rather than their document reading mode, when users are interacting with the grid. While in application mode, a screen reader user hears only focusable elements and content that labels focusable elements. So, screen reader users may unknowingly overlook elements contained in a grid that are either not focusable or not used to label a column or row.

### Keyboard Interaction For Data Grids

 The following keys provide grid navigation by moving focus among cells of the grid.
 Implementations of grid make these key commands available when an element in the grid has received focus, e.g., after a user has moved focus to the grid with `Tab`.

-
 `Right Arrow`: Moves focus one cell to the right. If focus is on the right-most cell in the row, focus does not move.
-
 `Left Arrow`: Moves focus one cell to the left. If focus is on the left-most cell in the row, focus does not move.
-
 `Down Arrow`: Moves focus one cell down. If focus is on the bottom cell in the column, focus does not move.
-
 `Up Arrow`: Moves focus one cell up. If focus is on the top cell in the column, focus does not move.
-
 `Page Down`: Moves focus down an author-determined number of rows, typically scrolling so the bottom row in the currently visible set of rows becomes one of the first visible rows. If focus is in the last row of the grid, focus does not move.
-
 `Page Up`: Moves focus up an author-determined number of rows, typically scrolling so the top row in the currently visible set of rows becomes one of the last visible rows. If focus is in the first row of the grid, focus does not move.
- `Home`: moves focus to the first cell in the row that contains focus.
- `End`: moves focus to the last cell in the row that contains focus.
- `Control`+- `Home`: moves focus to the first cell in the first row.
- `Control`+- `End`: moves focus to the last cell in the last row.

#### Note

- When the above grid navigation keys move focus, whether the focus is set on an element inside the cell or the grid cell depends on cell content. See Whether to Focus on a Cell or an Element Inside It.
- While navigation keys, such as arrow keys, are moving focus from cell to cell, they are not available to do something like operate a combobox or move an editing caret inside of a cell. If this functionality is needed, see Editing and Navigating Inside a Cell.
- If navigation functions can dynamically add more rows or columns to the DOM, key events that move focus to the beginning or end of the grid, such as `Control`+`End`, may move focus to the last row in the DOM rather than the last available row in the back-end data.

If a grid supports selection of cells, rows, or columns, the following keys are commonly used for these functions.

- `Control`+- `Space`: selects the column that contains the focus.
-
 `Shift`+`Space`: Selects the row that contains the focus. If the grid includes a column with checkboxes for selecting rows, this key can serve as a shortcut for checking the box when focus is not on the checkbox.
- `Control`+- `A`: Selects all cells.
- `Shift`+- `Right Arrow`: Extends selection one cell to the right.
- `Shift`+- `Left Arrow`: Extends selection one cell to the left.
- `Shift`+- `Down Arrow`: Extends selection one cell down.
- `Shift`+- `Up Arrow`: Extends selection one cell up.

##### Note

See Key Assignment Conventions for Common Functions for cut, copy, and paste key assignments.

## Layout Grids for Grouping Widgets

 The `grid` pattern can be used to group a set of interactive elements, such as links, buttons, or checkboxes.
 Since only one element in the entire grid is included in the tab sequence, grouping with a grid can dramatically reduce the number of tab stops on a page.
 This is especially valuable if scrolling through a list of elements dynamically loads more of those elements from a large data set, such as in a continuous list of suggested products on a shopping site.
 If elements in a list like this were in the tab sequence, keyboard users are effectively trapped in the list.
 If any elements in the group also have associated elements that appear on hover, the `grid` pattern is also useful for providing keyboard access to those contextual elements of the user interface.

 Unlike grids used to present data, A `grid` used for layout does not necessarily have header cells for labelling rows or columns and might contain only a single row or a single column.
 Even if it has multiple rows and columns, it may present a single, logically homogenous set of elements.
 For example, a list of recipients for a message may be a grid where each cell contains a link that represents a recipient.
 The grid may initially have a single row but then wrap into multiple rows as recipients are added.
 In such circumstances, grid navigation keys may also wrap so the user can read the list from beginning to end by pressing either `Right Arrow` or `Down Arrow`.
 While This type of focus movement wrapping can be very helpful in a layout grid, it would be disorienting if used in a data grid, especially for users of assistive technologies.

 Because arrow keys are used to move focus inside of a `grid`, a `grid` is both easier to build and use if the components it contains do not require the arrow keys to operate.
 If a cell contains an element like a listbox, then an extra key command to focus and activate the listbox is needed as well as a command for restoring the grid navigation functionality.
 Approaches to supporting this need are described in the section on Editing and Navigating Inside a Cell.

### Keyboard Interaction For Layout Grids

 The following keys provide grid navigation by moving focus among cells of the grid.
 Implementations of grid make these key commands available when an element in the grid has received focus, e.g., after a user has moved focus to the grid with `Tab`.

-
 `Right Arrow`: Moves focus one cell to the right. Optionally, if focus is on the right-most cell in the row, focus may move to the first cell in the following row. If focus is on the last cell in the grid, focus does not move.
-
 `Left Arrow`: Moves focus one cell to the left. Optionally, if focus is on the left-most cell in the row, focus may move to the last cell in the previous row. If focus is on the first cell in the grid, focus does not move.
-
 `Down Arrow`: Moves focus one cell down. Optionally, if focus is on the bottom cell in the column, focus may move to the top cell in the following column. If focus is on the last cell in the grid, focus does not move.
-
 `Up Arrow`: Moves focus one cell up. Optionally, if focus is on the top cell in the column, focus may move to the bottom cell in the previous column. If focus is on the first cell in the grid, focus does not move.
-
 `Page Down`(Optional): Moves focus down an author-determined number of rows, typically scrolling so the bottom row in the currently visible set of rows becomes one of the first visible rows. If focus is in the last row of the grid, focus does not move.
-
 `Page Up`(Optional): Moves focus up an author-determined number of rows, typically scrolling so the top row in the currently visible set of rows becomes one of the last visible rows. If focus is in the first row of the grid, focus does not move.
-
 `Home`: moves focus to the first cell in the row that contains focus. Optionally, if the grid has a single column or fewer than three cells per row, focus may instead move to the first cell in the grid.
-
 `End`: moves focus to the last cell in the row that contains focus. Optionally, if the grid has a single column or fewer than three cells per row, focus may instead move to the last cell in the grid.
- `Control`+- `Home`(optional): moves focus to the first cell in the first row.
- `Control`+- `End`(Optional): moves focus to the last cell in the last row.

#### Note

- When the above grid navigation keys move focus, whether the focus is set on an element inside the cell or the grid cell depends on cell content. See Whether to Focus on a Cell or an Element Inside It.
- While navigation keys, such as arrow keys, are moving focus from cell to cell, they are not available to do something like operate a combobox or move an editing caret inside of a cell. If this functionality is needed, see Editing and Navigating Inside a Cell.
- If navigation functions can dynamically add more rows or columns to the DOM, key events that move focus to the beginning or end of the grid, such as `control`+`End`, may move focus to the last row in the DOM rather than the last available row in the back-end data.

It would be unusual for a layout grid to provide functions that require cell selection. If it did, though, the following keys are commonly used for these functions.

- `Control`+- `Space`: selects the column that contains the focus.
-
 `Shift`+`Space`: Selects the row that contains the focus. If the grid includes a column with checkboxes for selecting rows, this key can serve as a shortcut for checking the box when focus is not on the checkbox.
- `Control`+- `A`: Selects all cells.
- `Shift`+- `Right Arrow`: Extends selection one cell to the right.
- `Shift`+- `Left Arrow`: Extends selection one cell to the left.
- `Shift`+- `Down Arrow`: Extends selection one cell down.
- `Shift`+- `Up Arrow`: Extends selection one cell up.

#### Note

See Key Assignment Conventions for Common Functions for cut, copy, and paste key assignments.

## WAI-ARIA Roles, States, and Properties

- The grid container has role grid.
- Each row container has role row and is either a DOM descendant of or owned by the `grid`element or an element with role rowgroup.
-
 Each cell is either a DOM descendant of or owned by a `row`element and has one of the following roles:- columnheader if the cell contains a title or header information for the column.
- rowheader if the cell contains title or header information for the row.
- gridcell if the cell does not contain column or row header information.

- If there is an element in the user interface that serves as a label for the grid, aria-labelledby is set on the grid element with a value that refers to the labelling element. Otherwise, a label is specified for the grid element using aria-label.
- If the grid has a caption or description, aria-describedby is set on the grid element with a value referring to the element containing the description.
- If the grid provides sort functions, aria-sort is set to an appropriate value on the header cell element for the sorted column or row as described in the Grid and Table Properties Practice.
-
 If the grid supports selection, when a cell or row is selected, the selected element has aria-selected set `true`. If the grid supports column selection and a column is selected, all cells in the column have`aria-selected`set to`true`.
-
 If the grid provides content editing functionality and contains cells that may have edit capabilities disabled in certain conditions, aria-readonly may be set `true`on cells where editing is disabled. If edit functions are disabled for all cells,`aria-readonly`may be set`true`on the grid element. Grids that do not provide editing functions do not include the`aria-readonly`attribute on any of their elements.
-
 If there are conditions where some rows or columns are hidden or not present in the DOM, e.g., data is dynamically loaded when scrolling or the grid provides functions for hiding rows or columns, the following properties are applied as described in the
 Grid and Table Properties Practice.
 - aria-colcount or aria-rowcount is set to the total number of columns or rows, respectively.
- aria-colindex or aria-rowindex is set to the position of a cell within a row or column, respectively.

-
 If the grid includes cells that span multiple rows or multiple columns, and if the `grid`role is NOT applied to an HTML`table`element, then aria-rowspan or aria-colspan is applied as described in the Grid and Table Properties Practice.

### Note

-
 If the element with the `grid`role is an HTML`table`element, then it is not necessary to use ARIA roles for rows and cells because the HTML elements have implied ARIA semantics. For example, an HTML`<TR>`has an implied ARIA role of`row`. A`grid`built from an HTML`table`that includes cells that span multiple rows or columns must use HTML`rowspan`and`colspan`and must not use`aria-rowspan`or`aria-colspan`.
- If rows or cells are included in a grid via aria-owns, they will be presented to assistive technologies after the DOM descendants of the `grid`element unless the DOM descendants are also included in the`aria-owns`attribute.

# Landmarks Pattern

## About This Pattern

Landmarks are a set of eight roles that identify the major sections of a page. Each landmark role enables assistive technology users to perceive the start and end of a feature of the high-level page structure that is usually conveyed visually with placement, spacing, color, or borders. For example, the main landmark designates the section that contains the main content of the page. In addition to conveying structure, landmarks enable browsers and assistive technologies to facilitate efficient keyboard navigation among sections of a page.

 Several landmark roles are implied by HTML elements.
 For example, the HTML `main` element automatically creates a main landmark region, and the HTML `nav` element creates a navigation landmark region.

Since landmarks are intended to help assistive technology users perceive the high-level structure of a page, their value diminishes as their number grows. For optimum value, a general rule of thumb is that a page contains seven or fewer landmark regions. Another best practice is to ensure that all content is contained within an appropriate landmark region. The Landmark Regions Practice describes ways of using HTML sectioning elements and ARIA landmark roles that will most benefit users.

## Examples

## Keyboard Interaction

Not applicable.

## WAI-ARIA Roles, States, and Properties

The Landmark Regions Practice describes the HTML elements, roles, properties, and usage guidelines for each of the landmark region roles.

# Link Pattern

## About This Pattern

A link widget provides an interactive reference to a resource. The target resource can be either external or local, i.e., either outside or within the current page or application.

### note

Authors are strongly encouraged to use a native host language link element, such as an HTML `<A>` element with an `href` attribute.
 As with other WAI-ARIA widget roles, applying the link role to an element will not cause browsers to enhance the element with standard link behaviors, such as navigation to the link target or context menu actions.
 When using the `link` role, providing these features of the element is the author's responsibility.

## Examples

Link Examples: Link widgets constructed from HTML `span` and `img` elements.

## Keyboard Interaction

- `Enter`: Executes the link and moves focus to the link target.
- `Shift`+- `F10`(Optional): Opens a context menu for the link.

## WAI-ARIA Roles, States, and Properties

The element containing the link text or graphic has role of link.

# Listbox Pattern

## About This Pattern

A listbox widget presents a list of options and allows a user to select one or more of them. A listbox that allows a single option to be chosen is a single-select listbox; one that allows multiple options to be selected is a multi-select listbox.

When screen readers present a listbox, they may render the name, state, and position of each option in the list. The name of an option is a string calculated by the browser, typically from the content of the option element. As a flat string, the name does not contain any semantic information. Thus, if an option contains a semantic element, such as a heading, screen reader users will not have access to the semantics. In addition, the interaction model conveyed by the listbox role to assistive technologies does not support interacting with elements inside of an option. Because of these traits of the listbox widget, it does not provide an accessible way to present a list of interactive elements, such as links, buttons, or checkboxes. To present a list of interactive elements, see the Grid Pattern.

Avoiding very long option names facilitates understandability and perceivability for screen reader users. The entire name of an option is spoken as a single unit of speech when the option is read. When too much information is spoken as the result of a single key press, it is difficult to understand. Long names inhibit perception by increasing the impact of interrupted speech because users typically have to re-read the entire option. And, if the user does not understand what is spoken, reading the name by character, word, or phrase may be a difficult operation for many screen reader users in the context of a listbox widget.

Sets of options where each option name starts with the same word or phrase can also significantly degrade usability for keyboard and screen reader users. Scrolling through the list to find a specific option becomes inordinately time consuming for a screen reader user who must listen to that word or phrase repeated before hearing what is unique about each option. For example, if a listbox for choosing a city were to contain options where each city name were preceded by a country name, and if many cities were listed for each country, a screen reader user would have to listen to the country name before hearing each city name. In such a scenario, it would be better to have 2 list boxes, one for country and one for city.

## Examples

- Scrollable Listbox Example: Single-select listbox that scrolls to reveal more options, similar to HTML `select`with`size`attribute greater than one.
- Example Listboxes with Rearrangeable Options: Examples of both single-select and multi-select listboxes with accompanying toolbars where options can be added, moved, and removed.
- Listbox Example with Grouped Options: Single-select listbox with grouped options, similar to an HTML `select`with`optgroup`children.

## Keyboard Interaction

For a vertically oriented listbox:

-
 When a single-select listbox receives focus:
 - If none of the options are selected before the listbox receives focus, the first option receives focus. Optionally, the first option may be automatically selected.
- If an option is selected before the listbox receives focus, focus is set on the selected option.

-
 When a multi-select listbox receives focus:
 - If none of the options are selected before the listbox receives focus, focus is set on the first option and there is no automatic change in the selection state.
- If one or more options are selected before the listbox receives focus, focus is set on the first option in the list that is selected.

-
 `Down Arrow`: Moves focus to the next option. Optionally, in a single-select listbox, selection may also move with focus.
-
 `Up Arrow`: Moves focus to the previous option. Optionally, in a single-select listbox, selection may also move with focus.
-
 `Home`(Optional): Moves focus to first option. Optionally, in a single-select listbox, selection may also move with focus. Supporting this key is strongly recommended for lists with more than five options.
-
 `End`(Optional): Moves focus to last option. Optionally, in a single-select listbox, selection may also move with focus. Supporting this key is strongly recommended for lists with more than five options.
-
 Type-ahead is recommended for all listboxes, especially those with more than seven options:
 - Type a character: focus moves to the next item with a name that starts with the typed character.
- Type multiple characters in rapid succession: focus moves to the next item with a name that starts with the string of characters typed.

-
 **Multiple Selection**: Authors may implement either of two interaction models to support multiple selection: a recommended model that does not require the user to hold a modifier key, such as`Shift`or`Control`, while navigating the list or an alternative model that does require modifier keys to be held while navigating in order to avoid losing selection states.-
 Recommended selection model -- holding modifier keys is not necessary:
 - `Space`: changes the selection state of the focused option.
- `Shift`+- `Down Arrow`(Optional): Moves focus to and toggles the selected state of the next option.
- `Shift`+- `Up Arrow`(Optional): Moves focus to and toggles the selected state of the previous option.
- `Shift`+- `Space`(Optional): Selects contiguous items from the most recently selected item to the focused item.
-
 `Control`+`Shift`+`Home`(Optional): Selects the focused option and all options up to the first option. Optionally, moves focus to the first option.
-
 `Control`+`Shift`+`End`(Optional): Selects the focused option and all options down to the last option. Optionally, moves focus to the last option.
-
 `Control`+`A`(Optional): Selects all options in the list. Optionally, if all options are selected, it may also unselect all options.

-
 Alternative selection model -- moving focus without holding a `Shift`or`Control`modifier unselects all selected options except the focused option:- `Shift`+- `Down Arrow`: Moves focus to and toggles the selection state of the next option.
- `Shift`+- `Up Arrow`: Moves focus to and toggles the selection state of the previous option.
- `Control`+- `Down Arrow`: Moves focus to the next option without changing its selection state.
- `Control`+- `Up Arrow`: Moves focus to the previous option without changing its selection state.
- `Control`+- `Space`Changes the selection state of the focused option.
- `Shift`+- `Space`(Optional): Selects contiguous items from the most recently selected item to the focused item.
-
 `Control`+`Shift`+`Home`(Optional): Selects the focused option and all options up to the first option. Optionally, moves focus to the first option.
-
 `Control`+`Shift`+`End`(Optional): Selects the focused option and all options down to the last option. Optionally, moves focus to the last option.
-
 `Control`+`A`(Optional): Selects all options in the list. Optionally, if all options are selected, it may also unselect all options.

-
 Recommended selection model -- holding modifier keys is not necessary:

### Note

- DOM focus (the active element) is functionally distinct from the selected state. For more details, see this description of differences between focus and selection.
-
 The `listbox`role supports the aria-activedescendant property, which provides an alternative to moving DOM focus among`option`elements when implementing keyboard navigation. For details, see Managing Focus in Composites Using aria-activedescendant.
- In a single-select listbox, moving focus may optionally unselect the previously selected option and select the newly focused option. This model of selection is known as "selection follows focus". Having selection follow focus can be very helpful in some circumstances and can severely degrade accessibility in others. For additional guidance, see Deciding When to Make Selection Automatically Follow Focus.
- If selecting or unselecting all options is an important function, implementing separate controls for these actions, such as buttons for "Select All" and "Unselect All", significantly improves accessibility.
-
 If the options in a listbox are arranged horizontally:
 - `Down Arrow`performs as- `Right Arrow`is described above, and vice versa.
- `Up Arrow`performs as- `Left Arrow`is described above, and vice versa.

## WAI-ARIA Roles, States, and Properties

- An element that contains or owns all the listbox options has role listbox.
-
 Each option in the listbox has role option and is contained in or owned by either:
 - The element with role `listbox`.
- An element with role group that is contained in or owned by the element with role `listbox`.

- The element with role
-
 Options contained in a `group`are referred to asgrouped options and form what is called anoption group. If a listbox contains grouped options, then:- All option groups contain at least one option.
- Each option group has an accessible name provided via aria-label or aria-labelledby.

- If the element with role `listbox`is not part of another widget, such as a combobox, then it has either a visible label referenced by aria-labelledby or a value specified for aria-label.
-
 If the listbox supports selection of more than one option, the element with role `listbox`has aria-multiselectable set to`true`. Otherwise,`aria-multiselectable`is either set to`false`or the default value of`false`is implied.
-
 The selection state of each selectable option is indicated with either aria-selected or aria-checked:
 -
 If the selection state is indicated with `aria-selected`, then`aria-checked`is not specified for any options. Alternatively, if the selection state is indicated with`aria-checked`, then`aria-selected`is not specified for any options. See notes below regarding considerations for which property to use and for details of the unusual conditions that might allow for both properties in the same listbox.
-
 If any options are selected, each selected option has either aria-selected or aria-checked set to `true`. No more than one option is selected at a time if the element with role`listbox`does*not*have aria-multiselectable set to`true`.
- All options that are selectable but not selected have either aria-selected or aria-checked set to `false`.
- Note that except in listboxes where selection follows focus, the selected state is distinct from focus. For more details, see this description of differences between focus and selection and Deciding When to Make Selection Automatically Follow Focus.

-
 If the selection state is indicated with
- If the complete set of available options is not present in the DOM due to dynamic loading as the user scrolls, their aria-setsize and aria-posinset attributes are set appropriately.
-
 If options are arranged horizontally, the element with role `listbox`has aria-orientation set to`horizontal`. The default value of`aria-orientation`for`listbox`is`vertical`.

### Note

-
 Some factors to consider when choosing whether to indicate selection with `aria-selected`or`aria-checked`are:-
 Some design systems use `aria-selected`for single-select widgets and`aria-checked`for multi-select widgets. In the absence of factors that would make an alternative convention more appropriate, this is a recommended convention.
-
 The language of instructions and the appearance of the interface might suggest one attribute is more appropriate than the other.
 For instance, do instructions say to select items? Or, is the visual indicator of selection a check mark?
- It is important to adopt a consistent convention for selection models across a site or app.

-
 Some design systems use
-
 Conditions that would permit a listbox to include both `aria-selected`and`aria-checked`are extremely rare. It is strongly recommended to avoid designing a listbox widget that would have the need for more than one type of state. If both states were to be used within a listbox, all the following conditions need to be satisfied:- The meaning and purpose of `aria-selected`is different from the meaning and purpose of`aria-checked`in the user interface.
- The user interface makes the meaning and purpose of each state apparent.
- The user interface provides a separate method for controlling each state.

- The meaning and purpose of
- If aria-owns is set on the listbox element to include elements that are not DOM children of the container, those elements will appear in the reading order in the sequence they are referenced and after any items that are DOM children. Scripts that manage focus need to ensure the visual focus order matches this assistive technology reading order.

# Menu and Menubar Pattern

## About This Pattern

 A menu is a widget that offers a list of choices to the user, such as a set of actions or functions.
 Menu widgets behave like native operating system menus, such as the menus that pull down from the menubars commonly found at the top of many desktop application windows.
 A menu is usually opened, or made visible, by activating a menu button, choosing an item in a menu that opens a sub menu, or by invoking a command, such as `Shift` + `F10` in Windows, that opens a context specific menu.
 When a user activates a choice in a menu, the menu usually closes unless the choice opened a submenu.

A menu that is visually persistent is a menubar. A menubar is typically horizontal and is often used to create a menu bar similar to those found near the top of the window in many desktop applications, offering the user quick access to a consistent set of commands.

A common convention for indicating that a menu item launches a dialog box is to append "…" (ellipsis) to the menu item label, e.g., "Save as …".

## Examples

- Editor Menubar Example: Demonstrates menu radios and menu checkboxes in submenus of a menubar that provides text formatting commands for a text field.
- Navigation Menubar Example: Demonstrates a menubar that provides site navigation.

## Keyboard Interaction

The following description of keyboard behaviors assumes:

- A horizontal `menubar`containing several`menuitem`,`menuitemradio`, or`menuitemcheckbox`elements.
- Some `menuitem`elements in the`menubar`have child submenus that contain vertically arranged items.
- Some of the `menuitem`elements in the submenus have child submenus with items that are also vertically arranged.

When reading the following descriptions, also keep in mind that:

- Focusable elements, which may have role `menuitem`,`menuitemradio`, or`menuitemcheckbox`, are referred to as items.
- If a behavior applies to only certain types of items, e.g., `menuitem`elements, the specific role name is used.
- Submenus, also known as popup menus, are elements with role `menu`.
- Except where noted, menus opened from a menubutton behave the same as menus opened from a menubar.

 When a `menu` opens, or when a `menubar` receives focus, keyboard focus is placed on the first item.
 Because `menubar` and `menu` elements are composite widgets as described in the practice for
 Keyboard Navigation Inside Components,
 `Tab` and `Shift` + `Tab` do not move focus among the items in the menu.
 Instead, the keyboard commands described in this section enable users to move focus among the elements in a `menubar` or `menu`.

-
 `Tab`and`Shift`+`Tab`:-
 Move focus into a `menubar`:- If focus is moving into the `menubar`for the first time, focus is set on the first`menuitem`.
-
 If the `menubar`has previously contained focus, focus is optionally set on the`menuitem`that last had focus. Otherwise, it is set on the first`menuitem`that is not disabled.

- If focus is moving into the
- When focus is on a `menuitem`in a`menu`or`menubar`, move focus out of the`menu`or`menubar`, and close all menus and submenus.
-
 Note that `Tab`and`Shift`+`Tab`do not move focus into a`menu`. Unlike a`menubar`, a`menu`is not visually persistent, and authors are responsible for ensuring focus moves to an item inside of a`menu`when the`menu`opens.

-
 Move focus into a
-
 `Enter`:- When focus is on a `menuitem`that has a submenu, opens the submenu and places focus on its first item.
- Otherwise, activates the item and closes the menu.

- When focus is on a
-
 `Space`:- (Optional): When focus is on a `menuitemcheckbox`, changes the state without closing the menu.
- (Optional): When focus is on a `menuitemradio`that is not checked, without closing the menu, checks the focused`menuitemradio`and unchecks any other checked`menuitemradio`element in the same group.
- (Optional): When focus is on a `menuitem`that has a submenu, opens the submenu and places focus on its first item.
- (Optional): When focus is on a `menuitem`that does not have a submenu, activates the`menuitem`and closes the menu.

- (Optional): When focus is on a
-
 `Down Arrow`:- When focus is on a `menuitem`in a`menubar`, and that`menuitem`has a submenu, opens the submenu and places focus on the first item in the submenu.
- When focus is in a `menu`, moves focus to the next item, optionally wrapping from the last to the first.

- When focus is on a
-
 `Up Arrow`:- When focus is in a `menu`, moves focus to the previous item, optionally wrapping from the first to the last.
- (Optional): When focus is on a `menuitem`in a`menubar`, and that`menuitem`has a submenu, opens the submenu and places focus on the last item in the submenu.

- When focus is in a
-
 `Right Arrow`:- When focus is in a `menubar`, moves focus to the next item, optionally wrapping from the last to the first.
- When focus is in a `menu`and on a`menuitem`that has a submenu, opens the submenu and places focus on its first item.
-
 When focus is in a `menu`and on an item that does not have a submenu, performs the following 3 actions:- Closes the submenu and any parent menus.
- Moves focus to the next item in the `menubar`.
- If focus is now on a `menuitem`with a submenu, either: (Recommended) opens the submenu of that`menuitem`without moving focus into the submenu, or opens the submenu of that`menuitem`and places focus on the first item in the submenu.
 `menubar`were not present, e.g., the menus were opened from a menubutton,`Right Arrow`would not do anything when focus is on an item that does not have a submenu.

- When focus is in a
-
 `Left Arrow`:- When focus is in a `menubar`, moves focus to the previous item, optionally wrapping from the first to the last.
- When focus is in a submenu of an item in a `menu`, closes the submenu and returns focus to the parent`menuitem`.
-
 When focus is in a submenu of an item in a `menubar`, performs the following 3 actions:- Closes the submenu.
- Moves focus to the previous item in the `menubar`.
- If focus is now on a `menuitem`with a submenu, either: (Recommended) opens the submenu of that`menuitem`without moving focus into the submenu, or opens the submenu of that`menuitem`and places focus on the first item in the submenu.

- When focus is in a
- `Home`: If arrow key wrapping is not supported, moves focus to the first item in the current- `menu`or- `menubar`.
- `End`: If arrow key wrapping is not supported, moves focus to the last item in the current- `menu`or- `menubar`.
- Any key that corresponds to a printable character (Optional): Move focus to the next item in the current menu whose label begins with that printable character.
- `Escape`: Close the menu that contains focus and return focus to the element or context, e.g., menu button or parent- `menuitem`, from which the menu was opened.

### Note

- Disabled menu items are focusable but cannot be activated.
- A separator in a menu is not focusable or interactive.
-
 If a menu is opened or a menubar receives focus as a result of a context action, `Escape`or`Enter`may return focus to the invoking context. For example, a rich text editor may have a menubar that receives focus when a shortcut key, e.g.,`alt`+`F10`, is pressed while editing. In this case, pressing`Escape`or activating a command from the menu may return focus to the editor.
-
 Although it is recommended that authors avoid doing so, some implementations of navigation menubars may have `menuitem`elements that both perform a function and open a submenu. In such implementations,`Enter`and`Space`perform a navigation function, e.g., load new content, while`Down Arrow`, in a horizontal menubar, opens the submenu associated with that same`menuitem`.
-
 When items in a `menubar`are arranged vertically and items in`menu`containers are arranged horizontally:- `Down Arrow`performs as- `Right Arrow`is described above, and vice versa.
- `Up Arrow`performs as- `Left Arrow`is described above, and vice versa.

## WAI-ARIA Roles, States, and Properties

- A menu is a container of items that represent choices. The element serving as the menu has a role of either menu or menubar.
- The items contained in a menu are child elements of the containing menu or menubar and have any of the following roles:
-
 If activating a menuitem opens a submenu, the menuitem is known as a parent menuitem.
 A submenu's `menu`element is:- Contained inside the same `menu`element as its parent`menuitem`.
- Is the sibling element immediately following its parent `menuitem`.

- Contained inside the same
- A parent menuitem has aria-haspopup set to either `menu`or`true`.
- A parent menuitem has aria-expanded set to `false`when its child menu is not visible and set to`true`when the child menu is visible.
-
 One of the following approaches is used to enable scripts to move focus among items in a menu as described in the practice for
 Keyboard Navigation Inside Components:
 - The menu container has `tabindex`set to`-1`or`0`and aria-activedescendant set to the ID of the focused item.
- Each item in the menu has `tabindex`set to`-1`, except in a menubar, where the first item has`tabindex`set to`0`.

- The menu container has
- When a menuitemcheckbox or menuitemradio is checked, aria-checked is set to `true`.
- When a menu item is disabled, aria-disabled is set to `true`.
- Items in a menu may be divided into groups by placing an element with a role of separator between groups. For example, this technique should be used when a menu contains a set of menuitemradio items.
- All separators should have aria-orientation consistent with the separator's orientation.
-
 If a menubar has a visible label, the element with role `menubar`has aria-labelledby set to a value that refers to the labelling element. Otherwise, the menubar element has a label provided by aria-label.
-
 If a menubar is vertically oriented, it has aria-orientation set to `vertical`. The default value of`aria-orientation`for a menubar is`horizontal`.
-
 An element with role `menu`either has:- aria-labelledby set to a value that refers to the menuitem or button that controls its display.
- A label provided by aria-label.

-
 If a menu is horizontally oriented, it has aria-orientation set to `horizontal`. The default value of`aria-orientation`for a menu is`vertical`.

### Note

If aria-owns is set on the menu container to include elements that are not DOM children of the container, those elements will appear in the reading order in the sequence they are referenced and after any items that are DOM children. Scripts that manage focus need to ensure the visual focus order matches this assistive technology reading order.

# Menu Button Pattern

## About This Pattern

A menu button is a button that opens a menu as described in the Menu and Menubar Pattern. It is often styled as a typical push button with a downward pointing arrow or triangle to hint that activating the button will display a menu.

## Examples

- Action Menu Button Example Using aria-activedescendant: A button that opens a menu of actions or commands where focus in the menu is managed using aria-activedescendant.
- Action Menu Button Example Using element.focus(): A menu button made from an HTML `button`element that opens a menu of actions or commands where focus in the menu is managed using`element.focus()`.
- Navigation Menu Button: A menu button made from an HTML `a`element that opens a menu of items that behave as links.

## Keyboard Interaction

-
 With focus on the button:
 - `Enter`: opens the menu and places focus on the first menu item.
- `Space`: Opens the menu and places focus on the first menu item.
- (Optional) `Down Arrow`: opens the menu and moves focus to the first menu item.
- (Optional) `Up Arrow`: opens the menu and moves focus to the last menu item.

- The keyboard behaviors needed after the menu is open are described in the Menu and Menubar Pattern.

## WAI-ARIA Roles, States, and Properties

- The element that opens the menu has role button.
- The element with role `button`has aria-haspopup set to either`menu`or`true`.
-
 When the menu is displayed, the element with role `button`has aria-expanded set to`true`. When the menu is hidden,`aria-expanded`is set to`false`.
- The element that contains the menu items displayed by activating the button has role menu.
-
 Optionally, the element with role `button`has a value specified for aria-controls that refers to the element with role`menu`.
- Additional roles, states, and properties needed for the menu element are described in the Menu and Menubar Pattern.

# Meter Pattern

## About This Pattern

A meter is a graphical display of a numeric value that varies within a defined range. For example, a meter could be used to depict a device's current battery percentage or a car's fuel level.

### Note

- A `meter`should not be used to represent a value like the current world population since it does not have a meaningful maximum limit.
-
 The `meter`should not be used to indicate progress, such as loading or percent completion of a task. To communicate progress, use the progressbar role instead.

## Example

## Keyboard Interaction

Not applicable.

## WAI-ARIA Roles, States, and Properties

- The element serving as the meter has a role of meter.
- The meter has aria-valuenow set to a decimal value between `aria-valuemin`and`aria-valuemax`representing the current value of the meter.
- The meter has aria-valuemin set to a decimal value less than `aria-valuemax`.
- The meter has aria-valuemax set to a decimal value greater than `aria-valuemin`.
-
 Assistive technologies often present `aria-valuenow`as a percentage. If conveying the value of the meter only in terms of a percentage would not be user friendly, the aria-valuetext property is set to a string that makes the meter value understandable. For example, a battery meter value might be conveyed as`aria-valuetext="50% (6 hours) remaining"`.
-
 If the meter has a visible label, it is referenced by aria-labelledby on the element with role `meter`. Otherwise, the element with role`meter`has a label provided by aria-label.

# Radio Group Pattern

## About This Pattern

A radio group is a set of checkable buttons, known as radio buttons, where no more than one of the buttons can be checked at a time. Some implementations may initialize the set with all buttons in the unchecked state in order to force the user to check one of the buttons before moving past a certain point in the workflow.

## Examples

- Radio Group Example Using Roving tabindex
- Radio Group Example Using aria-activedescendant
- Rating Radio Group Example: Radio group that provides input for a five-star rating scale.

## Keyboard Interaction

### For Radio Groups Not Contained in a Toolbar

This section describes the keyboard interaction implemented for most radio groups. For the special case of a radio group nested inside a toolbar, use the keyboard interaction described in the following section.

-
 `Tab`and`Shift`+`Tab`: Move focus into and out of the radio group. When focus moves into a radio group:- If a radio button is checked, focus is set on the checked button.
- If none of the radio buttons are checked, focus is set on the first radio button in the group.

- `Space`: checks the focused radio button if it is not already checked.
-
 `Right Arrow`and`Down Arrow`: move focus to the next radio button in the group, uncheck the previously focused button, and check the newly focused button. If focus is on the last button, focus moves to the first button.
-
 `Left Arrow`and`Up Arrow`: move focus to the previous radio button in the group, uncheck the previously focused button, and check the newly focused button. If focus is on the first button, focus moves to the last button.

### Note

The initial focus behavior described above differs slightly from the behavior provided by some browsers for native HTML radio groups.
 In some browsers, if none of the radio buttons are selected, moving focus into the radio group with `Shift` + `Tab` will place focus on the last radio button instead of the first radio button.

### For Radio Group Contained in a Toolbar

 Because arrow keys are used to navigate among elements of a toolbar and the `Tab` key moves focus in and out of a toolbar, when a radio group is nested inside a toolbar, the keyboard interaction of the radio group is slightly different from that of a radio group that is not inside of a toolbar.
 For instance, users need to be able to navigate among all toolbar elements, including the radio buttons, without changing which radio button is checked.
 So, when navigating through a radio group in a toolbar with arrow keys, the button that is checked does not change.
 The keyboard interaction of a radio group nested in a toolbar is as follows.

-
 `Space`: If the focused radio button is not checked, unchecks the currently checked radio button and checks the focused radio button. Otherwise, does nothing.
-
 `Enter`(optional): If the focused radio button is not checked, unchecks the currently checked radio button and checks the focused radio button. Otherwise, does nothing.
-
 `Right Arrow`:- When focus is on a radio button and that radio button is **not**the last radio button in the radio group, moves focus to the next radio button.
- When focus is on the last radio button in the radio group and that radio button is **not**the last element in the toolbar, moves focus to the next element in the toolbar.
- When focus is on the last radio button in the radio group and that radio button is also the last element in the toolbar, moves focus to the first element in the toolbar.

- When focus is on a radio button and that radio button is
-
 `Left Arrow`:- When focus is on a radio button and that radio button is **not**the first radio button in the radio group, moves focus to the previous radio button.
- When focus is on the first radio button in the radio group and that radio button is **not**the first element in the toolbar, moves focus to the previous element in the toolbar.
- When focus is on the first radio button in the radio group and that radio button is also the first element in the toolbar, moves focus to the last element in the toolbar.

- When focus is on a radio button and that radio button is
-
 `Down Arrow`(optional): Moves focus to the next radio button in the radio group. If focus is on the last radio button in the radio group, moves focus to the first radio button in the group.
-
 `Up Arrow`(optional): Moves focus to the previous radio button in the radio group. If focus is on the first radio button in the radio group, moves focus to the last radio button in the group.

#### Note

Radio buttons in a toolbar are frequently styled in a manner that appears more like toggle buttons. For an example, See the Simple Editor Toolbar Example.

## WAI-ARIA Roles, States, and Properties

- The radio buttons are contained in or owned by an element with role radiogroup.
- Each radio button element has role radio.
-
 If a radio button is checked, the `radio`element has aria-checked set to`true`. If it is not checked, it has aria-checked set to`false`.
- Each `radio`element is labelled by its content, has a visible label referenced by aria-labelledby, or has a label specified with aria-label.
- The `radiogroup`element has a visible label referenced by aria-labelledby or has a label specified with aria-label.
- If elements providing additional information about either the radio group or each radio button are present, those elements are referenced by the `radiogroup`element or`radio`elements with the aria-describedby property.

# Slider Pattern

## About This Pattern

A slider is an input where the user selects a value from within a given range. Sliders typically have a slider thumb that can be moved along a bar, rail, or track to change the value of the slider.

### Warning

Some users of touch-based assistive technologies may experience difficulty utilizing widgets that implement this slider pattern because the gestures their assistive technology provides for operating sliders may not yet generate the necessary output. To change the slider value, touch-based assistive technologies need to respond to user gestures for increasing and decreasing the value by synthesizing key events. This is a new convention that may not be fully implemented by some assistive technologies. Authors should fully test slider widgets using assistive technologies on devices where touch is a primary input mechanism before considering incorporation into production systems.

## Examples

- Color Viewer Slider Example: Basic horizontal sliders that illustrate setting numeric values for a color picker.
- Vertical Temperature Slider Example: Demonstrates using `aria-orientation`to specify vertical orientation and`aria-valuetext`to communicate unit of measure for a temperature input.
- Rating Slider Example: Horizontal slider that demonstrates using `aria-valuetext`to make it easy for assistive technology users to understand the meaning of the current value chosen on a ten-point satisfaction scale.
- Media Seek Slider Example: Horizontal slider that demonstrates using `aria-valuetext`to communicate current and maximum values of time in media to make the values easy to understand for assistive technology users by converting the total number of seconds to minutes and seconds.

## Keyboard Interaction

- `Right Arrow`: Increase the value of the slider by one step.
- `Up Arrow`: Increase the value of the slider by one step.
- `Left Arrow`: Decrease the value of the slider by one step.
- `Down Arrow`: Decrease the value of the slider by one step.
- `Home`: Set the slider to the first allowed value in its range.
- `End`: Set the slider to the last allowed value in its range.
- `Page Up`(Optional): Increase the slider value by an amount larger than the step change made by- `Up Arrow`.
- `Page Down`(Optional): Decrease the slider value by an amount larger than the step change made by- `Down Arrow`.

### Note

- Focus is placed on the slider (the visual object that the mouse user would move, also known as the thumb.
- In some circumstances, reversing the direction of the value change for the keys specified above, e.g., having `Up Arrow`decrease the value, could create a more intuitive experience.

## WAI-ARIA Roles, States, and Properties

- The element serving as the focusable slider control has role slider.
- The slider element has the aria-valuenow property set to a decimal value representing the current value of the slider.
- The slider element has the aria-valuemin property set to a decimal value representing the minimum allowed value of the slider.
- The slider element has the aria-valuemax property set to a decimal value representing the maximum allowed value of the slider.
- If the value of `aria-valuenow`is not user-friendly, e.g., the day of the week is represented by a number, the aria-valuetext property is set to a string that makes the slider value understandable, e.g., "Monday".
- If the slider has a visible label, it is referenced by aria-labelledby on the slider element. Otherwise, the slider element has a label provided by aria-label.
-
 If the slider is vertically oriented, it has aria-orientation set to `vertical`. The default value of`aria-orientation`for a slider is`horizontal`.

# Slider (Multi-Thumb) Pattern

## About This Pattern

A multi-thumb slider implements the Slider Pattern but includes two or more thumbs, often on a single rail. Each thumb sets one of the values in a group of related values. For example, in a product search, a two-thumb slider could be used to enable users to set the minimum and maximum price limits for the search. In many two-thumb sliders, the thumbs are not allowed to pass one another, such as when the slider sets the minimum and maximum values for a range. For example, in a price range selector, the maximum value of the thumb that sets the lower end of the range is limited by the current value of the thumb that sets the upper end of the range. Conversely, the minimum value of the upper end thumb is limited by the current value of the lower end thumb. However, in some multi-thumb sliders, each thumb sets a value that does not depend on the other thumb values.

### Warning

Some users of touch-based assistive technologies may experience difficulty utilizing widgets that implement this slider pattern because the gestures their assistive technology provides for operating sliders may not yet generate the necessary output. To change the slider value, touch-based assistive technologies need to respond to user gestures for increasing and decreasing the value by synthesizing key events. This is a new convention that may not be fully implemented by some assistive technologies. Authors should fully test slider widgets using assistive technologies on devices where touch is a primary input mechanism before considering incorporation into production systems.

## Example

Horizontal Multi-Thumb Slider Example: Demonstrates a two-thumb slider for picking a price range for a hotel reservation.

## Keyboard Interaction

- Each thumb is in the page tab sequence and has the keyboard interactions described in the Slider Pattern.
- The tab order remains constant regardless of thumb value and visual position within the slider. For example, if the value of a thumb changes such that it moves past one of the other thumbs, the tab order does not change.

## WAI-ARIA Roles, States, and Properties

- Each element serving as a focusable slider thumb has role slider.
- Each slider element has the aria-valuenow property set to a decimal value representing the current value of the slider.
- Each slider element has the aria-valuemin property set to a decimal value representing the minimum allowed value of the slider.
- Each slider element has the aria-valuemax property set to a decimal value representing the maximum allowed value of the slider.
- When the range (e.g. minimum and/or maximum value) of another slider is dependent on the current value of a slider, the values of aria-valuemin or aria-valuemax of the dependent sliders are updated when the value changes.
-
 If a value of `aria-valuenow`is not user-friendly, e.g., the day of the week is represented by a number, the aria-valuetext property is set to a string that makes the slider value understandable, e.g., "Monday".
- If a slider has a visible label, it is referenced by aria-labelledby on the slider element. Otherwise, the slider element has a label provided by aria-label.
-
 If a slider is vertically oriented, it has aria-orientation set to `vertical`. The default value of`aria-orientation`for a slider is`horizontal`.

# Spinbutton Pattern

## About This Pattern

A spinbutton is an input widget that restricts its value to a set or range of discrete values. For example, in a widget that enables users to set an alarm, a spinbutton could allow users to select a number from 0 to 59 for the minute of an hour.

Spinbuttons often have three components, including a text field that displays the current value, an increase button, and a decrease button. The text field is usually the only focusable component because the increase and decrease functions are keyboard accessible via arrow keys. Typically, the text field also allows users to directly edit the value.

 If the range is large, a spinbutton may support changing the value in both small and large steps.
 For instance, in the alarm example, the user may be able to move by 1 minute with `Up Arrow` and `Down Arrow` and by 10 minutes with `Page Up` and `Page Down`.

## Example

Quantity Spin Button Example: A set of three spin buttons to collect the quantities of adults, children, and animals for a hotel reservation.

## Keyboard Interaction

- `Up Arrow`: Increases the value.
- `Down Arrow`: Decreases the value.
- `Home`: If the spinbutton has a minimum value, sets the value to its minimum.
- `End`: If the spinbutton has a maximum value, sets the value to its maximum.
- `Page Up`(Optional): Increases the value by a larger step than- `Up Arrow`.
- `Page Down`(Optional): Decreases the value by a larger step than- `Down Arrow`.
-
 If the spinbutton text field allows directly editing the value, the following keys are supported:
 - Standard single line text editing keys appropriate for the device platform (see note below).
- Printable Characters: Type characters in the textbox. Note that many implementations allow only certain characters as part of the value and prevent input of any other characters. For example, an hour-and-minute spinner would allow only integer values from 0 to 59, the colon ':', and the letters 'AM' and 'PM'. Any other character input does not change the contents of the text field nor the value of the spinbutton.

### Note

- Focus remains on the text field during operation.
-
 Standard single line text editing keys appropriate for the device platform:
 - include keys for input, cursor movement, selection, and text manipulation.
- Standard key assignments for editing functions depend on the device operating system.
- The most robust approach for providing text editing functions is to rely on browsers, which supply them for HTML inputs with type text and for elements with the `contenteditable`HTML attribute.
- **IMPORTANT:**Be sure that JavaScript does not interfere with browser-provided text editing functions by capturing key events for the keys used to perform them.

## WAI-ARIA Roles, States, and Properties

- The focusable element serving as the spinbutton has role spinbutton. This is typically an element that supports text input.
- The spinbutton element has the aria-valuenow property set to a decimal value representing the current value of the spinbutton.
- The spinbutton element has the aria-valuemin property set to a decimal value representing the minimum allowed value of the spinbutton if it has a known minimum value.
- The spinbutton element has the aria-valuemax property set to a decimal value representing the maximum allowed value of the spinbutton if it has a known maximum value.
-
 If the value of `aria-valuenow`is not user-friendly, e.g., the day of the week is represented by a number, the aria-valuetext property is set on the spinbutton element to a string that makes the spinbutton value understandable, e.g., "Monday".
- If the spinbutton has a visible label, it is referenced by aria-labelledby on the spinbutton element. Otherwise, the spinbutton element has a label provided by aria-label.
-
 The spinbutton element has aria-invalid set to `true`if the value is outside the allowed range. Note that most implementations prevent input of invalid values, but in some scenarios, blocking all invalid input may not be practical.

# Switch Pattern

## About This Pattern

 A switch is an input widget that allows users to choose one of two values: on

 or off

.
 Switches are similar to checkboxes and toggle buttons, which can also serve as binary inputs.
 One difference, however, is that switches can only be used for binary input while checkboxes and toggle buttons allow implementations the option of supporting a third middle state.
 Checkboxes can be checked

 or not checked

 and can optionally also allow for a partially checked

 state.
 Toggle buttons can be pressed

 or not pressed

 and can optionally allow for a partially pressed

 state.

 Since switch, checkbox, and toggle button all offer binary input, they are often functionally interchangeable.
 Choose the role that best matches both the visual design and semantics of the user interface.
 For instance, there are some circumstances where the semantics of on or off

 would be easier for assistive technology users to understand than the semantics of checked or unchecked

, and vice versa.
 Consider a widget for turning lights on or off.
 In this case, screen reader output of Lights switch on

 is more user friendly than Lights checkbox checked

.
 However, if the same input were in a group of inputs labeled Which of the following must be included in your pre-takeoff procedures?

, Lights checkbox checked

 would make more sense.

**Important:** it is critical the label on a switch does not change when its state changes.

## Examples

- Switch Example: A switch based on a `div`element that turns a notification preference on and off.
- Switch Example Using HTML Button: A group of 2 switches based on HTML `button`elements that turn lights on and off.
- Switch Example Using HTML Checkbox Input: A group of 2 switches based on HTML `input[type="checkbox"]`elements that turn accessibility preferences on and off.

## Keyboard Interaction

- `Space`: When focus is on the switch, changes the state of the switch.
- `Enter`(Optional): When focus is on the switch, changes the state of the switch.

## WAI-ARIA Roles, States, and Properties

- The switch has role switch.
-
 The switch has an accessible label provided by one of the following:
 - Visible text content contained within the element with role `switch`.
- A visible label referenced by the value of aria-labelledby set on the element with role `switch`.
- aria-label set on the element with role `switch`.

- Visible text content contained within the element with role
- When `on`, the switch element has state aria-checked set to`true`.
- When `off`, the switch element has state aria-checked set to`false`.
- If the switch element is an HTML `input[type="checkbox"]`, it uses the HTML`checked`attribute instead of the`aria-checked`property.
-
 If a set of switches is presented as a logical group with a visible label, either:
 - The switches are included in an element with role group that has the property aria-labelledby set to the ID of the element containing the group label.
- The set is contained in an HTML `fieldset`and the label for the set is contained in an HTML`legend`element.

- If the presentation includes additional descriptive static text relevant to a switch or switch group, the switch or switch group has the property aria-describedby set to the ID of the element containing the description.

# Table Pattern

## About This Pattern

 Like an HTML `table` element, a WAI-ARIA table is a static tabular structure containing one or more rows that each contain one or more cells; it is not an interactive widget.
 Thus, its cells are not focusable or selectable.
 The grid pattern is used to make an interactive widget that has a tabular structure.

However, tables are often used to present a combination of information and interactive widgets. Since a table is not a widget, each widget contained in a table is a separate stop in the page tab sequence. If the number of widgets is large, replacing the table with a grid can dramatically reduce the length of the page tab sequence because a grid is a composite widget that can contain other widgets.

### Note

As with other WAI-ARIA roles that have a native host language equivalent, authors are strongly encouraged to use a native HTML `table` element whenever possible.
 This is especially important with role `table` because it is a new feature of WAI-ARIA 1.1.
 It is thus advisable to test implementations thoroughly with each browser and assistive technology combination that could be used by the target audience.

## Examples

- Table Example: ARIA table made using HTML `div`and`span`elements.
- Sortable Table Example: Basic HTML table that illustrates implementation of `aria-sort`in the headers of sortable columns.

## Keyboard Interaction

Not applicable.

## WAI-ARIA Roles, States, and Properties

- The table container has role table.
- Each row container has role row and is either a DOM descendant of or owned by the `table`element or an element with role rowgroup.
-
 Each cell is either a DOM descendant of or owned by a `row`element and has one of the following roles:- columnheader if the cell contains a title or header information for the column.
- rowheader if the cell contains title or header information for the row.
- cell if the cell does not contain column or row header information.

- If there is an element in the user interface that serves as a label for the table, aria-labelledby is set on the table element with a value that refers to the labelling element. Otherwise, a label is specified for the table element using aria-label.
- If the table has a caption or description, aria-describedby is set on the table element with a value referring to the element containing the description.
- If the table contains sortable columns or rows, aria-sort is set to an appropriate value on the header cell element for the sorted column or row as described in the Grid and Table Properties Practice.
-
 If there are conditions where some rows or columns are hidden or not present in the DOM, e.g., there are widgets on the page for hiding rows or columns, the following properties are applied as described in the
 Grid and Table Properties Practice.
 - aria-colcount or aria-rowcount is set to the total number of columns or rows, respectively.
- aria-colindex or aria-rowindex is set to the position of a cell within a row or column, respectively.

- If the table includes cells that span multiple rows or multiple columns, then aria-rowspan or aria-colspan is applied as described in the Grid and Table Properties Practice.

### Note

If rows or cells are included in a table via aria-owns, they will be presented to assistive technologies after the DOM descendants of the `table` element unless the DOM descendants are also included in the `aria-owns` attribute.

# Tabs Pattern

## About This Pattern

Tabs are a set of layered sections of content, known as tab panels, that display one panel of content at a time. Each tab panel has an associated tab element, that when activated, displays the panel. The list of tab elements is arranged along one edge of the currently displayed panel, most commonly the top edge.

Terms used to describe this design pattern include:

- Tabs or Tabbed Interface
- A set of tab elements and their associated tab panels.
- Tab List
- A set of tab elements contained in a tablist element.
- tab
- An element in the tab list that serves as a label for one of the tab panels and can be activated to display that panel.
- tabpanel
- The element that contains the content associated with a tab.

When a tabbed interface is initialized, one tab panel is displayed and its associated tab is styled to indicate that it is active. When the user activates one of the other tab elements, the previously displayed tab panel is hidden, the tab panel associated with the activated tab becomes visible, and the tab is considered "active".

## Examples

- Tabs With Automatic Activation: A tabs widget where tabs are automatically activated and their panel is displayed when they receive focus.
- Tabs With Manual Activation: A tabs widget where users activate a tab and display its panel by pressing `Space`or`Enter`.

## Keyboard Interaction

For the tab list:

-
 `Tab`:- When focus moves into the tab list, places focus on the active `tab`element.
- When the tab list contains the focus, moves focus to the next element in the page tab sequence outside the tablist, which is the tabpanel unless the first element containing meaningful content inside the tabpanel is focusable.

- When focus moves into the tab list, places focus on the active
-
 When focus is on a tab element in a horizontal tab list:
 -
 `Left Arrow`: moves focus to the previous tab. If focus is on the first tab, moves focus to the last tab. Optionally, activates the newly focused tab (See note below).
-
 `Right Arrow`: Moves focus to the next tab. If focus is on the last tab element, moves focus to the first tab. Optionally, activates the newly focused tab (See note below).

-

-
 When focus is on a tab in a tablist with either horizontal or vertical orientation:
 - `Space or Enter`: Activates the tab if it was not activated automatically on focus.
-
 `Home`(Optional): Moves focus to the first tab. Optionally, activates the newly focused tab (See note below).
-
 `End`(Optional): Moves focus to the last tab. Optionally, activates the newly focused tab (See note below).
- `Shift`+- `F10`: If the tab has an associated popup menu, opens the menu.
-
 `Delete`(Optional): If deletion is allowed, deletes (closes) the current tab element and its associated tab panel, sets focus on the tab following the tab that was closed, and optionally activates the newly focused tab. If there is not a tab that followed the tab that was deleted, e.g., the deleted tab was the right-most tab in a left-to-right horizontal tab list, sets focus on and optionally activates the tab that preceded the deleted tab. If the application allows all tabs to be deleted, and the user deletes the last remaining tab in the tab list, the application moves focus to another element that provides a logical work flow. As an alternative to`Delete`, or in addition to supporting`Delete`, the delete function is available in a context menu.

### Note

- It is recommended that tabs activate automatically when they receive focus as long as their associated tab panels are displayed without noticeable latency. This typically requires tab panel content to be preloaded. Otherwise, automatic activation slows focus movement, which significantly hampers users' ability to navigate efficiently across the tab list. For additional guidance, see Deciding When to Make Selection Automatically Follow Focus.
-
 When a tab list has its aria-orientation set to `vertical`:- `Down Arrow`performs as- `Right Arrow`is described above.
- `Up Arrow`performs as- `Left Arrow`is described above.

- If the tab list is horizontal, it does not listen for `Down Arrow`or`Up Arrow`so those keys can provide their normal browser scrolling functions even when focus is inside the tab list.
- When the tabpanel does not contain any focusable elements or the first element with content is not focusable, the tabpanel should set `tabindex="0"`to include it in the tab sequence of the page.

## WAI-ARIA Roles, States, and Properties

- The element that serves as the container for the set of tabs has role tablist.
- Each element that serves as a tab has role tab and is contained within the element with role `tablist`.
- Each element that contains the content panel for a `tab`has role tabpanel.
-
 If the tab list has a visible label, the element with role `tablist`has aria-labelledby set to a value that refers to the labelling element. Otherwise, the`tablist`element has a label provided by aria-label.
- Each element with role `tab`has the property aria-controls referring to its associated`tabpanel`element.
- The active `tab`element has the state aria-selected set to`true`and all other`tab`elements have it set to`false`.
- Each element with role `tabpanel`has the property aria-labelledby referring to its associated`tab`element.
- If a `tab`element has a popup menu, it has the property aria-haspopup set to either`menu`or`true`.
-
 If the `tablist`element is vertically oriented, it has the property aria-orientation set to`vertical`. The default value of`aria-orientation`for a`tablist`element is`horizontal`.

# Toolbar Pattern

## About This Pattern

A toolbar is a container for grouping a set of controls, such as buttons, menubuttons, or checkboxes.

 When a set of controls is visually presented as a group, the `toolbar` role can be used to communicate the presence and purpose of the grouping to screen reader users.
 Grouping controls into toolbars can also be an effective way of reducing the number of tab stops in the keyboard interface.

To optimize the benefit of toolbar widgets:

-
 Implement focus management so the keyboard tab sequence includes one stop for the toolbar and arrow keys move focus among the controls in the toolbar.
 -
 In horizontal toolbars, `Left Arrow`and`Right Arrow`navigate among controls.`Up Arrow`and`Down Arrow`can duplicate`Left Arrow`and`Right Arrow`, respectively, or can be reserved for operating controls, such as spin buttons that require vertical arrow keys to operate.
-
 In vertical toolbars, `Up Arrow`and`Down Arrow`navigate among controls.`Left Arrow`and`Right Arrow`can duplicate`Up Arrow`and`Down Arrow`, respectively, or can be reserved for operating controls, such as horizontal sliders that require horizontal arrow keys to operate.
- In toolbars with multiple rows of controls, `Left Arrow`and`Right Arrow`can provide navigation that wraps from row to row, leaving the option of reserving vertical arrow keys for operating controls.

-
 In horizontal toolbars,
- Avoid including controls whose operation requires the pair of arrow keys used for toolbar navigation. If unavoidable, include only one such control and make it the last element in the toolbar. For example, in a horizontal toolbar, a textbox could be included as the last element.
- Use toolbar as a grouping element only if the group contains 3 or more controls.

## Example

Toolbar Example: A toolbar that uses roving tabindex to manage focus and contains several types of controls, including toggle buttons, radio buttons, a menu button, a spin button, a checkbox, and a link.

## Keyboard Interaction

-
 `Tab`and`Shift`+`Tab`: Move focus into and out of the toolbar. When focus moves into a toolbar:- If focus is moving into the toolbar for the first time, focus is set on the first control that is not disabled.
- If the toolbar has previously contained focus, focus is optionally set on the control that last had focus. Otherwise, it is set on the first control that is not disabled.

-
 For a horizontal toolbar (the default):
 -
 `Left Arrow`: Moves focus to the previous control. Optionally, focus movement may wrap from the first element to the last element.
-
 `Right Arrow`: Moves focus to the next control. Optionally, focus movement may wrap from the last element to the first element.

-

- `Home`(Optional): Moves focus to first element.
- `End`(Optional): Moves focus to last element.

### Note

-
 If the items in a toolbar are arranged vertically:
 - `Down Arrow`performs as- `Right Arrow`is described above.
- `Up Arrow`performs as- `Left Arrow`is described above.

- Typically, disabled elements are not focusable when navigating with a keyboard. However, in circumstances where discoverability of a function is crucial, it may be helpful if disabled controls are focusable so screen reader users are more likely to be aware of their presence. For additional guidance, see Focusability of disabled controls.
- In applications where quick access to a toolbar is important, such as accessing an editor's toolbar from its text area, a documented shortcut key for moving focus from the relevant context to its corresponding toolbar is recommended.

## WAI-ARIA Roles, States, and Properties

- The element that serves as the toolbar container has role toolbar.
- If the toolbar has a visible label, it is referenced by aria-labelledby on the toolbar element. Otherwise, the toolbar element has a label provided by aria-label.
-
 If the controls are arranged vertically, the toolbar element has aria-orientation set to `vertical`. The default orientation is horizontal.

# Tooltip Pattern

## About This Pattern

 **NOTE:** This design pattern is work in progress; it does not yet have task force consensus.
 Progress and discussions are captured in issue 128.

 A tooltip is a popup that displays information related to an element when the element receives keyboard focus or the mouse hovers over it.
 It typically appears after a small delay and disappears when `Escape` is pressed or on mouse out.

Tooltip widgets do not receive focus. A hover that contains focusable elements can be made using a non-modal dialog.

## Example

Work to develop a tooltip example is tracked by issue 127.

## Keyboard Interaction

`Escape`: Dismisses the Tooltip.

### Note

- Focus stays on the triggering element while the tooltip is displayed.
- If the tooltip is invoked when the trigger element receives focus, then it is dismissed when it no longer has focus (onBlur).
- If the tooltip is invoked when a pointing cursor moves over the trigger element, then it remains open as long as the cursor is over the trigger or the tooltip.

## WAI-ARIA Roles, States, and Properties

- The element that serves as the tooltip container has role tooltip.
- The element that triggers the tooltip references the tooltip element with aria-describedby.

# Tree View Pattern

## About This Pattern

A tree view widget presents a hierarchical list. Any item in the hierarchy may have child items, and items that have children may be expanded or collapsed to show or hide the children. For example, in a file system navigator that uses a tree view to display folders and files, an item representing a folder can be expanded to reveal the contents of the folder, which may be files, folders, or both.

When using a keyboard to navigate a tree, a visual keyboard indicator informs the user which item is focused. If the tree allows the user to choose just one item for an action, then it is known as a single-select tree. In some implementations of single-select tree, the focused item also has a selected state; this is known as selection follows focus. However, in multi-select trees, which enable the user to select more than one item for an action, the selected state is always independent of the focus. For example, in a typical file system navigator, the user can move focus to select any number of files for an action, such as copy or move. It is important that the visual design distinguish between items that are selected and the item that has focus. For more details, see this description of differences between focus and selection and Deciding When to Make Selection Automatically Follow Focus.

## Examples

- File Directory Treeview Example Using Computed Properties: A file selector tree that demonstrates browser support for automatically computing `aria-level`,`aria-posinset`and`aria-setsize`based on DOM structure.
- File Directory Treeview Example Using Declared Properties: A file selector tree that demonstrates how to explicitly define values for `aria-level`,`aria-posinset`and`aria-setsize`.
- Navigation Treeview Example: A tree that provides navigation to a set of web pages and demonstrates browser support for automatically computing `aria-level`,`aria-posinset`and`aria-setsize`based on DOM structure.

## Terms

Terms for describing tree views include:

- Node
- An item in a tree.
- Root Node
- Node at the base of the tree; it may have one or more child nodes but does not have a parent node.
- Child Node
- Node that has a parent; any node that is not a root node is a child node.
- End Node
- Node that does not have any child nodes; an end node may be either a root node or a child node.
- Parent Node
- Node with one or more child nodes. It can be open (expanded) or closed (collapsed).
- Open Node
- Parent node that is expanded so its child nodes are visible.
- Closed Node
- Parent node that is collapsed so the child nodes are not visible.

## Keyboard Interaction

For a vertically oriented tree:

-
 When a single-select tree receives focus:
 - If none of the nodes are selected before the tree receives focus, focus is set on the first node.
- If a node is selected before the tree receives focus, focus is set on the selected node.

-
 When a multi-select tree receives focus:
 - If none of the nodes are selected before the tree receives focus, focus is set on the first node.
- If one or more nodes are selected before the tree receives focus, focus is set on the first selected node.

-
 `Right arrow`:- When focus is on a closed node, opens the node; focus does not move.
- When focus is on a open node, moves focus to the first child node.
- When focus is on an end node, does nothing.

-
 `Left arrow`:- When focus is on an open node, closes the node.
- When focus is on a child node that is also either an end node or a closed node, moves focus to its parent node.
- When focus is on a root node that is also either an end node or a closed node, does nothing.

- `Down Arrow`: Moves focus to the next node that is focusable without opening or closing a node.
- `Up Arrow`: Moves focus to the previous node that is focusable without opening or closing a node.
- `Home`: Moves focus to the first node in the tree without opening or closing a node.
- `End`: Moves focus to the last node in the tree that is focusable without opening a node.
-
 `Enter`: activates a node, i.e., performs its default action. For parent nodes, one possible default action is to open or close the node. In single-select trees where selection does not follow focus (see note below), the default action is typically to select the focused node.
-
 Type-ahead is recommended for all trees, especially for trees with more than 7 root nodes:
 - Type a character: focus moves to the next node with a name that starts with the typed character.
- Type multiple characters in rapid succession: focus moves to the next node with a name that starts with the string of characters typed.

- `*`(Optional): Expands all siblings that are at the same level as the current node.
-
 **Selection in multi-select trees:**Authors may implement either of two interaction models to support multiple selection: a recommended model that does not require the user to hold a modifier key, such as`Shift`or`Control`, while navigating the list or an alternative model that does require modifier keys to be held while navigating in order to avoid losing selection states.-
 Recommended selection model -- holding a modifier key while moving focus is not necessary:
 - `Space`: Toggles the selection state of the focused node.
- `Shift`+- `Down Arrow`(Optional): Moves focus to and toggles the selection state of the next node.
- `Shift`+- `Up Arrow`(Optional): Moves focus to and toggles the selection state of the previous node.
- `Shift`+- `Space`(Optional): Selects contiguous nodes from the most recently selected node to the current node.
-
 `Control`+`Shift`+`Home`(Optional): Selects the node with focus and all nodes up to the first node. Optionally, moves focus to the first node.
-
 `Control`+`Shift`+`End`(Optional): Selects the node with focus and all nodes down to the last node. Optionally, moves focus to the last node.
-
 `Control`+`A`(Optional): Selects all nodes in the tree. Optionally, if all nodes are selected, it can also unselect all nodes.

-
 Alternative selection model -- Moving focus without holding the `Shift`or`Control`modifier unselects all selected nodes except for the focused node:- `Shift`+- `Down Arrow`: Moves focus to and toggles the selection state of the next node.
- `Shift`+- `Up Arrow`: Moves focus to and toggles the selection state of the previous node.
- `Control`+- `Down Arrow`: Without changing the selection state, moves focus to the next node.
- `Control`+- `Up Arrow`: Without changing the selection state, moves focus to the previous node.
- `Control`+- `Space`: Toggles the selection state of the focused node.
- `Shift`+- `Space`(Optional): Selects contiguous nodes from the most recently selected node to the current node.
-
 `Control`+`Shift`+`Home`(Optional): Selects the node with focus and all nodes up to the first node. Optionally, moves focus to the first node.
-
 `Control`+`Shift`+`End`(Optional): Selects the node with focus and all nodes down to the last node. Optionally, moves focus to the last node.
-
 `Control`+`A`(Optional): Selects all nodes in the tree. Optionally, if all nodes are selected, it can also unselect all nodes.

-
 Recommended selection model -- holding a modifier key while moving focus is not necessary:

### Note

- DOM focus (the active element) is functionally distinct from the selected state. For more details, see this description of differences between focus and selection.
-
 The `tree`role supports the aria-activedescendant property, which provides an alternative to moving DOM focus among`treeitem`elements when implementing keyboard navigation. For details, see Managing Focus in Composites Using aria-activedescendant.
- In a single-select tree, moving focus may optionally unselect the previously selected node and select the newly focused node. This model of selection is known as "selection follows focus". Having selection follow focus can be very helpful in some circumstances and can severely degrade accessibility in others. For additional guidance, see Deciding When to Make Selection Automatically Follow Focus.
- If selecting or unselecting all nodes is an important function, implementing separate controls for these actions, such as buttons for "Select All" and "Unselect All", significantly improves accessibility.
-
 If the nodes in a tree are arranged horizontally:
 - `Down Arrow`performs as- `Right Arrow`is described above, and vice versa.
- `Up Arrow`performs as- `Left Arrow`is described above, and vice versa.

## WAI-ARIA Roles, States, and Properties

- All tree nodes are contained in or owned by an element with role tree.
- Each element serving as a tree node has role treeitem.
- Each root node is contained in the element with role `tree`or referenced by an aria-owns property set on the`tree`element.
- Each parent node contains or owns an element with role group.
- Each child node is contained in or owned by an element with role group that is contained in or owned by the node that serves as the parent of that child.
-
 Each element with role `treeitem`that serves as a parent node has aria-expanded set to`false`when the node is in a closed state and set to`true`when the node is in an open state. End nodes do not have the`aria-expanded`attribute because, if they were to have it, they would be incorrectly described to assistive technologies as parent nodes.
-
 If the tree supports selection of more than one node, the element with role `tree`has aria-multiselectable set to`true`. Otherwise,`aria-multiselectable`is either set to`false`or the default value of`false`is implied.
-
 The selection state of each selectable node is indicated with either aria-selected or aria-checked:
 -
 If the selection state is indicated with `aria-selected`, then`aria-checked`is not specified for any nodes. Alternatively, if the selection state is indicated with`aria-checked`, then`aria-selected`is not specified for any nodes. See notes below regarding considerations for which property to use and for details of the unusual conditions that might allow for both properties in the same tree.
-
 If any nodes are selected, each selected node has either aria-selected or aria-checked set to `true`. No more than one node is selected at a time if the element with role`tree`does*not*have aria-multiselectable set to`true`.
- All nodes that are selectable but not selected have either aria-selected or aria-checked set to `false`.
- If the tree contains nodes that are not selectable, neither aria-selected nor aria-checked is present on those nodes.
- Note that except in trees where selection follows focus, the selected state is distinct from focus. For more details, see this description of differences between focus and selection and Deciding When to Make Selection Automatically Follow Focus.

-
 If the selection state is indicated with
- The element with role `tree`has either a visible label referenced by aria-labelledby or a value specified for aria-label.
- If the complete set of available nodes is not present in the DOM due to dynamic loading as the user moves focus in or scrolls the tree, each node has aria-level, aria-setsize, and aria-posinset specified.
-
 If the `tree`element is horizontally oriented, it has aria-orientation set to`horizontal`. The default value of`aria-orientation`for a tree is`vertical`.

### Note

-
 Some factors to consider when choosing whether to indicate selection with `aria-selected`or`aria-checked`are:-
 Some design systems use `aria-selected`for single-select widgets and`aria-checked`for multi-select widgets. In the absence of factors that would make an alternative convention more appropriate, this is a recommended convention.
-
 The language of instructions and the appearance of the interface might suggest one attribute is more appropriate than the other.
 For instance, do instructions say to select items? Or, is the visual indicator of selection a check mark?
- It is important to adopt a consistent convention for selection models across a site or app.

-
 Some design systems use
-
 Conditions that would permit a tree to include both `aria-selected`and`aria-checked`are extremely rare. It is strongly recommended to avoid designing a tree widget that would have the need for more than one type of state. If both states were to be used within a tree, all the following conditions need to be satisfied:- The meaning and purpose of `aria-selected`is different from the meaning and purpose of`aria-checked`in the user interface.
- The user interface makes the meaning and purpose of each state apparent.
- The user interface provides a separate method for controlling each state.

- The meaning and purpose of
- If aria-owns is set on the tree container to include elements that are not DOM children of the container, those elements will appear in the reading order in the sequence they are referenced and after any items that are DOM children. Scripts that manage focus need to ensure the visual focus order matches this assistive technology reading order.

# Treegrid Pattern

## About This Pattern

 A treegrid widget presents a hierarchical data grid consisting of tabular information that is editable or interactive.
 Any row in the hierarchy may have child rows, and rows with children may be expanded or collapsed to show or hide the children.
 For example, in a `treegrid` used to display messages and message responses for a e-mail discussion list, messages with responses would be in rows that can be expanded to reveal the response messages.

 In a treegrid both rows and cells are focusable.
 Every row and cell contains a focusable element or is itself focusable, regardless of whether individual cell content is editable or interactive.
 There is one exception: if column header cells do not provide functions, such as sort or filter, they do not need to be focusable.
 One reason it is important for all cells to be able to receive or contain keyboard focus is that screen readers will typically be in their application reading mode, rather than their document reading mode, when users are interacting with the grid.
 While in application mode, a screen reader user hears only focusable elements and content that labels focusable elements.
 So, screen reader users may unknowingly overlook elements contained in a `treegrid` that are either not focusable or not used to label a column or row.

 When using a keyboard to navigate a `treegrid`, a visual keyboard indicator informs users which row or cell is focused.
 If the `treegrid` allows users to choose just one item for an action, then it is known as a single-select `treegrid`, and the item with focus also has a selected state.
 However, in a multi-select `treegrid`, which enables users to select more than one row or cell for an action, the selected state is independent of the focus.
 For example, in a hierarchical e-mail discussion grid, users can move focus to select any number of rows for an action, such as delete or move.
 It is important that the visual design distinguish between items that are selected and the item that has focus.
 For more details, see this description of differences between focus and selection.

## Example

E-mail Inbox `treegrid` Example: A treegrid for navigating an e-mail inbox that demonstrates three keyboard navigation models -- rows first, cells first, and cells only.

## Keyboard Interaction

 The following keys provide `treegrid` navigation by moving focus among rows and cells of the grid.
 Implementations of `treegrid` make these key commands available when an element in the grid has received focus, e.g., after a user has moved focus to the grid with `Tab`.
 Moving focus into the grid may result in the first cell or the first row being focused.
 Whether focus goes to a cell or the row depends on author preferences and whether row focus is supported, since some implementations of `treegrid` may not provide focus to rows.

- `Enter`: If cell-only focus is enabled and focus is on the first cell with the- `aria-expanded`property, opens or closes the child rows.,Otherwise, performs the default action for the cell.
-
 `Tab`: If the row containing focus contains focusable elements (e.g., inputs, buttons, links, etc.), moves focus to the next input in the row. If focus is on the last focusable element in the row, moves focus out of the`treegrid`widget to the next focusable element.
-
 `Right Arrow`:- If focus is on a collapsed row, expands the row.
- If focus is on an expanded row or is on a row that does not have child rows, moves focus to the first cell in the row.
- If focus is on the right-most cell in a row, focus does not move.
- If focus is on any other cell, moves focus one cell to the right.

-
 `Left Arrow`:- If focus is on an expanded row, collapses the row.
- If focus is on a collapsed row or on a row that does not have child rows, focus does not move.
- If focus is on the first cell in a row and row focus is supported, moves focus to the row.
- If focus is on the first cell in a row and row focus is not supported, focus does not move.
- If focus is on any other cell, moves focus one cell to the left.

-
 `Down Arrow`:- If focus is on a row, moves focus one row down. If focus is on the last row, focus does not move.
- If focus is on a cell, moves focus one cell down. If focus is on the bottom cell in the column, focus does not move.

-
 `Up Arrow`:- If focus is on a row, moves focus one row up. If focus is on the first row, focus does not move.
- If focus is on a cell, moves focus one cell up. If focus is on the top cell in the column, focus does not move.

-
 `Page Down`:- If focus is on a row, moves focus down an author-determined number of rows, typically scrolling so the bottom row in the currently visible set of rows becomes one of the first visible rows. If focus is in the last row, focus does not move.
- If focus is on a cell, moves focus down an author-determined number of cells, typically scrolling so the bottom row in the currently visible set of rows becomes one of the first visible rows. If focus is in the last row, focus does not move.

-
 `Page Up`:- If focus is on a row, moves focus up an author-determined number of rows, typically scrolling so the top row in the currently visible set of rows becomes one of the last visible rows. If focus is in the first row, focus does not move.
- If focus is on a cell, moves focus up an author-determined number of cells, typically scrolling so the top row in the currently visible set of rows becomes one of the last visible rows. If focus is in the first row, focus does not move.

-
 `Home`:- If focus is on a row, moves focus to the first row. If focus is in the first row, focus does not move.
- If focus is on a cell, moves focus to the first cell in the row. If focus is in the first cell of the row, focus does not move.

-
 `End`:- If focus is on a row, moves focus to the last row. If focus is in the last row, focus does not move.
- If focus is on a cell, moves focus to the last cell in the row. If focus is in the last cell of the row, focus does not move.

-
 `Control`+`Home`:- If focus is on a row, moves focus to the first row. If focus is in the first row, focus does not move.
- If focus is on a cell, moves focus to the cell in the first row in the same column as the cell that had focus. If focus is in the first row, focus does not move.

-
 `Control`+`End`:- If focus is on a row, moves focus to the last row. If focus is in the last row, focus does not move.
- If focus is on a cell, moves focus to the cell in the last row in the same column as the cell that had focus. If focus is in the last row, focus does not move.

### Note

-
 When the above `treegrid`navigation keys move focus, whether the focus is set on an element inside the cell or on the cell depends on cell content. See Whether to Focus on a Cell or an Element Inside It.
- While navigation keys, such as arrow keys, are moving focus from cell to cell, they are not available to do something like operate a combobox or move an editing caret inside of a cell. If this functionality is needed, see Editing and Navigating Inside a Cell.
- If navigation functions can dynamically add more rows or columns to the DOM, key events that move focus to the beginning or end of the grid, such as `control`+`End`, may move focus to the last row in the DOM rather than the last available row in the back-end data.

If a treegrid supports selection of cells, rows, or columns, the following keys are commonly used for these functions.

-
 `Control`+`Space`:- If focus is on a row, selects all cells.
- If focus is on a cell, selects the column that contains the focus.

-
 `Shift`+`Space`:- If focus is on a row, selects the row.
- If focus is on a cell, selects the row that contains the focus. If the treegrid includes a column with checkboxes for selecting rows, this key can serve as a shortcut for checking the box when focus is not on the checkbox.

- `Control`+- `A`: Selects all cells.
-
 `Shift`+`Right Arrow`:- If focus is on a row, does not change selection.
- if focus is on a cell, extends selection one cell to the right.

-
 `Shift`+`Left Arrow`:- If focus is on a row, does not change selection.
- if focus is on a cell, extends selection one cell to the left.

-
 `Shift`+`Down Arrow`:- If focus is on a row, extends selection to all the cells in the next row.
- If focus is on a cell, extends selection one cell down.

-
 `Shift`+`Up Arrow`:- If focus is on a row, extends selection to all the cells in the previous row.
- If focus is on a cell, extends selection one cell up.

### Note

See Key Assignment Conventions for Common Functions for cut, copy, and paste key assignments.

## WAI-ARIA Roles, States, and Properties

- The treegrid container has role treegrid.
- Each row container has role row and is either a DOM descendant of or owned by the `treegrid`element or an element with role rowgroup.
-
 Each cell is either a DOM descendant of or owned by a `row`element and has one of the following roles:- columnheader if the cell contains a title or header information for the column.
- rowheader if the cell contains title or header information for the row.
- gridcell if the cell does not contain column or row header information.

-
 A `row`that can be expanded or collapsed to show or hide a set of child rows is a parent row. Each parent`row`has the aria-expanded state set on either the`row`element or on a cell contained in the`row`. The`aria-expanded`state is set to`false`when the child rows are not displayed and set to`true`when the child rows are displayed. Rows that do not control display of child rows do not have the`aria-expanded`attribute because, if they were to have it, they would be incorrectly described to assistive technologies as parent rows.
-
 If the treegrid supports selection of more than one row or cell, it is a multi-select treegrid and the element with role `treegrid`has aria-multiselectable set to`true`. Otherwise, it is a single-select treegrid, and`aria-multiselectable`is either set to`false`or the default value of`false`is implied.
- If the treegrid is a single-select treegrid, aria-selected is set to `true`on the selected row or cell, and it is not present on any other row or cell in the treegrid.
-
 if the treegrid is a multi-select treegrid:
 - All selected rows or cells have aria-selected set to `true`.
- All rows and cells that are not selected have aria-selected set to `false`.

- All selected rows or cells have aria-selected set to
- If there is an element in the user interface that serves as a label for the treegrid, aria-labelledby is set on the grid element with a value that refers to the labelling element. Otherwise, a label is specified for the grid element using aria-label.
- If the treegrid has a caption or description, aria-describedby is set on the grid element with a value referring to the element containing the description.
- If the treegrid provides sort functions, aria-sort is set to an appropriate value on the header cell element for the sorted column or row as described in the Grid and Table Properties Practice.
-
 If the treegrid provides content editing functionality and contains cells that may have edit capabilities disabled in certain conditions, aria-readonly is set to `true`on cells where editing is disabled. If edit functions are disabled for all cells, instead of setting`aria-readonly`to`true`on every cell,`aria-readonly`may be set to`true`on the`treegrid`element. Treegrids that do not provide cell content editing functions do not include the`aria-readonly`attribute on any of their elements.
-
 If there are conditions where some rows or columns are hidden or not present in the DOM, e.g., data is dynamically loaded when scrolling or the treegrid provides functions for hiding rows or columns, the following properties are applied as described in the
 Grid and Table Properties Practice.
 - aria-colcount or aria-rowcount is set to the total number of columns or rows, respectively.
- aria-colindex or aria-rowindex is set to the position of a cell within a row or column, respectively.

-
 If the treegrid includes cells that span multiple rows or multiple columns, and if the `treegrid`role is NOT applied to an HTML`table`element, then aria-rowspan or aria-colspan is applied as described in the Grid and Table Properties Practice.

### Note

- A `treegrid`built from an HTML`table`that includes cells that span multiple rows or columns must use HTML`rowspan`and`colspan`and must not use`aria-rowspan`or`aria-colspan`.
- If rows or cells are included in a treegrid via aria-owns, they will be presented to assistive technologies after the DOM descendants of the `treegrid`element unless the DOM descendants are also included in the`aria-owns`attribute.

# Window Splitter Pattern

## About This Pattern

 **NOTE:** ARIA 1.1 introduced changes to the separator role so it behaves as a widget when focusable.
 While this pattern has been revised to match the ARIA 1.1 specification, the task force will not complete its review until a functional example that matches the ARIA 1.1 specification is complete.
 Progress on this pattern is tracked by issue 129.

A window splitter is a moveable separator between two sections, or panes, of a window that enables users to change the relative size of the panes. A Window Splitter can be either variable or fixed. A fixed splitter toggles between two positions whereas a variable splitter can be adjusted to any position within an allowed range.

A window splitter has a value that represents the size of one of the panes, which, in this pattern, is called the primary pane. When the splitter has its minimum value, the primary pane has its smallest size and the secondary pane has its largest size. The splitter also has an accessible name that matches the name of the primary pane.

 For example, consider a book reading application with a primary pane for the table of contents and a secondary pane that displays content from a section of the book.
 The two panes are divided by a vertical splitter labelled "Table of Contents".
 When the table of contents pane has its maximum size, the splitter has a value of `100`, and when the table of contents is completely collapsed, the splitter has a value of `0`.

Note that the term "primary pane" does not describe the importance or purpose of content inside the pane.

## Example

Work to develop an example window splitter widget is tracked by issue 130.

## Keyboard Interaction

- `Left Arrow`: Moves a vertical splitter to the left.
- `Right Arrow`: Moves a vertical splitter to the right.
- `Up Arrow`: Moves a horizontal splitter up.
- `Down Arrow`: Moves a horizontal splitter down.
-
 `Enter`: If the primary pane is not collapsed, collapses the pane. If the pane is collapsed, restores the splitter to its previous position.
-
 `Home`(Optional): Moves splitter to the position that gives the primary pane its smallest allowed size. This may completely collapse the primary pane.
-
 `End`(Optional): Moves splitter to the position that gives the primary pane its largest allowed size. This may completely collapse the secondary pane.
- `F6`(Optional): Cycle through window panes.

### Note

A fixed size splitter omits implementation of the arrow keys.

## WAI-ARIA Roles, States, and Properties

- The element that serves as the focusable splitter has role separator.
- The separator element has the aria-valuenow property set to a decimal value representing the current position of the separator.
-
 The separator element has the aria-valuemin property set to a decimal value that represents the position where the primary pane has its minimum size.
 This is typically `0`.
-
 The separator element has the aria-valuemax property set to a decimal value that represents the position where the primary pane has its maximum size.
 This is typically `100`.
- If the primary pane has a visible label, it is referenced by aria-labelledby on the separator element. Otherwise, the separator element has a label provided by aria-label.
- The separator element has aria-controls referring to the primary pane.
