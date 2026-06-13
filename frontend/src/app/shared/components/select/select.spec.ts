import { TestBed, ComponentFixture } from '@angular/core/testing';
import { CustomSelectComponent } from './select';

describe('CustomSelectComponent', () => {
  let component: CustomSelectComponent;
  let fixture: ComponentFixture<CustomSelectComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CustomSelectComponent]
    }).compileComponents();

    fixture = TestBed.createComponent(CustomSelectComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should set value and emit', () => {
    vi.spyOn(component.valueChange, 'emit');
    component.value = 'test';
    expect(component.value).toBe('test');
    expect(component.valueChange.emit).toHaveBeenCalledWith('test');
  });

  it('should toggle dropdown', () => {
    expect(component.isOpen).toBe(false);
    component.toggleDropdown();
    expect(component.isOpen).toBe(true);
    component.toggleDropdown();
    expect(component.isOpen).toBe(false);
  });

  it('should format value properly', () => {
    component.options = [{ label: 'Test Label', value: 'TEST' }];
    component.value = 'TEST';
    expect(component.getSelectedLabel()).toBe('Test Label');

    component.value = 'UNKNOWN';
    expect(component.getSelectedLabel()).toBe('Select option');
  });


  it('should select option', () => {
    component.options = ['opt1', 'opt2'];
    component.selectOption('opt1');
    expect(component.value).toBe('opt1');
    expect(component.isOpen).toBe(false);
  });

  it('should get label correctly', () => {
    expect(component.getLabel('opt1')).toBe('opt1');
    component.bindLabel = 'name';
    expect(component.getLabel({ name: 'Opt 2' })).toBe('Opt 2');
  });

  it('should filter options', () => {
    component.options = ['Apple', 'Banana', 'Cherry'];
    component.searchQuery.set('ap');
    expect(component.filteredOptions()).toEqual(['Apple']);
  });

  it('should handle keyboard events', () => {
    component.options = ['A', 'B'];
    component.isOpen = true;
    component.handleKeyDown(new KeyboardEvent('keydown', { key: 'Escape' }));
    expect(component.isOpen).toBe(false);

    component.handleKeyDown(new KeyboardEvent('keydown', { key: 'ArrowDown' }));
    expect(component.isOpen).toBe(true);
    expect(component.focusedIndex).toBe(0);

    component.handleKeyDown(new KeyboardEvent('keydown', { key: 'Enter' }));
    expect(component.value).toBe('A');
  });
});
