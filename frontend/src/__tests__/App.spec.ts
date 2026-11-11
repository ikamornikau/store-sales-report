import { describe, it, expect, vi, afterEach } from 'vitest'
import { flushPromises, mount } from '@vue/test-utils'

import App from '../App.vue'

describe('App', () => {
  afterEach(() => {
    vi.unstubAllGlobals()
  })

  it('loads the report screen from this origin', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        status: 200,
        json: () => Promise.resolve([]),
        text: () => Promise.resolve('[]'),
      }),
    )

    const wrapper = mount(App)
    await flushPromises()

    const request = vi.mocked(fetch).mock.calls[0]?.[0]
    expect(String(request)).toBe(`${window.location.origin}/api/v1/stores/`)
    expect(wrapper.text()).toContain('Store sales')
    expect(wrapper.text()).toContain('No stores are available.')
    wrapper.unmount()
  })
})
