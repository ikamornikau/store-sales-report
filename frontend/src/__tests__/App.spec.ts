import { describe, it, expect, vi, afterEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import App from '../App.vue'

describe('App', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('shows the stores response text', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        text: () => Promise.resolve('Hello from stores list handler!'),
      }),
    )

    const wrapper = mount(App)
    await flushPromises()

    const request = vi.mocked(fetch).mock.calls[0]?.[0]
    expect(String(request)).toBe(`${window.location.origin}/api/v1/stores/`)
    expect(wrapper.text()).toContain('Hello from stores list handler!')
  })
})
