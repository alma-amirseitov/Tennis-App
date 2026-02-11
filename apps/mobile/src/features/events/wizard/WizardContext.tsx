import { createContext, useContext, useState, type ReactNode } from 'react';
import type { EventType, CompositionType } from '@/shared/api/events';

export interface WizardData {
  event_type: EventType | '';
  composition_type: CompositionType | '';
  level_min: number;
  level_max: number;
  max_participants: number;
  sets_count: number;
  location_name: string;
  start_date: string;
  start_time: string;
  end_time: string;
  title: string;
  description: string;
  community_id: string;
  price: number;
}

const DEFAULT: WizardData = {
  event_type: '',
  composition_type: '',
  level_min: 2.0,
  level_max: 5.0,
  max_participants: 2,
  sets_count: 3,
  location_name: '',
  start_date: '',
  start_time: '',
  end_time: '',
  title: '',
  description: '',
  community_id: '',
  price: 0,
};

interface WizardCtx {
  data: WizardData;
  update: (partial: Partial<WizardData>) => void;
  step: number;
  setStep: (n: number) => void;
  totalSteps: number;
}

const Ctx = createContext<WizardCtx | null>(null);

export function WizardProvider({ children }: { children: ReactNode }) {
  const [data, setData] = useState<WizardData>(DEFAULT);
  const [step, setStep] = useState(0);
  const update = (partial: Partial<WizardData>) => setData((prev) => ({ ...prev, ...partial }));
  return (
    <Ctx.Provider value={{ data, update, step, setStep, totalSteps: 8 }}>
      {children}
    </Ctx.Provider>
  );
}

export function useWizard() {
  const ctx = useContext(Ctx);
  if (!ctx) throw new Error('useWizard must be inside WizardProvider');
  return ctx;
}
