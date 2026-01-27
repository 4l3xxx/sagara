import { Injectable } from '@nestjs/common';
import { InjectRepository } from '@nestjs/typeorm';
import { Repository } from 'typeorm';
import { ConsultationRequest } from './entities/consultation-request.entity';
import { CreateConsultationDto } from './dto/create-consultation.dto';

@Injectable()
export class ContactService {
  constructor(@InjectRepository(ConsultationRequest) private repo: Repository<ConsultationRequest>) {}

  async createConsultation(dto: CreateConsultationDto): Promise<any> {
    const consultation = this.repo.create(dto);
    await this.repo.save(consultation);
    return { success: true, message: 'Consultation request submitted successfully' };
  }

  async findAll(): Promise<any> {
    const consultations = await this.repo.find({ order: { createdAt: 'DESC' } });
    return { success: true, data: consultations };
  }
}
