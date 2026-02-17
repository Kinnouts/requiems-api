## Overview

The Requiems API Dashboard is fully responsive and optimized for mobile, tablet,
and desktop devices using Tailwind CSS's mobile-first approach.

## Breakpoints

We use Tailwind's default breakpoints:

- **Mobile**: < 640px (base, no prefix)
- **Small (sm)**: ≥ 640px
- **Medium (md)**: ≥ 768px
- **Large (lg)**: ≥ 1024px
- **Extra Large (xl)**: ≥ 1280px
- **2XL**: ≥ 1536px

## Key Responsive Patterns

### 1. Layouts

#### Dashboard Layout (`layouts/dashboard.html.erb`)

```erb
<div class="lg:grid lg:grid-cols-12 lg:gap-8">
  <aside class="lg:col-span-3">
    <!-- Sidebar: Full width on mobile, 3 cols on desktop -->
  </aside>
  <main class="lg:col-span-9 mt-8 lg:mt-0">
    <!-- Main content: Full width on mobile, 9 cols on desktop -->
  </main>
</div>
```

- **Mobile**: Sidebar and content stack vertically
- **Desktop**: Sidebar (25%) and content (75%) side-by-side

#### Admin Layout (`layouts/admin.html.erb`)

- Same grid pattern as dashboard layout
- Admin sidebar sticky on desktop, normal flow on mobile

### 2. Navigation

#### Desktop Navigation

```erb
<div class="hidden sm:ml-8 sm:flex sm:space-x-8">
  <!-- Navigation links visible on sm+ -->
</div>
```

#### Mobile Navigation

```erb
<div class="flex items-center sm:hidden">
  <button data-controller="dropdown">
    <!-- Hamburger menu button -->
  </button>
</div>

<div class="sm:hidden hidden" data-dropdown-target="menu">
  <!-- Mobile menu panel -->
</div>
```

- **Mobile**: Hamburger menu with dropdown
- **Desktop**: Horizontal navigation bar

### 3. Grid Layouts

#### Stats Cards

```erb
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
  <!-- Cards adapt based on screen size -->
</div>
```

- **Mobile**: 1 column
- **Tablet**: 2 columns
- **Desktop**: 4 columns

#### Content Grid

```erb
<div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
  <div class="lg:col-span-2"><!-- Main content --></div>
  <div><!-- Sidebar --></div>
</div>
```

### 4. Tables

All tables wrapped in overflow containers:

```erb
<div class="overflow-x-auto">
  <table class="min-w-full divide-y divide-gray-200">
    <!-- Table content -->
  </table>
</div>
```

- **Mobile**: Horizontal scroll enabled
- **Desktop**: Full table display

### 5. Typography

```erb
<h1 class="text-2xl sm:text-3xl font-bold">
  <!-- Smaller on mobile, larger on desktop -->
</h1>
```

### 6. Spacing

Responsive padding and margins:

```erb
<div class="px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
  <!-- Less padding on mobile, more on desktop -->
</div>
```

### 7. Forms

```erb
<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
  <div class="md:col-span-2">
    <!-- Search input spans 2 columns on tablet+ -->
  </div>
  <div>
    <!-- Filter dropdowns -->
  </div>
</div>
```

## Component Responsiveness

### Buttons

```erb
<button class="w-full sm:w-auto px-4 py-2">
  <!-- Full width on mobile, auto on desktop -->
</button>
```

### Modals

```erb
<div class="w-full max-w-md">
  <!-- Constrained width with mobile-first approach -->
</div>
```

### Cards

```erb
<div class="bg-white rounded-lg shadow p-4 sm:p-6">
  <!-- Less padding on mobile -->
</div>
```

## Mobile-Specific Optimizations

### 1. Touch Targets

All interactive elements have minimum 44px touch targets:

```erb
<button class="px-4 py-3">
  <!-- Adequate padding for touch -->
</button>
```

### 2. Sticky Navigation

```erb
<nav class="sticky top-0 z-40">
  <!-- Navbar stays visible on scroll -->
</nav>
```

### 3. Truncated Text

```erb
<p class="truncate max-w-xs">
  <!-- Long text truncated on small screens -->
</p>
```

### 4. Hidden Elements

```erb
<!-- Hide on mobile, show on desktop -->
<div class="hidden lg:block">

<!-- Show on mobile, hide on desktop -->
<div class="lg:hidden">
```

## Testing Checklist

### Mobile (375px - 639px)

- ✅ Navbar collapses to hamburger menu
- ✅ Tables scroll horizontally
- ✅ Forms stack vertically
- ✅ Cards display in single column
- ✅ Touch targets are adequately sized
- ✅ Text is readable without zooming

### Tablet (640px - 1023px)

- ✅ Navigation partially visible
- ✅ Grid layouts show 2 columns
- ✅ Sidebar remains stacked or side-by-side
- ✅ Tables display properly

### Desktop (1024px+)

- ✅ Full navigation visible
- ✅ Sidebars fixed/sticky
- ✅ Multi-column grid layouts
- ✅ Tables display without scrolling
- ✅ All content properly spaced

## Browser Support

- Chrome 90+
- Firefox 88+
- Safari 14+
- Edge 90+

## Accessibility

All responsive elements maintain:

- Proper ARIA labels
- Keyboard navigation
- Screen reader compatibility
- Sufficient color contrast
- Focus indicators

## Performance

- Mobile-first CSS (base styles for mobile, media queries for desktop)
- No layout shift during load
- Optimized images with responsive sizing
- Minimal JavaScript for responsive features

## Common Patterns

### Dashboard/Admin Sidebar

```css
/* Mobile: Full width, stacks above content */
/* Desktop: 3-column width, sticky positioning */
class="lg:col-span-3 sticky top-8"
```

### Content Areas

```css
/* Mobile: Full width with margin-top */
/* Desktop: 9-column width, no top margin */
class="lg:col-span-9 mt-8 lg:mt-0"
```

### Action Buttons

```css
/* Mobile: Full width, stacked */
/* Desktop: Inline, auto width */
class="flex flex-col sm:flex-row gap-4"
```

## Future Enhancements

- [ ] Add PWA manifest for mobile app-like experience
- [ ] Implement touch gestures for table navigation
- [ ] Add dark mode support
- [ ] Optimize for foldable devices
- [ ] Add tablet-specific layouts (768px - 1024px)

## Resources

- [Tailwind CSS Responsive Design](https://tailwindcss.com/docs/responsive-design)
- [Mobile-First CSS](https://developer.mozilla.org/en-US/docs/Web/Progressive_web_apps/Responsive/Mobile_first)
- [Touch Target Guidelines](https://www.w3.org/WAI/WCAG21/Understanding/target-size.html)
