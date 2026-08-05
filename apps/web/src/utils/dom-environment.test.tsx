import { it, expect, describe } from 'vitest';
import { render, screen } from '@testing-library/react';

// Proves the component-test wiring: happy-dom environment, React 19 render,
// and the jest-dom matchers loaded in vitest.setup.ts.
describe('component test environment', () => {
  it('renders a component and matches with jest-dom', () => {
    function Hello({ name }: { name: string }) {
      return <h1>Halo {name}</h1>;
    }

    render(<Hello name="MPP" />);

    expect(screen.getByRole('heading')).toHaveTextContent('Halo MPP');
  });
});
