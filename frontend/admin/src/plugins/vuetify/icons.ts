import type { IconAliases, IconProps } from 'vuetify'

import checkboxChecked from '@images/svg/checkbox-checked.svg'
import checkboxIndeterminate from '@images/svg/checkbox-indeterminate.svg'
import checkboxUnchecked from '@images/svg/checkbox-unchecked.svg'
import radioChecked from '@images/svg/radio-checked.svg'
import radioUnchecked from '@images/svg/radio-unchecked.svg'

const customIcons: Record<string, unknown> = {
  'mdi-checkbox-blank-outline': checkboxUnchecked,
  'mdi-checkbox-marked': checkboxChecked,
  'mdi-minus-box': checkboxIndeterminate,
  'mdi-radiobox-marked': radioChecked,
  'mdi-radiobox-blank': radioUnchecked,
}

const aliases: Partial<IconAliases> = {
  calendar: 'bx bx-calendar',
  collapse: 'bx bx-chevron-up',
  complete: 'bx bx-check',
  cancel: 'bx bx-x',
  close: 'bx bx-x',
  delete: 'bx bx-bxs-x-circle',
  clear: 'bx bx-x-circle',
  success: 'bx bx-check-circle',
  info: 'bx bx-info-circle',
  warning: 'bx bx-error',
  error: 'bx bx-error-circle',
  prev: 'bx bx-chevron-left',
  ratingEmpty: 'bx bx-star',
  ratingFull: 'bx bx-bxs-star',
  ratingHalf: 'bx bx-bxs-star-half',
  next: 'bx bx-chevron-right',
  delimiter: 'bx bx-circle',
  sort: 'bx bx-up-arrow-alt',
  expand: 'bx bx-chevron-down',
  menu: 'bx bx-menu',
  subgroup: 'bx bx-caret-down',
  dropdown: 'bx bx-chevron-down',
  edit: 'bx bx-pencil',
  loading: 'bx bx-refresh',
  first: 'bx bx-skip-previous',
  last: 'bx bx-skip-next',
  unfold: 'bx bx-move-vertical',
  file: 'bx bx-paperclip',
  plus: 'bx bx-plus',
  minus: 'bx bx-minus',
  sortAsc: 'bx bx-up-arrow-alt',
  sortDesc: 'bx bx-down-arrow-alt',
}

export const iconify = {
  component: (props: IconProps) => {
    // Load custom SVG directly instead of going through icon component
    if (typeof props.icon === 'string') {
      const iconComponent = customIcons[props.icon]

      if (iconComponent)
        return h(iconComponent)
    }

    return h(
      props.tag,
      {
        ...props,

        // As we are using class based icons
        class: [props.icon],

        // Remove used props from DOM rendering
        tag: undefined,
        icon: undefined,
      },
    )
  },
}

export const icons = {
  defaultSet: 'iconify',
  aliases,
  sets: {
    iconify,
  },
}
