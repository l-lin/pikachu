import {Component, OnInit} from "@angular/core";
import {PikaServiceComponent} from "./pika-service.component";
import {PikaService} from "./pika-service";
import {PikaServiceService} from "./pika-service.service";

@Component({
    selector: 'pikachu',
    template: `
        <h1>Pikachu</h1>
        <ul>
            <li *ngFor="let service of services">
                <pika-service [service]="service"></pika-service>
            </li>
        </ul>
    `,
    directives: [PikaServiceComponent],
    providers: [PikaServiceService]
})
export class AppComponent implements OnInit {
    services: PikaService[];

    constructor(private pikaServiceService: PikaServiceService) {}

    ngOnInit():any {
        this.pikaServiceService.getServices().then(services => this.services = services);
    }
}
