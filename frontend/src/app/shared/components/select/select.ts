import { Component, Input, Output, EventEmitter, forwardRef, ElementRef, HostListener } from '@angular/core';
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-select',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './select.html',
  styleUrl: './select.scss',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => CustomSelectComponent),
      multi: true
    }
  ]
})
export class CustomSelectComponent implements ControlValueAccessor {
  @Input() options: any[] = [];
  @Input() bindValue?: string;
  @Input() bindLabel?: string;
  @Input() placeholder: string = 'Select option';
  @Input() disabled: boolean = false;
  @Input() id?: string;

  private _value: any = null;

  @Input()
  get value(): any {
    return this._value;
  }

  set value(val: any) {
    if (this._value !== val) {
      this._value = val;
      this.valueChange.emit(val);
      this.onChange(val);
    }
  }

  @Output() valueChange = new EventEmitter<any>();

  isOpen = false;
  focusedIndex = -1;

  onChange: any = () => {};
  onTouched: any = () => {};

  constructor(private elementRef: ElementRef) {}

  writeValue(value: any): void {
    this._value = value;
  }

  registerOnChange(fn: any): void {
    this.onChange = fn;
  }

  registerOnTouched(fn: any): void {
    this.onTouched = fn;
  }

  setDisabledState(isDisabled: boolean): void {
    this.disabled = isDisabled;
  }

  toggleDropdown() {
    if (this.disabled) return;
    this.isOpen = !this.isOpen;
    if (this.isOpen) {
      this.focusedIndex = this.options.findIndex(opt => this.getValue(opt) === this.value);
      if (this.focusedIndex === -1) {
        this.focusedIndex = 0;
      }
    }
  }

  selectOption(option: any) {
    if (this.disabled) return;
    const val = this.getValue(option);
    this.value = val;
    this.isOpen = false;
    this.onTouched();
  }

  getLabel(option: any): string {
    if (option === null || option === undefined) return '';
    if (typeof option === 'object') {
      return this.bindLabel ? option[this.bindLabel] : (option.label || option.name || String(option));
    }
    return String(option);
  }

  getValue(option: any): any {
    if (option === null || option === undefined) return null;
    if (typeof option === 'object') {
      return this.bindValue ? option[this.bindValue] : (option.value !== undefined ? option.value : option);
    }
    return option;
  }

  getSelectedLabel(): string {
    const selected = this.options.find(opt => this.getValue(opt) === this.value);
    return selected ? this.getLabel(selected) : this.placeholder;
  }

  @HostListener('document:click', ['$event'])
  onClickOutside(event: Event) {
    if (!this.elementRef.nativeElement.contains(event.target)) {
      this.isOpen = false;
    }
  }

  @HostListener('keydown', ['$event'])
  handleKeyDown(event: KeyboardEvent) {
    if (this.disabled) return;

    if (event.key === 'Escape') {
      this.isOpen = false;
      event.preventDefault();
    } else if (event.key === 'ArrowDown') {
      if (!this.isOpen) {
        this.isOpen = true;
        this.focusedIndex = 0;
      } else {
        this.focusedIndex = (this.focusedIndex + 1) % this.options.length;
      }
      this.scrollToFocused();
      event.preventDefault();
    } else if (event.key === 'ArrowUp') {
      if (!this.isOpen) {
        this.isOpen = true;
        this.focusedIndex = this.options.length - 1;
      } else {
        this.focusedIndex = (this.focusedIndex - 1 + this.options.length) % this.options.length;
      }
      this.scrollToFocused();
      event.preventDefault();
    } else if (event.key === 'Enter' || event.key === ' ') {
      if (this.isOpen) {
        if (this.focusedIndex >= 0 && this.focusedIndex < this.options.length) {
          this.selectOption(this.options[this.focusedIndex]);
        }
      } else {
        this.isOpen = true;
      }
      event.preventDefault();
    }
  }

  private scrollToFocused() {
    setTimeout(() => {
      const activeEl = this.elementRef.nativeElement.querySelector('[role="option"].bg-zinc-900');
      if (activeEl) {
        activeEl.scrollIntoView({ block: 'nearest' });
      }
    }, 0);
  }
}
