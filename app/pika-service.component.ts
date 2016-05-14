import {Component, Input} from "@angular/core";

import {PikaService} from "./pika-service"

@Component({
    selector: 'pika-service',
    template: `
        <div *ngIf="service">
            <h2>{{ service.id }} - {{ service.name }}</h2>
        </div>
    `
})
export class PikaServiceComponent{
    @Input()
    service: PikaService
}
