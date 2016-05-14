import {Injectable} from "@angular/core"
import {SERVICES} from "./pika-service.mock";

@Injectable()
export class PikaServiceService {
    getServices() {
        return Promise.resolve(SERVICES)
    }
}